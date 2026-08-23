// Tolerant-reader behavior of the unified frame envelope (AXI-1851).
//
// The contract shared by the REST archive (SessionsListFrames) and the live
// telemetry WebSocket ("live and archive differ only in cardinality"):
// an unknown frame kind surfaces as an explicit unknown variant carrying the
// raw JSON — never an error, never a silent drop; unknown fields inside known
// kinds are ignored (kept as extra properties); unknown span_type/log_type
// values parse without error.
//
// Hand-written. regen.sh excludes this file from the rsync --delete sweep;
// if a regen changes the generated union's tolerance, these tests fail loudly.
package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const unknownKindFrame = `{"kind":"telemetry_v2","payload":{"nested":true}}`

const framesResponsePage = `{
	"frames": [
		{"kind":"span","name":"Screen.observe","phase":"end","span_id":"aaaaaaaaaaaaaaaa",
		 "span_type":"quantum_leap","start_time_unix_nano":1,"end_time_unix_nano":2,
		 "status":{"code":"ok","message":""},"trace_id":"b","brand_new_field":{"x":1}},
		{"kind":"telemetry_v2","payload":{"nested":true}},
		{"kind":"log","body":"hello","log_type":"quantum_log","severity":"info",
		 "time_unix_nano":3,"trace_id":"b"}
	],
	"inference_costs":{},"limit":100,"offset":0,"retention_expired":false,"sdk_call_costs":{}
}`

func TestFramesTolerantReader(t *testing.T) {
	t.Run("unknown kind surfaces as explicit unknown variant with raw JSON", func(t *testing.T) {
		var item RunSessionFramesResponseFramesItem
		require.NoError(t, json.Unmarshal([]byte(unknownKindFrame), &item))
		assert.Equal(t, "telemetry_v2", item.GetKind())
		assert.Nil(t, item.GetLog())
		assert.Nil(t, item.GetSpan())

		// The raw JSON is carried verbatim: re-marshaling an unknown frame
		// yields the original bytes, so nothing is lost in a read/write pass.
		remarshaled, err := json.Marshal(item)
		require.NoError(t, err)
		assert.JSONEq(t, unknownKindFrame, string(remarshaled))
	})

	t.Run("full page parses with nothing dropped", func(t *testing.T) {
		var response RunSessionFramesResponse
		require.NoError(t, json.Unmarshal([]byte(framesResponsePage), &response))
		require.Len(t, response.GetFrames(), 3)
		assert.NotNil(t, response.GetFrames()[0].GetSpan())
		assert.Equal(t, "telemetry_v2", response.GetFrames()[1].GetKind())
		assert.NotNil(t, response.GetFrames()[2].GetLog())
	})

	t.Run("unknown field in a known kind is kept as an extra property", func(t *testing.T) {
		var response RunSessionFramesResponse
		require.NoError(t, json.Unmarshal([]byte(framesResponsePage), &response))
		span := response.GetFrames()[0].GetSpan()
		require.NotNil(t, span)
		assert.Contains(t, span.GetExtraProperties(), "brand_new_field")
	})

	t.Run("unknown span_type and log_type values parse", func(t *testing.T) {
		var response RunSessionFramesResponse
		require.NoError(t, json.Unmarshal([]byte(framesResponsePage), &response))
		assert.Equal(t, "quantum_leap", response.GetFrames()[0].GetSpan().GetSpanType())
		assert.Equal(t, "quantum_log", response.GetFrames()[2].GetLog().GetLogType())
	})
}
