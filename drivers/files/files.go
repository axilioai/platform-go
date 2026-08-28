// Package files is the hand-written convenience layer for the org file
// library: upload a local file, manage the library, and deliver a file to a
// phone, on top of the generated files + phones REST clients.
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
// PUSH sends a library file to a phone, SEND does both. These call the unified
// /files collection (a file's provenance is the source attribute: source=upload
// for what you put in) and /phones/{id}/deliveries for the record of what we
// sent to a phone.
package files

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	platformgo "github.com/axilioai/platform-go"
	client "github.com/axilioai/platform-go/client"
)

// MaxDeliveryBytes mirrors the server's per-delivery ceiling (100 MiB): the
// phone downloads over its own cellular link, so the bound belongs to that
// transport — the library itself stores files up to 1 GiB, including files
// no phone can receive. Mirrored rather than fetched, and the server stays
// authoritative (it rejects an oversize push regardless); this constant only
// lets Send refuse BEFORE uploading a file that could never be delivered,
// which would otherwise be retained in the library by a failed one-shot
// call. The number is pinned by a backend regression test (AXI-1581), so a
// server-side change shows up as a failing test there, not a silent drift.
const MaxDeliveryBytes int64 = 100 << 20 // 100 MiB

// ErrTooLargeForDelivery marks a Send refused before upload because the local
// file exceeds MaxDeliveryBytes. Match with errors.Is; the wrapped message
// carries the sizes. Upload deliberately does NOT apply this bound — the
// library accepts what phones cannot receive, and Upload is the library door.
var ErrTooLargeForDelivery = errors.New("file exceeds the 100 MiB phone-delivery limit")

// defaultMIME is sent when the extension maps to nothing known; the backend
// MIME whitelist has the final say.
const defaultMIME = "application/octet-stream"

// extMIME pins the content type for every extension the API's upload
// whitelist accepts. mime.TypeByExtension reads the HOST's MIME database,
// which varies by machine and container image and does not reliably know
// .heic — so relying on it alone made a supported format upload as
// application/octet-stream on some machines and succeed on others. The API
// pins Content-Type into the presigned PUT, so getting this wrong is not a
// cosmetic difference: the upload is registered under a type the server then
// rejects.
var extMIME = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
	".gif":  "image/gif",
	".heic": "image/heic",
	".mp4":  "video/mp4",
	".webm": "video/webm",
	".mov":  "video/quicktime",
	".3gp":  "video/3gpp",
	".mkv":  "video/x-matroska",
}

// detectMIME resolves a filename to a content type: our own table first, the
// host database second, and the generic default last.
func detectMIME(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	if m, ok := extMIME[ext]; ok {
		return m
	}
	if m := mime.TypeByExtension(ext); m != "" {
		// TypeByExtension can return parameters ("text/plain; charset=utf-8").
		// The API compares the bare type, so hand it the bare type.
		if base, _, err := mime.ParseMediaType(m); err == nil {
			return base
		}
		return m
	}
	return defaultMIME
}

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
	httpClient *http.Client
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

// WithWait makes Push and Send block until the phone reports terminal status
// (or timeout elapses), returning the latest delivery either way. A zero or
// negative timeout uses the 60s default.
func WithWait(timeout time.Duration) Option {
	return func(c *config) {
		c.wait = true
		if timeout > 0 {
			c.timeout = timeout
		}
	}
}

