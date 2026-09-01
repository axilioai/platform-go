package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFailureEvidenceRedactsErrorDetails(t *testing.T) {
	got := redactedFailure(errors.New("GET /sessions/secret-id Authorization: secret-token"))
	if strings.Contains(got, "secret-id") || strings.Contains(got, "secret-token") {
		t.Fatalf("redacted failure leaked details: %s", got)
	}
	if !strings.Contains(got, "error_type=") || !strings.Contains(got, "fingerprint=") {
		t.Fatalf("redacted failure lacks safe diagnostics: %s", got)
	}
	if strings.Contains(got, "details") {
		t.Fatalf("redacted failure contains detail marker: %s", got)
	}
}

func TestReplayMatrixPasses(t *testing.T) {
	rec := new(recorder)
	runReplay(context.Background(), rec)
	if !rec.passed() {
		for _, result := range rec.Results {
			if result.Status != "PASS" {
				t.Errorf("%s: %s", result.ID, result.ObservedRedacted)
			}
		}
	}
	want := []string{
		"GO-REPLAY-01", "GO-REPLAY-02", "GO-REPLAY-03", "GO-REPLAY-04",
		"GO-REPLAY-05", "GO-REPLAY-06", "GO-REPLAY-07",
	}
	if len(rec.Results) != len(want) {
		t.Fatalf("checks = %d, want %d", len(rec.Results), len(want))
	}
	for i := range want {
		if rec.Results[i].ID != want[i] {
			t.Errorf("check[%d] = %s, want %s", i, rec.Results[i].ID, want[i])
		}
	}
}

func TestValidateTarget(t *testing.T) {
	tests := []struct {
		name, environment, baseURL string
		wantErr                    bool
	}{
		{"staging exact", "staging", "https://staging-api.axilio.ai", false},
		{"dev loopback", "dev", "http://127.0.0.1:8080", false},
		{"production origin", "staging", "https://api.axilio.ai", true},
		{"production environment", "production", "https://api.axilio.ai", true},
		{"arbitrary staging", "staging", "https://example.invalid", true},
		{"remote cleartext dev", "dev", "http://example.invalid", true},
		{"path", "dev", "https://dev.example.invalid/api/v1", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateTarget(tt.environment, tt.baseURL, "")
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTargetRefusesTrailingDotProductionAndArbitraryHTTPSDev(t *testing.T) {
	for _, target := range []string{"https://api.axilio.ai.", "https://example.invalid"} {
		if _, err := validateTarget("dev", target, ""); err == nil {
			t.Fatalf("validateTarget(dev, %q) succeeded", target)
		}
	}
}

func TestValidateTargetAllowsOnlyApprovedRemoteDevOrigin(t *testing.T) {
	approved := "https://dev-api.example.invalid"
	if got, err := validateTarget("dev", approved, approved); err != nil || got != approved {
		t.Fatalf("validateTarget approved dev = %q, %v", got, err)
	}
	if _, err := validateTarget("dev", approved, "https://other-dev.example.invalid"); err == nil {
		t.Fatal("validateTarget accepted mismatched dev origin")
	}
}

func TestOutputIncludesFixtureAndFileProvenance(t *testing.T) {
	tempDir := t.TempDir()
	manifestPath := filepath.Join(tempDir, "fixtures.json")
	manifestBody := []byte(`{"manifest_version":1,"environment":"dev","seed_revision":"seed-sha","fixtures":{"normal_session":{"id":"normal"},"normal_empty_session":{"id":"empty"},"expired_session":{"id":"expired"}}}`)
	if err := os.WriteFile(manifestPath, manifestBody, 0o600); err != nil {
		t.Fatal(err)
	}
	outputPath := filepath.Join(tempDir, "go.json")
	if err := execute(context.Background(), "dev", "http://127.0.0.1", "", manifestPath, outputPath, "candidate-sha", "artifact-sha", true, ""); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	var got output
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(manifestBody)
	if got.SeedRevision != "seed-sha" || got.FixtureManifestSHA256 != fmt.Sprintf("%x", digest) {
		t.Fatalf("fixture provenance = seed %q hash %q", got.SeedRevision, got.FixtureManifestSHA256)
	}
	for _, result := range got.Results {
		if result.EvidenceFile != filepath.Base(outputPath) {
			t.Fatalf("%s evidence_file = %q", result.ID, result.EvidenceFile)
		}
	}
}
