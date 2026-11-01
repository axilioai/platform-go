package internal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringifyJSONRedactsSecrets(t *testing.T) {
	raw := json.RawMessage(`{"id":"we_1","secret":"whsec_x","nested":[{"token":"tok_x","n":12345678901234567890}],"key_value":null}`)

	printed, err := StringifyJSON(raw)
	require.NoError(t, err)
	assert.Equal(t, `{
  "id": "we_1",
  "secret": "[redacted]",
  "nested": [
    {
      "token": "[redacted]",
      "n": 12345678901234567890
    }
  ],
  "key_value": null
}`, printed)
}

func TestStringifyJSONUnchangedWithoutSecrets(t *testing.T) {
	value := struct {
		Name  string   `json:"name"`
		Tags  []string `json:"tags"`
		Count int      `json:"count"`
		Link  string   `json:"link"`
	}{"b", []string{"x", "y"}, 2, "a<b>&c"}

	want, err := json.MarshalIndent(value, "", "  ")
	require.NoError(t, err)
	printed, err := StringifyJSON(value)
	require.NoError(t, err)
	assert.Equal(t, string(want), printed)
}

func TestStringifyJSONRedactsRunStreamURL(t *testing.T) {
	out, err := StringifyJSON(json.RawMessage(`{"id":"run_1","stream_url":"https://connect.test/events?token=tok_stream"}`))
	require.NoError(t, err)
	assert.NotContains(t, out, "tok_stream")
	assert.Contains(t, out, "[redacted]")
}
