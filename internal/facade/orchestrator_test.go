package facade

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_MissingDirectory(t *testing.T) {
	err := Run("/path/to/nonexistent/dir/12345", 1)
	if err == nil {
		t.Fatalf("Expected error for missing directory, got nil")
	}
}

func TestRun_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	err := Run(tempDir, 1)
	if err != nil {
		t.Fatalf("Expected no error for empty directory, got %v", err)
	}
}

func TestRun_Success(t *testing.T) {
	tempDir := t.TempDir()
	txtPath := filepath.Join(tempDir, "test.txt")
	os.WriteFile(txtPath, []byte("listen silent"), 0644)

	err := Run(tempDir, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
