// Package files is the hand-written convenience layer for the org file
// library: upload a local file and push it to a phone, on top of the generated
// files + phones REST clients.
//
//	f, _ := files.Upload(ctx, c, "./demo.mp4")           // register + PUT bytes
//	files.Push(ctx, c, "phn_abc", f.ID)                  // reuse across phones
//	files.Send(ctx, c, "phn_abc", "./demo.mp4",          // one-shot: upload + push
//		files.WithWait(60*time.Second))                  // ...and wait for delivery
//
// Free functions taking the generated *client.Client — the Go twin of
// platform-python's client.files.upload / client.phones.send_file, and the same
// idiom as drivers/mobile (the generated types can't be extended, so the
// value-add lives beside them). Preserved across regen by scripts/regen.sh's
// drivers/ exclude.
package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	platformgo "github.com/axilioai/platform-go"
	client "github.com/axilioai/platform-go/client"
)

// defaultMIME is sent when the extension maps to nothing known; the backend
// MIME whitelist has the final say.
const defaultMIME = "application/octet-stream"

const (
	defaultWaitTimeout = 60 * time.Second
	defaultPollEvery   = 2 * time.Second
)

// terminal reports whether a delivery has reached a final status.
func terminal(s platformgo.FileDeliverySummaryStatus) bool {
	return s == platformgo.FileDeliverySummaryStatusDelivered ||
		s == platformgo.FileDeliverySummaryStatusFailed
}

type config struct {
	filename   string
	mimeType   string
	collection string
	wait       bool
	timeout    time.Duration
	pollEvery  time.Duration
}

// Option configures Upload / Push / Send. Each function reads only the options
// that apply to it (documented per function).
type Option func(*config)

// WithFilename overrides the registered filename (default: base of the path).
func WithFilename(name string) Option { return func(c *config) { c.filename = name } }

// WithMimeType overrides the content type (default: guessed from the extension).
func WithMimeType(mimeType string) Option { return func(c *config) { c.mimeType = mimeType } }

// WithCollection overrides the MediaStore collection (DCIM / Pictures / Movies;
// default by media class server-side).
func WithCollection(collection string) Option { return func(c *config) { c.collection = collection } }

// WithWait makes Send block until the phone reports terminal status (or timeout
// elapses), returning the latest delivery either way. A zero or negative
// timeout uses the 60s default.
func WithWait(timeout time.Duration) Option {
	return func(c *config) {
		c.wait = true
		if timeout > 0 {
			c.timeout = timeout
		}
	}
}

// WithPollInterval sets the poll cadence used while waiting (default 2s).
func WithPollInterval(d time.Duration) Option {
	return func(c *config) {
		if d > 0 {
			c.pollEvery = d
		}
	}
}

func resolve(opts ...Option) config {
	c := config{timeout: defaultWaitTimeout, pollEvery: defaultPollEvery}
	for _, o := range opts {
		o(&c)
	}
	return c
}

// Upload registers a local file and uploads its bytes to the org library.
// Returns the FileSummary; its ID is what Push / Send take. Reads WithFilename,
// WithMimeType.
func Upload(ctx context.Context, c *client.Client, path string, opts ...Option) (*platformgo.FileSummary, error) {
	cfg := resolve(opts...)
	name := cfg.filename
	if name == "" {
		name = filepath.Base(path)
	}
	mimeType := cfg.mimeType
	if mimeType == "" {
		if mimeType = mime.TypeByExtension(filepath.Ext(name)); mimeType == "" {
			mimeType = defaultMIME
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	registered, err := c.Files.Create(ctx, &platformgo.FileCreateRequest{
		Filename:  name,
		MimeType:  mimeType,
		SizeBytes: int64(len(data)),
	})
	if err != nil {
		return nil, err
	}
	// The presigned PUT goes straight to object storage: no Axilio auth header,
	// and Content-Type must match what was registered (the push HeadObject-
	// verifies size + type). ContentLength matches size_bytes.
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, registered.UploadURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mimeType)
	req.ContentLength = int64(len(data))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("upload PUT to storage failed: %s: %s", resp.Status, body)
	}
	return registered.File, nil
}

// Push sends an already-uploaded library file to a phone. Returns the delivery
// (status dispatched once the phone acks). Reads WithCollection.
func Push(ctx context.Context, c *client.Client, phoneID, fileID string, opts ...Option) (*platformgo.FileDeliverySummary, error) {
	cfg := resolve(opts...)
	req := &platformgo.PhonesPushFileRequest{PhoneID: phoneID, FileID: fileID}
	if cfg.collection != "" {
		col := platformgo.PhonesPushFileRequestCollection(cfg.collection)
		req.Collection = &col
	}
	resp, err := c.Phones.PushFile(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Delivery, nil
}

// Send uploads a local file and pushes it to a phone in one call. Returns the
// delivery after dispatch; with WithWait it polls until terminal status (or
// timeout), returning the latest delivery either way — inspect Status / Error.
// Reads all options.
func Send(ctx context.Context, c *client.Client, phoneID, path string, opts ...Option) (*platformgo.FileDeliverySummary, error) {
	cfg := resolve(opts...)
	f, err := Upload(ctx, c, path, opts...)
	if err != nil {
		return nil, err
	}
	delivery, err := Push(ctx, c, phoneID, f.ID, opts...)
	if err != nil {
		return nil, err
	}
	if !cfg.wait {
		return delivery, nil
	}
	return awaitTerminal(ctx, c, phoneID, delivery, cfg.timeout, cfg.pollEvery)
}

func awaitTerminal(
	ctx context.Context,
	c *client.Client,
	phoneID string,
	delivery *platformgo.FileDeliverySummary,
	timeout, pollEvery time.Duration,
) (*platformgo.FileDeliverySummary, error) {
	deadline := time.Now().Add(timeout)
	limit := int64(100)
	for !terminal(delivery.Status) && time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return delivery, ctx.Err()
		case <-time.After(pollEvery):
		}
		page, err := c.Phones.ListFiles(ctx, &platformgo.PhonesListFilesRequest{PhoneID: phoneID, Limit: &limit})
		if err != nil {
			return delivery, err
		}
		for _, cand := range page.Deliveries {
			if cand.ID == delivery.ID {
				delivery = cand
				break
			}
		}
	}
	return delivery, nil
}
