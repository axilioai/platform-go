package api

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	core "github.com/axilioai/platform-go/core"
)

func TestIsAccessibilityUnavailable(t *testing.T) {
	conflict := func(body string) error {
		return core.NewAPIError(http.StatusConflict, nil, errors.New(body))
	}
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"allocate conflict", conflict(`{"status":409,"title":"Conflict","detail":"accessibility_unavailable: phone \"p1\" does not support accessibility mode"}`), true},
		{"wrapped", fmt.Errorf("allocate: %w", conflict(`{"detail":"accessibility_unavailable: phone \"p1\" does not support accessibility mode"}`)), true},
		{"other conflict", conflict(`{"status":409,"detail":"session name already in use"}`), false},
		{"same detail, not a conflict", core.NewAPIError(http.StatusBadRequest, nil, errors.New(`{"detail":"accessibility_unavailable: x"}`)), false},
		{"non-JSON body", conflict("accessibility_unavailable: x"), false},
		{"not an API error", errors.New("accessibility_unavailable: x"), false},
		{"nil", nil, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsAccessibilityUnavailable(c.err); got != c.want {
				t.Fatalf("IsAccessibilityUnavailable = %v, want %v", got, c.want)
			}
		})
	}
}
