package changelog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetLatestVersion(t *testing.T) {
	tmpDir := t.TempDir()
	changelogPath := filepath.Join(tmpDir, "CHANGELOG.md")

	content := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.11.0] - 2026-03-09
### Changed
- Updated Lucide icons from 0.575.0 to 0.577.0

## [v0.10.0] - 2026-02-15
### Changed
- Updated Lucide icons from 0.470.0 to 0.475.0

---

[v0.11.0]: https://github.com/kaugesaar/lucide-go/releases/tag/v0.11.0
[v0.10.0]: https://github.com/kaugesaar/lucide-go/releases/tag/v0.10.0
`

	if err := os.WriteFile(changelogPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test changelog: %v", err)
	}

	mgr := New(changelogPath)
	version, err := mgr.GetLatestVersion()
	if err != nil {
		t.Fatalf("GetLatestVersion() error: %v", err)
	}

	if version != "v0.11.0" {
		t.Errorf("GetLatestVersion() = %q, want %q", version, "v0.11.0")
	}
}

func TestGetLatestVersionNoVersions(t *testing.T) {
	tmpDir := t.TempDir()
	changelogPath := filepath.Join(tmpDir, "CHANGELOG.md")

	content := `# Changelog

All notable changes to this project will be documented in this file.
`

	if err := os.WriteFile(changelogPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test changelog: %v", err)
	}

	mgr := New(changelogPath)
	_, err := mgr.GetLatestVersion()
	if err == nil {
		t.Error("GetLatestVersion() should error when no versions exist")
	}
}

func TestGetLatestVersionFileNotFound(t *testing.T) {
	mgr := New("/nonexistent/CHANGELOG.md")
	_, err := mgr.GetLatestVersion()
	if err == nil {
		t.Error("GetLatestVersion() should error when file doesn't exist")
	}
}

func TestGetReleaseNotes(t *testing.T) {
	tmpDir := t.TempDir()
	changelogPath := filepath.Join(tmpDir, "CHANGELOG.md")

	content := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.2.0] - 2025-12-01
### Changed
- Updated Lucide icons from 0.400.0 to 0.450.0
- Added 10 new icon(s)

## [v0.1.0] - 2025-10-26
### Added
- Initial release

---

[v0.2.0]: https://github.com/kaugesaar/lucide-go/releases/tag/v0.2.0
[v0.1.0]: https://github.com/kaugesaar/lucide-go/releases/tag/v0.1.0
`

	if err := os.WriteFile(changelogPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test changelog: %v", err)
	}

	mgr := New(changelogPath)

	notes, err := mgr.GetReleaseNotes("v0.2.0")
	if err != nil {
		t.Fatalf("GetReleaseNotes() error: %v", err)
	}

	if notes == "" {
		t.Fatal("GetReleaseNotes() returned empty string")
	}

	if want := "Updated Lucide icons from 0.400.0 to 0.450.0"; !contains(notes, want) {
		t.Errorf("GetReleaseNotes() missing %q", want)
	}

	if want := "Added 10 new icon(s)"; !contains(notes, want) {
		t.Errorf("GetReleaseNotes() missing %q", want)
	}
}

func TestGetReleaseNotesVersionNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	changelogPath := filepath.Join(tmpDir, "CHANGELOG.md")

	content := `# Changelog

## [v0.1.0] - 2025-10-26
### Added
- Initial release
`

	if err := os.WriteFile(changelogPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test changelog: %v", err)
	}

	mgr := New(changelogPath)
	_, err := mgr.GetReleaseNotes("v0.99.0")
	if err == nil {
		t.Error("GetReleaseNotes() should error for nonexistent version")
	}
}

func TestNextVersionFromChangelog(t *testing.T) {
	tests := []struct {
		name             string
		changelogVersion string
		wantVersion      string
	}{
		{
			name:             "single digit minor",
			changelogVersion: "v0.5.0",
			wantVersion:      "v0.6.0",
		},
		{
			name:             "double digit minor",
			changelogVersion: "v0.11.0",
			wantVersion:      "v0.12.0",
		},
		{
			name:             "high minor version",
			changelogVersion: "v1.99.0",
			wantVersion:      "v1.100.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			changelogPath := filepath.Join(tmpDir, "CHANGELOG.md")

			content := "# Changelog\n\nAll notable changes.\n\nFormat based on Keep a Changelog.\n\n## [" + tt.changelogVersion + "] - 2026-01-01\n### Changed\n- Update\n"
			if err := os.WriteFile(changelogPath, []byte(content), 0o644); err != nil {
				t.Fatalf("failed to write changelog: %v", err)
			}

			mgr := New(changelogPath)
			version, err := mgr.GetLatestVersion()
			if err != nil {
				t.Fatalf("GetLatestVersion() error: %v", err)
			}

			if version != tt.changelogVersion {
				t.Fatalf("GetLatestVersion() = %q, want %q", version, tt.changelogVersion)
			}

			got := bumpMinor(version)
			if got != tt.wantVersion {
				t.Errorf("bumpMinor(%q) = %q, want %q", version, got, tt.wantVersion)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
