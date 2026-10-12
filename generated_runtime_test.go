package api_test

// Behaviour of the repo-owned overrides of Fern's generated runtime
// (internal/explicit_fields.go, internal/retrier.go, internal/stringer.go),
// checked through the generated client the way callers use it. regen.sh keeps
// those files, and this one, out of every regen; see the comment there.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	api "github.com/axilioai/platform-go"
	"github.com/axilioai/platform-go/client"
	"github.com/axilioai/platform-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCreateRequestSetterKeepsRuns(t *testing.T) {
	start := time.Date(2026, 10, 12, 9, 0, 0, 0, time.UTC)
	request := &api.RunCreateRequest{WorkflowID: "wf_123", Runs: []*api.RunConfig{{}}}
	request.SetStartAt(&start)

	body, err := json.Marshal(request)
	require.NoError(t, err)
	assert.JSONEq(t, `{"runs":[{}],"start_at":"2026-10-12T09:00:00Z"}`, string(body))
}

func TestRunCreateRequestExplicitNilStartAt(t *testing.T) {
	request := &api.RunCreateRequest{Runs: []*api.RunConfig{{}}}
	request.SetStartAt(nil)

	body, err := json.Marshal(request)
	require.NoError(t, err)
	assert.JSONEq(t, `{"runs":[{}],"start_at":null}`, string(body))
}

func TestWebhookEndpointSetterKeepsFields(t *testing.T) {
	endpoint := &api.WebhookEndpoint{ID: "we_123", Type: api.WebhookEndpointTypeWebhook}
	endpoint.SetCreatedAt(time.Date(2026, 10, 12, 9, 0, 0, 0, time.UTC))

	body, err := json.Marshal(endpoint)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"id": "we_123",
		"type": "webhook",
		"disabled": false,
		"created_at": "2026-10-12T09:00:00Z"
	}`, string(body))
}

// countingServer answers each request with the next status in statuses
// (repeating the last), and counts the requests it saw.
func countingServer(t *testing.T, statuses ...int) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var count atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := int(count.Add(1))
		status := statuses[min(n, len(statuses))-1]
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(status)
		if status < 300 {
			_, _ = w.Write([]byte(`{"id":"we_123","secret":"whsec_test","endpoints":[]}`))
		}
	}))
	t.Cleanup(server.Close)
	return server, &count
}

func TestWebhookCreateIsNotRetriedAfterServerError(t *testing.T) {
	server, count := countingServer(t, http.StatusInternalServerError, http.StatusCreated)
	c := client.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("ax_test"))

	url := "https://example.com/hook"
	_, err := c.Webhooks.CreateEndpoint(context.Background(), &api.WebhookEndpointCreateRequest{URL: &url})
	require.Error(t, err)
	assert.Equal(t, int32(1), count.Load(), "a POST that may have taken effect must not be sent twice")
}

func TestWebhookCreateIsRetriedAfterTooManyRequests(t *testing.T) {
	server, count := countingServer(t, http.StatusTooManyRequests, http.StatusCreated)
	c := client.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("ax_test"))

	url := "https://example.com/hook"
	_, err := c.Webhooks.CreateEndpoint(context.Background(), &api.WebhookEndpointCreateRequest{URL: &url})
	require.NoError(t, err)
	assert.Equal(t, int32(2), count.Load(), "a 429 was refused before it ran, so it is safe to resend")
}

func TestListIsRetriedAfterServerError(t *testing.T) {
	server, count := countingServer(t, http.StatusInternalServerError, http.StatusOK)
	c := client.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("ax_test"))

	_, err := c.Webhooks.ListEndpoints(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int32(2), count.Load())
}

func TestStringRedactsSecrets(t *testing.T) {
	cases := []struct {
		name   string
		value  interface{ String() string }
		secret string
		kept   string
	}{
		{"webhook endpoint create", &api.WebhookEndpointCreateResponse{}, "whsec_create", `"we_1"`},
		{"webhook secret rotation", &api.WebhookRotateSecretResponse{}, "whsec_rotate", `"we_1"`},
		{"API key create", &api.APIKeyCreateResponse{}, "ax_live_key", `"key_1"`},
		{"phone allocation live view", &api.PhoneAllocateResponse{}, "tok_live_view", `"ph_1"`},
		{"phone allocation control", &api.PhoneAllocateResponse{}, "tok_control", `"ph_1"`},
		{"phone allocation telemetry", &api.PhoneAllocateResponse{}, "tok_telemetry", `"ph_1"`},
	}
	payloads := map[string]string{
		"webhook endpoint create":    `{"id":"we_1","type":"webhook","created_at":"2026-10-12T09:00:00Z","secret":"whsec_create"}`,
		"webhook secret rotation":    `{"id":"we_1","secret":"whsec_rotate"}`,
		"API key create":             `{"id":"key_1","name":"ci","created_at":"2026-10-12T09:00:00Z","key_value":"ax_live_key"}`,
		"phone allocation live view": `{"phone_id":"ph_1","session_id":"ss_1","workflow_started_at":"2026-10-12T09:00:00Z","live_view_url":"https://live.test/s?token=tok_live_view"}`,
		"phone allocation control":   `{"phone_id":"ph_1","session_id":"ss_1","workflow_started_at":"2026-10-12T09:00:00Z","control_url":"wss://connect.test/ws/control?token=tok_control"}`,
		"phone allocation telemetry": `{"phone_id":"ph_1","session_id":"ss_1","workflow_started_at":"2026-10-12T09:00:00Z","telemetry_url":"wss://connect.test/ws/telemetry?token=tok_telemetry"}`,
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, json.Unmarshal([]byte(payloads[tc.name]), tc.value))
			printed := tc.value.String()
			assert.NotContains(t, printed, tc.secret)
			assert.Contains(t, printed, `"[redacted]"`)
			assert.Contains(t, printed, tc.kept, "the rest of the value still prints")
		})
	}
}
