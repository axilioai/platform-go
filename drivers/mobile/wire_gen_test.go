package mobile

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// TestWireGenUpToDate is the SDK-side end of the DCP honesty chain: the backend
// drift gate forces the contract to equal the Go server, and this forces the
// committed wire_gen.go to equal a fresh generation from the vendored contract.
// Together they guarantee the SDK's typed layer equals what the server speaks.
// It regenerates into a temp file and byte-compares, so any drift in the
// methods, error kinds, protocol version, or param frames fails here.
func TestWireGenUpToDate(t *testing.T) {
	root := repoRoot(t)
	// Not a .go name: `go run gen.go x.go` would treat x.go as a second source.
	tmp := filepath.Join(t.TempDir(), "wire_gen.generated")

	cmd := exec.Command("go", "run", "scripts/gen_dcp_wire.go", tmp)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generator failed: %v\n%s", err, out)
	}

	want, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("read regenerated: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(root, "drivers", "mobile", "wire_gen.go"))
	if err != nil {
		t.Fatalf("read committed: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("drivers/mobile/wire_gen.go is out of date with contracts/dcp-asyncapi.json; " +
			"run: go run scripts/gen_dcp_wire.go")
	}
}

// repoRoot returns the module root (two levels up from this test file).
func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "..", "..")
}
