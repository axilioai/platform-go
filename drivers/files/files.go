// Package files is the hand-written convenience layer for the org file
// library: upload a local file, manage the library, and deliver a file to a
// phone, on top of the generated uploads + phones REST clients.
//
//	f, _ := files.Upload(ctx, c, "./demo.mp4")           // register + PUT + complete
//	files.Push(ctx, c, "phn_abc", f.ID)                  // reuse across phones
//	files.Send(ctx, c, "phn_abc", "./demo.mp4",          // one-shot: upload + deliver
//		files.WithWait(60*time.Second))                  // ...and wait for delivery
//	files.List(ctx, c)                                   // what's in the library
//	files.Delete(ctx, c, f.ID)                           // free the quota again
//
// Free functions taking the generated *client.Client — the Go twin of
// platform-python's client.files helpers, and the same idiom as drivers/mobile
// (the generated types can't be extended, so the value-add lives beside them).
// Preserved across regen by scripts/regen.sh's drivers/ exclude.
//
// Vocabulary, one word per concept: UPLOAD puts a local file into the library,
// PUSH sends a library file to a phone, SEND does both. The API these call is
// named by direction — /uploads for what you put in, /phones/{id}/deliveries
// for the record of what we sent to a phone.
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
	registered, err := c.Uploads.Create(ctx, &platformgo.FileCreateRequest{
		Filename:  name,
		MimeType:  mimeType,
		SizeBytes: int64(len(data)),
	})
	if err != nil {
		return nil, err
	}
	// The presigned PUT goes straight to object storage: no Axilio auth header,
	// and Content-Type must match what was registered (the signature pins both
	// type and length). ContentLength matches size_bytes.
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
	// Completion is what makes the file deliverable: the server verifies the
	// object landed at the declared size and type, checks the content really is
	// the media it claims, and flips the row to ready. Skipping it leaves a file
	// stuck 'uploading' that every delivery would reject — previously the first
	// push did this verification lazily, which only worked because Send always
	// pushed. Upload on its own now finishes the job.
	completed, err := c.Uploads.Complete(ctx, &platformgo.UploadsCompleteRequest{
		UploadID: registered.File.GetID(),
	})
	if err != nil {
		return nil, err
	}
	return completed.File, nil
}

// List returns one page of the org's library along with its standing usage
// against the storage quota, so a caller can show "X of Y" without a second
// call. Options are the generated request's own (limit/offset/search/sort).
func List(ctx context.Context, c *client.Client, request *platformgo.UploadsListRequest) (*platformgo.FileListResponse, error) {
	if request == nil {
		request = &platformgo.UploadsListRequest{}
	}
	return c.Uploads.List(ctx, request)
}

// Delete removes a file from the org's library: the stored object, the entry,
// and its delivery history. This is the other half of a quota — without it a
// caller can fill the library and has no supported way to clear it.
func Delete(ctx context.Context, c *client.Client, uploadID string) error {
	_, err := c.Uploads.Delete(ctx, &platformgo.UploadsDeleteRequest{UploadID: uploadID})
	return err
}

// Push sends an already-uploaded library file to a phone. Returns the delivery
// (status dispatched once the phone acks). Reads WithCollection.
func Push(ctx context.Context, c *client.Client, phoneID, fileID string, opts ...Option) (*platformgo.FileDeliverySummary, error) {
	cfg := resolve(opts...)
	req := &platformgo.FileDeliveryCreateRequest{PhoneID: phoneID, FileID: fileID}
	if cfg.collection != "" {
		col := platformgo.FileDeliveryCreateRequestCollection(cfg.collection)
		req.Collection = &col
	}
	resp, err := c.Phones.CreateDelivery(ctx, req)
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
	for !terminal(delivery.Status) && time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return delivery, ctx.Err()
		case <-time.After(pollEvery):
		}
		// Fetch OUR delivery by id. This used to page the newest 100 deliveries
		// and scan for a match, which silently lost the target on a busy phone:
		// once 100 newer pushes landed, the delivery being waited on fell off
		// the page, the loop stopped updating it, and the caller got a stale
		// non-terminal record back as if it had simply timed out. The per-
		// delivery endpoint has no such window.
		latest, err := c.Phones.GetDelivery(ctx, &platformgo.PhonesGetDeliveryRequest{
			PhoneID:    phoneID,
			DeliveryID: delivery.GetID(),
		})
		if err != nil {
			return delivery, err
		}
		delivery = latest
	}
	return delivery, nil
}
