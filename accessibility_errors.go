package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	core "github.com/axilioai/platform-go/core"
)

// AccessibilityUnavailableKind is the stable token that leads the detail of
// the conflict Phones.Allocate returns when the request sets Accessibility
// to true and names a phone (PhoneID) that does not support accessibility
// mode. It is also the session end reason when the phone cannot confirm
// accessibility mode before the session goes live.
const AccessibilityUnavailableKind = "accessibility_unavailable"

// IsAccessibilityUnavailable reports whether err is the allocate conflict
// for a phone that does not support accessibility mode: an HTTP 409 whose
// problem detail starts with AccessibilityUnavailableKind. Retrying the same
// request cannot succeed; drop PhoneID (to claim any capable phone), or
// omit Accessibility to get it only where the phone supports it.
//
// This file is hand-written (scripts/regen.sh keeps it across regens).
func IsAccessibilityUnavailable(err error) bool {
	var apiErr *core.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusConflict {
		return false
	}
	inner := apiErr.Unwrap()
	if inner == nil {
		return false
	}
	var body struct {
		Detail string `json:"detail"`
	}
	if json.Unmarshal([]byte(inner.Error()), &body) != nil {
		return false
	}
	return strings.HasPrefix(body.Detail, AccessibilityUnavailableKind+":")
}
