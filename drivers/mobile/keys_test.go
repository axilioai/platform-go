package mobile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// TestNamedKeysMatchContract pins the hand-written key constants to the
// contract's enum (AXI-2025): the wire generator does not emit enum values,
// so a name the phone learns would otherwise reach Go users only as a raw
// string. A key present on one side and not the other fails here.
func TestNamedKeysMatchContract(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "contracts", "dcp-asyncapi.json"))
	if err != nil {
		t.Fatal(err)
	}
	var contract struct {
		Components struct {
			Schemas struct {
				KeyboardKeyPressParams struct {
					Properties struct {
						Key struct {
							Enum []string `json:"enum"`
						} `json:"key"`
					} `json:"properties"`
				} `json:"KeyboardKeyPressParams"`
			} `json:"schemas"`
		} `json:"components"`
	}
	if err := json.Unmarshal(raw, &contract); err != nil {
		t.Fatal(err)
	}
	want := slices.Clone(contract.Components.Schemas.KeyboardKeyPressParams.Properties.Key.Enum)
	if len(want) == 0 {
		t.Fatal("contract has no KeyboardKeyPressParams.key enum")
	}
	got := []string{KeyEnter, KeyCapsLock}
	slices.Sort(want)
	slices.Sort(got)
	if !slices.Equal(got, want) {
		t.Fatalf("named keys: constants %v, contract enum %v", got, want)
	}
}
