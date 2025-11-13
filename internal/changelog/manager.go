package changelog

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const defaultChangelog = "CHANGELOG.md"

// Entry represents a single changelog entry.
type Entry struct {
	Version      string
	Date         time.Time
	CurrentTag   string
	NewTag       string
	IconsAdded   int
	IconsRemoved int
}

// Manager handles reading and updating the changelog file.
type Manager struct {
	Path string
}

// New creates a new changelog manager for the given file path.
func New(path string) *Manager {
	if path == "" {
		path = defaultChangelog
	}
	return &Manager{Path: path}
}

// AddEntry adds a new unreleased entry to the changelog.
func (m *Manager) AddEntry(entry Entry) error {
	content, err := os.ReadFile(m.Path)
	if err != nil {
		return fmt.Errorf("failed to read changelog: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Find the line after the header (line 6 in the current format)
	// Format:
	// 0: # Changelog
	// 1: (blank)
	// 2: All notable changes...
	// 3: (blank)
	// 4: The format is based...
	// 5: (blank)
	// 6: ## [v0.3.0] - date  <- Insert before this
	insertIndex := 6
	if len(lines) < insertIndex {
		return fmt.Errorf("changelog format unexpected: too few lines")
	}

	entryText := formatEntry(entry)

	entryLines := strings.Split(entryText, "\n")

	newLines := make([]string, 0, len(lines)+len(entryLines))
	newLines = append(newLines, lines[:insertIndex]...)
	newLines = append(newLines, entryLines...)
	newLines = append(newLines, lines[insertIndex:]...)

	newContent := strings.Join(newLines, "\n")

	if err := os.WriteFile(m.Path, []byte(newContent), 0o644); err != nil {
		return fmt.Errorf("failed to write changelog: %w", err)
	}

	return nil
}

func formatEntry(e Entry) string {
	dateStr := e.Date.Format("2006-01-02")

	var b strings.Builder
	b.WriteString(fmt.Sprintf("## [Unreleased] - %s\n", dateStr))
	b.WriteString("### Changed\n")
	b.WriteString(fmt.Sprintf("- Updated Lucide icons from %s to %s\n", e.CurrentTag, e.NewTag))

	if e.IconsAdded > 0 {
		b.WriteString(fmt.Sprintf("- Added %d new icon(s)\n", e.IconsAdded))
	}

	if e.IconsRemoved > 0 {
		b.WriteString(fmt.Sprintf("- Removed %d icon(s)\n", e.IconsRemoved))
	}

	b.WriteString("\n")

	return b.String()
}
