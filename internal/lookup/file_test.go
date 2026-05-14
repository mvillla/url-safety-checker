package lookup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadURLsFileIgnoresCommentsAndEmptyLines(t *testing.T) {
	path := writeTestFile(t, `
# comment

malware.test/bad
  
bad.example/path
`)

	got, err := LoadURLsFile(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := []string{"malware.test/bad", "bad.example/path"}
	assertStringSlicesEqual(t, got, want)
}

func TestLoadURLsFileReturnsErrorForMissingFile(t *testing.T) {
	_, err := LoadURLsFile(filepath.Join(t.TempDir(), "missing.txt"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func writeTestFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "urls.txt")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	return path
}

func assertStringSlicesEqual(t *testing.T, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expected %d items, got %d: %#v", len(want), len(got), got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected item %d to be %q, got %q", i, want[i], got[i])
		}
	}
}