// WithHTTPClient overrides the client used for the storage PUT (default
// http.DefaultClient). Useful for a custom timeout or transport, and for
// testing against a stub server. A nil client is ignored rather than
// installed, so a zero value can't disable uploads.
func WithHTTPClient(h *http.Client) Option {
	return func(c *config) {
		if h != nil {
			c.httpClient = h
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
	c := config{
		timeout:    defaultWaitTimeout,
		pollEvery:  defaultPollEvery,
		httpClient: http.DefaultClient,
	}
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

	// Stream the file rather than reading it into memory. The library accepts
	// up to 1 GB per file, so os.ReadFile would allocate the whole payload —
	// tolerable under the old 5 MiB cap, not under this one. Passing the open
	// *os.File as the request body lets net/http copy it through in chunks.
	fh, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer fh.Close()
	info, err := fh.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory", path)
	}
	size := info.Size()

	name := cfg.filename
	if name == "" {
		name = filepath.Base(path)
	}
	mimeType := cfg.mimeType
	if mimeType == "" {
		mimeType = detectMIME(name)
	}

	registered, err := c.Files.Create(ctx, &platformgo.FileCreateRequest{
		Filename:  name,
		MimeType:  mimeType,
		SizeBytes: size,
	})
	if err != nil {
		return nil, fmt.Errorf("register upload: %w", err)
	}
	// The presigned PUT goes straight to object storage: no Axilio auth header,
	// and Content-Type must match what was registered (the signature pins both
	// type and length). Setting ContentLength explicitly keeps the request out
	// of chunked encoding, which the signature would reject.
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, registered.UploadURL, fh)
	if err != nil {
		return nil, fmt.Errorf("build storage request: %w", err)
	}
	req.Header.Set("Content-Type", mimeType)
	req.ContentLength = size
	resp, err := cfg.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload %s to storage: %w", name, err)
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
	completed, err := c.Files.Complete(ctx, &platformgo.FilesCompleteRequest{
		FileID: registered.File.GetID(),
	})
	if err != nil {
		return nil, fmt.Errorf("complete upload %s: %w", registered.File.GetID(), err)
	}
	return completed.File, nil
}

// List returns one page of the org's library along with its standing usage
// against the storage quota, so a caller can show "X of Y" without a second
// call. Options are the generated request's own (limit/offset/search/sort).
func List(ctx context.Context, c *client.Client, request *platformgo.FilesListRequest) (*platformgo.FileListResponse, error) {
	if request == nil {
		request = &platformgo.FilesListRequest{}
	}
	return c.Files.List(ctx, request)
}

// Delete removes a file from the org's library: the stored object, the entry,
// and its delivery history. This is the other half of a quota — without it a
// caller can fill the library and has no supported way to clear it.
func Delete(ctx context.Context, c *client.Client, fileID string) error {
	_, err := c.Files.Delete(ctx, &platformgo.FilesDeleteRequest{FileID: fileID})
	return err
}

// Push sends an already-uploaded library file to a phone. Returns the delivery
// after dispatch; with WithWait it polls until the phone reports terminal
// status (or timeout), returning the latest delivery either way — inspect
// Status / Error. Reads WithCollection, WithWait, WithPollInterval.
//
// Waiting used to live only in Send, which made pushing one stored file to
// several phones the one flow that could not be waited on — the case Push
// exists for. The wait now belongs to the delivery, not to how the file got
// into the library.
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
	if !cfg.wait {
		return resp.Delivery, nil
	}
	return awaitTerminal(ctx, c, phoneID, resp.Delivery, cfg.timeout, cfg.pollEvery)
}

// Send uploads a local file and pushes it to a phone in one call: Upload then
// Push, with every option forwarded to both. Returns whatever Push returns, so
// WithWait behaves identically here and on a bare Push. Reads all options.
//
// Send preflights MaxDeliveryBytes before any request goes out (AXI-1581):
// a file the delivery endpoint would refuse must not be uploaded first, or
// the failed one-shot call leaves it retained in the library, consuming
// quota. The check lives on Send and not Upload because only Send promises
// delivery; Upload keeps the library's own 1 GiB contract.
func Send(ctx context.Context, c *client.Client, phoneID, path string, opts ...Option) (*platformgo.FileDeliverySummary, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}
	if info.Size() > MaxDeliveryBytes {
		return nil, fmt.Errorf("%s is %d bytes: %w (%d bytes); the library itself stores up to 1 GiB — use Upload and Push separately if you only need it stored",
			filepath.Base(path), info.Size(), ErrTooLargeForDelivery, MaxDeliveryBytes)
	}
	f, err := Upload(ctx, c, path, opts...)
	if err != nil {
		return nil, err
	}
	return Push(ctx, c, phoneID, f.ID, opts...)
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
