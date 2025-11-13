package changelog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	mgr := New("")
	if mgr.Path != defaultChangelog {
		t.Errorf("New(\"\").Path = %q, want %q", mgr.Path, defaultChangelog)
	}

	customPath := "custom-changelog.md"
	mgr = New(customPath)
	if mgr.Path != customPath {
		t.Errorf("New(%q).Path = %q, want %q", customPath, mgr.Path, customPath)
	}
}

func TestAddEntry(t *testing.T) {
	tmpDir := t.TempDir()
	changelogPath := filepath.Join(tmpDir, "CHANGELOG.md")

	initialContent := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.0] - 2025-10-26
### Added
- Initial release
`

	if err := os.WriteFile(changelogPath, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("failed to create test changelog: %v", err)
	}

	mgr := New(changelogPath)

	entry := Entry{
		Date:         time.Date(2025, 11, 13, 0, 0, 0, 0, time.UTC),
		CurrentTag:   "0.1.0",
		NewTag:       "0.2.0",
		IconsAdded:   5,
		IconsRemoved: 2,
	}

	err := mgr.AddEntry(entry)
	if err != nil {
		t.Fatalf("AddEntry() failed: %v", err)
	}

	content, err := os.ReadFile(changelogPath)
	if err != nil {
		t.Fatalf("failed to read changelog: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "## [Unreleased] - 2025-11-13") {
		t.Error("changelog doesn't contain expected date header")
	}

	if !strings.Contains(contentStr, "Updated Lucide icons from 0.1.0 to 0.2.0") {
		t.Error("changelog doesn't contain update message")
	}

	if !strings.Contains(contentStr, "Added 5 new icon(s)") {
		t.Error("changelog doesn't contain icons added message")
	}

	if !strings.Contains(contentStr, "Removed 2 icon(s)") {
		t.Error("changelog doesn't contain icons removed message")
	}

	if !strings.Contains(contentStr, "## [v0.1.0] - 2025-10-26") {
		t.Error("changelog lost old entry")
	}

	unreleasedIdx := strings.Index(contentStr, "## [Unreleased]")
	oldEntryIdx := strings.Index(contentStr, "## [v0.1.0]")
	if unreleasedIdx > oldEntryIdx {
		t.Error("new entry should come before old entry")
	}
}

func TestAddEntryNoIconChanges(t *testing.T) {
	tmpDir := t.TempDir()
	changelogPath := filepath.Join(tmpDir, "CHANGELOG.md")

	initialContent := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.0] - 2025-10-26
### Added
- Initial release
`

	if err := os.WriteFile(changelogPath, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("failed to create test changelog: %v", err)
	}

	mgr := New(changelogPath)

	entry := Entry{
		Date:         time.Date(2025, 11, 13, 0, 0, 0, 0, time.UTC),
		CurrentTag:   "0.1.0",
		NewTag:       "0.2.0",
		IconsAdded:   0,
		IconsRemoved: 0,
	}

	err := mgr.AddEntry(entry)
	if err != nil {
		t.Fatalf("AddEntry() failed: %v", err)
	}

	content, _ := os.ReadFile(changelogPath)
	contentStr := string(content)

	if strings.Contains(contentStr, "Added 0") {
		t.Error("changelog should not mention adding 0 icons")
	}

	if strings.Contains(contentStr, "Removed 0") {
		t.Error("changelog should not mention removing 0 icons")
	}

	if !strings.Contains(contentStr, "Updated Lucide icons from 0.1.0 to 0.2.0") {
		t.Error("changelog should contain update message even with no icon changes")
	}
}

func TestAddEntryInvalidChangelog(t *testing.T) {
	tmpDir := t.TempDir()
	changelogPath := filepath.Join(tmpDir, "CHANGELOG.md")

	shortContent := `# Changelog

Too short
`

	if err := os.WriteFile(changelogPath, []byte(shortContent), 0o644); err != nil {
		t.Fatalf("failed to create test changelog: %v", err)
	}

	mgr := New(changelogPath)

	entry := Entry{
		Date:       time.Now(),
		CurrentTag: "0.1.0",
		NewTag:     "0.2.0",
	}

	err := mgr.AddEntry(entry)
	if err == nil {
		t.Error("AddEntry() should fail with too few lines")
	}
}

func TestAddEntryFileNotFound(t *testing.T) {
	mgr := New("/nonexistent/changelog.md")

	entry := Entry{
		Date:       time.Now(),
		CurrentTag: "0.1.0",
		NewTag:     "0.2.0",
	}

	err := mgr.AddEntry(entry)
	if err == nil {
		t.Error("AddEntry() should fail when file doesn't exist")
	}
}

func TestFormatEntry(t *testing.T) {
	entry := Entry{
		Date:         time.Date(2025, 11, 13, 0, 0, 0, 0, time.UTC),
		CurrentTag:   "0.1.0",
		NewTag:       "0.2.0",
		IconsAdded:   10,
		IconsRemoved: 5,
	}

	result := formatEntry(entry)

	if !strings.Contains(result, "## [Unreleased] - 2025-11-13") {
		t.Error("formatEntry() missing header")
	}

	if !strings.Contains(result, "### Changed") {
		t.Error("formatEntry() missing Changed section")
	}

	if !strings.Contains(result, "Updated Lucide icons from 0.1.0 to 0.2.0") {
		t.Error("formatEntry() missing update line")
	}

	if !strings.Contains(result, "Added 10 new icon(s)") {
		t.Error("formatEntry() missing added line")
	}

	if !strings.Contains(result, "Removed 5 icon(s)") {
		t.Error("formatEntry() missing removed line")
	}

	if !strings.HasSuffix(result, "\n\n") {
		t.Error("formatEntry() should end with double newline")
	}
}

func TestFormatEntryNoChanges(t *testing.T) {
	entry := Entry{
		Date:         time.Date(2025, 11, 13, 0, 0, 0, 0, time.UTC),
		CurrentTag:   "0.1.0",
		NewTag:       "0.2.0",
		IconsAdded:   0,
		IconsRemoved: 0,
	}

	result := formatEntry(entry)

	if strings.Contains(result, "Added 0") || strings.Contains(result, "Removed 0") {
		t.Error("formatEntry() should not include zero counts")
	}

	if !strings.Contains(result, "Updated Lucide icons from 0.1.0 to 0.2.0") {
		t.Error("formatEntry() missing update line")
	}
}
