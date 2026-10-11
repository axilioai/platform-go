package internal

// Repo-owned: regen.sh keeps this file out of every regen. Fern's version
// printed every field, so logging a create or rotate response wrote its
// one-time secret to the log. See generated_runtime_test.go.

import (
	"bytes"
	"encoding/json"
)

// secretKeys are the JSON keys whose values String() never prints: values the
// API returns once and that grant access on their own (webhook signing
// secrets, API keys, live-view and sign-in tokens). The fields stay readable on
// the struct; only the printed form hides them.
var secretKeys = map[string]bool{
	"secret":        true,
	"key_value":     true,
	"token":         true,
	"live_view_url": true,
}

// StringifyJSON returns a pretty JSON string representation of
// the given value, with secret values replaced by "[redacted]".
func StringifyJSON(value interface{}) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	redacted, err := redactSecrets(raw)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := json.Indent(&out, redacted, "", "  "); err != nil {
		return "", err
	}
	return out.String(), nil
}

// redactSecrets rewrites a JSON document token by token, keeping key order and
// number precision, with the value of every secret key replaced. A null stays
// null, so an absent secret still reads as absent.
func redactSecrets(raw []byte) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var out bytes.Buffer
	if err := copyRedacted(decoder, &out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func copyRedacted(decoder *json.Decoder, out *bytes.Buffer) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return writeJSON(out, token)
	}

	closing := json.Delim(']')
	if delim == '{' {
		closing = '}'
	}
	out.WriteByte(byte(delim))
	for first := true; decoder.More(); first = false {
		if !first {
			out.WriteByte(',')
		}
		if delim == '[' {
			if err := copyRedacted(decoder, out); err != nil {
				return err
			}
			continue
		}

		key, err := decoder.Token()
		if err != nil {
			return err
		}
		if err := writeJSON(out, key); err != nil {
			return err
		}
		out.WriteByte(':')
		if !secretKeys[key.(string)] {
			if err := copyRedacted(decoder, out); err != nil {
				return err
			}
			continue
		}
		var secret json.RawMessage
		if err := decoder.Decode(&secret); err != nil {
			return err
		}
		if string(secret) == "null" {
			out.WriteString("null")
			continue
		}
		out.WriteString(`"[redacted]"`)
	}
	if _, err := decoder.Token(); err != nil { // the closing delimiter
		return err
	}
	out.WriteByte(byte(closing))
	return nil
}

func writeJSON(out *bytes.Buffer, value interface{}) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	out.Write(encoded)
	return nil
}
