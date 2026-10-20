package api

import "testing"

// The synced spec's only server is the relative "/api/v1"; scripts/regen.sh
// swaps in the production host at generation time so the generated clients
// default to it (matching platform-python's Client default). This pins that
// swap: if a regen bypasses the script or the swap stops working, the default
// goes back to "/api/v1" and every call needs an explicit WithBaseURL. The
// file survives regens via the script's rsync exclude.
func TestDefaultEnvironmentIsProductionAPIHost(t *testing.T) {
	if got, want := Environments.Default, "https://api.axilio.ai/api/v1"; got != want {
		t.Fatalf("Environments.Default = %q, want %q (regen.sh servers swap lost?)", got, want)
	}
}
