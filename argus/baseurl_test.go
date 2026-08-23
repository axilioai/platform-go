package api

import "testing"

// The synced argus spec carries no servers array; scripts/regen_argus.sh
// injects the production host at generation time so the generated clients
// default to it (matching platform-python's hardcoded wrapper default). This
// pins that injection: if a regen bypasses the script or the injection stops
// working, the default silently becomes "" and every call needs an explicit
// WithBaseURL. The file survives regens via the script's rsync exclude.
func TestDefaultEnvironmentIsProductionArgusHost(t *testing.T) {
	if got, want := Environments.Default, "https://argus.axilio.ai"; got != want {
		t.Fatalf("Environments.Default = %q, want %q (regen_argus.sh servers injection lost?)", got, want)
	}
}
