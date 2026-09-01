package main

import (
	"context"
	"testing"
)

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
			_, err := validateTarget(tt.environment, tt.baseURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
