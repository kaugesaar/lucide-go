package lucide

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v78/github"
)

// DownloadAndExtract downloads a release asset and extracts it to the destination directory.
// It expects the asset to be a zip file containing an icons/ directory with SVG files.
func DownloadAndExtract(ctx context.Context, client *Client, asset *github.ReleaseAsset, destDir string) error {
	downloadURL := asset.GetBrowserDownloadURL()
	if downloadURL == "" {
		return fmt.Errorf("asset has no download URL")
	}

	tmpFile, err := os.CreateTemp("", "lucide-icons-*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) //nolint:errcheck // Best-effort cleanup

	if err := downloadFile(ctx, downloadURL, tmpFile); err != nil {
		tmpFile.Close() //nolint:errcheck // Cleanup on error path
		return fmt.Errorf("failed to download asset: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := extractIcons(tmpFile.Name(), destDir); err != nil {
		return fmt.Errorf("failed to extract icons: %w", err)
	}

	return nil
}

// DownloadAndExtractTarball downloads a source tarball and extracts icon files from it.
// It expects the tarball to contain an icons/ directory with SVG and JSON files.
func DownloadAndExtractTarball(ctx context.Context, archiveURL, destDir string) error {
	tmpFile, err := os.CreateTemp("", "lucide-source-*.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) //nolint:errcheck // Best-effort cleanup

	if err := downloadFile(ctx, archiveURL, tmpFile); err != nil {
		tmpFile.Close() //nolint:errcheck // Cleanup on error path
		return fmt.Errorf("failed to download tarball: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := extractIconsFromTarball(tmpFile.Name(), destDir); err != nil {
		return fmt.Errorf("failed to extract icons from tarball: %w", err)
	}

	return nil
}

func extractIconsFromTarball(tarballPath, destDir string) error {
	f, err := os.Open(tarballPath)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck // File cleanup

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close() //nolint:errcheck // Gzip reader cleanup

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}

	if err := clearDirectory(destDir); err != nil {
		return fmt.Errorf("failed to clear destination directory: %w", err)
	}

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		parts := strings.Split(header.Name, "/")
		if len(parts) < 3 || parts[1] != "icons" {
			continue
		}

		filename := parts[len(parts)-1]
		if !strings.HasSuffix(filename, ".svg") && !strings.HasSuffix(filename, ".json") {
			continue
		}

		destPath := filepath.Join(destDir, filename)
		if err := writeFile(destPath, tr); err != nil {
			return fmt.Errorf("failed to extract %s: %w", header.Name, err)
		}
	}

	return nil
}

func writeFile(destPath string, r io.Reader) error {
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close() //nolint:errcheck // File will be synced by Copy

	_, err = io.Copy(dest, r)
	return err
}

func downloadFile(ctx context.Context, url string, dest *os.File) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck // Response body cleanup

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	_, err = io.Copy(dest, resp.Body)
	return err
}

func extractIcons(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close() //nolint:errcheck // Zip reader cleanup

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}

	if err := clearDirectory(destDir); err != nil {
		return fmt.Errorf("failed to clear destination directory: %w", err)
	}

	for _, f := range r.File {
		if !strings.HasPrefix(f.Name, "icons/") {
			continue
		}

		if !strings.HasSuffix(f.Name, ".svg") && !strings.HasSuffix(f.Name, ".json") {
			continue
		}

		filename := filepath.Base(f.Name)
		destPath := filepath.Join(destDir, filename)

		if err := extractFile(f, destPath); err != nil {
			return fmt.Errorf("failed to extract %s: %w", f.Name, err)
		}
	}

	return nil
}

func extractFile(f *zip.File, destPath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close() //nolint:errcheck // Read closer cleanup

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close() //nolint:errcheck // File will be synced by Copy

	_, err = io.Copy(dest, rc)
	return err
}

func clearDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}
