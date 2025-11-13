package lucide

import (
	"fmt"
	"os"
	"strings"
)

const versionFile = ".lucide-version"

// GetCurrentVersion reads the current Lucide version from the version file.
func GetCurrentVersion() (string, error) {
	content, err := os.ReadFile(versionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("version file not found: %s", versionFile)
		}
		return "", fmt.Errorf("failed to read version file: %w", err)
	}

	version := strings.TrimSpace(string(content))
	if version == "" {
		return "", fmt.Errorf("version file is empty")
	}

	return version, nil
}

// SetCurrentVersion writes the given version to the version file.
func SetCurrentVersion(version string) error {
	content := version + "\n"
	if err := os.WriteFile(versionFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write version file: %w", err)
	}
	return nil
}
