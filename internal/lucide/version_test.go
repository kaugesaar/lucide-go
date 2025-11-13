package lucide

import (
	"os"
	"testing"
)

func TestGetCurrentVersion(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("failed to restore directory: %v", err)
		}
	}()

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}

	_, err := GetCurrentVersion()
	if err == nil {
		t.Error("GetCurrentVersion() should return error when file doesn't exist")
	}

	if err := os.WriteFile(".lucide-version", []byte("0.553.0\n"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	version, err := GetCurrentVersion()
	if err != nil {
		t.Errorf("GetCurrentVersion() failed: %v", err)
	}
	if version != "0.553.0" {
		t.Errorf("GetCurrentVersion() = %q, want %q", version, "0.553.0")
	}

	if err := os.WriteFile(".lucide-version", []byte(""), 0o644); err != nil {
		t.Fatalf("failed to create empty test file: %v", err)
	}

	_, err = GetCurrentVersion()
	if err == nil {
		t.Error("GetCurrentVersion() should return error for empty file")
	}

	if err := os.WriteFile(".lucide-version", []byte("  \n  \n"), 0o644); err != nil {
		t.Fatalf("failed to create whitespace test file: %v", err)
	}

	_, err = GetCurrentVersion()
	if err == nil {
		t.Error("GetCurrentVersion() should return error for whitespace-only file")
	}
}

func TestSetCurrentVersion(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("failed to restore directory: %v", err)
		}
	}()

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}

	err := SetCurrentVersion("1.2.3")
	if err != nil {
		t.Errorf("SetCurrentVersion() failed: %v", err)
	}

	content, err := os.ReadFile(".lucide-version")
	if err != nil {
		t.Errorf("failed to read version file: %v", err)
	}

	expected := "1.2.3\n"
	if string(content) != expected {
		t.Errorf("version file content = %q, want %q", string(content), expected)
	}

	err = SetCurrentVersion("2.0.0")
	if err != nil {
		t.Errorf("SetCurrentVersion() failed on update: %v", err)
	}

	content, _ = os.ReadFile(".lucide-version")
	expected = "2.0.0\n"
	if string(content) != expected {
		t.Errorf("updated version file content = %q, want %q", string(content), expected)
	}
}

func TestVersionRoundTrip(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("failed to restore directory: %v", err)
		}
	}()

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}

	versions := []string{"0.1.0", "1.2.3", "10.20.30"}

	for _, want := range versions {
		if err := SetCurrentVersion(want); err != nil {
			t.Errorf("SetCurrentVersion(%q) failed: %v", want, err)
		}

		got, err := GetCurrentVersion()
		if err != nil {
			t.Errorf("GetCurrentVersion() after setting %q failed: %v", want, err)
		}

		if got != want {
			t.Errorf("round trip failed: set %q, got %q", want, got)
		}
	}
}
