package changelog

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// GetNextVersion determines the next minor version by reading git tags.
// Returns the next version in format "v0.5.0".
func GetNextVersion() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "v0.1.0", nil
	}

	currentVersion := strings.TrimSpace(string(output))

	re := regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)$`)
	matches := re.FindStringSubmatch(currentVersion)
	if len(matches) != 4 {
		return "", fmt.Errorf("invalid version format: %s", currentVersion)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])

	nextVersion := fmt.Sprintf("v%d.%d.0", major, minor+1)
	return nextVersion, nil
}

// GetLatestVersion parses the topmost version from the changelog.
// Returns the version string (e.g., "v0.5.0") or empty string if not found.
func (m *Manager) GetLatestVersion() (string, error) {
	content, err := m.readContent()
	if err != nil {
		return "", err
	}

	lines := strings.Split(content, "\n")

	re := regexp.MustCompile(`^## \[(v\d+\.\d+\.\d+)\]`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 2 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("no version found in changelog")
}

// GetReleaseNotes extracts the changelog content for a specific version.
// Returns the release notes as a string.
func (m *Manager) GetReleaseNotes(version string) (string, error) {
	content, err := m.readContent()
	if err != nil {
		return "", err
	}

	lines := strings.Split(content, "\n")

	versionHeader := fmt.Sprintf("## [%s]", version)
	startIdx := -1
	for i, line := range lines {
		if strings.HasPrefix(line, versionHeader) {
			startIdx = i
			break
		}
	}

	if startIdx == -1 {
		return "", fmt.Errorf("version %s not found in changelog", version)
	}

	var notes []string
	for i := startIdx + 1; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, "## [") || strings.HasPrefix(line, "---") {
			break
		}

		notes = append(notes, line)
	}

	for len(notes) > 0 && strings.TrimSpace(notes[len(notes)-1]) == "" {
		notes = notes[:len(notes)-1]
	}

	return strings.Join(notes, "\n"), nil
}

// readContent is a helper to read the changelog file.
func (m *Manager) readContent() (string, error) {
	content, err := os.ReadFile(m.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read changelog: %w", err)
	}
	return string(content), nil
}
