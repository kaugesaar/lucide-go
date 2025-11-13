package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/kaugesaar/lucide-go/internal/changelog"
	"github.com/kaugesaar/lucide-go/internal/generator"
	"github.com/kaugesaar/lucide-go/internal/lucide"
)

const (
	iconsDir   = "lucide-icons"
	outputFile = "icons.go"
)

type UpdateResult struct {
	HasUpdates    bool   `json:"has_updates"`
	CurrentTag    string `json:"current_tag"`
	LatestTag     string `json:"latest_tag"`
	IconsAdded    int    `json:"icons_added"`
	IconsRemoved  int    `json:"icons_removed"`
	ReleaseURL    string `json:"release_url"`
	ReleaseNotes  string `json:"release_notes,omitempty"`
	ChangelogPath string `json:"changelog_path,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "check":
		if err := runCheck(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "update":
		if err := runUpdate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "generate":
		if err := runGenerate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`lucide-go tool

Usage:
  tool check      Check if updates are available
  tool update     Download latest release and regenerate icons
  tool generate   Regenerate icons from current icon files
  tool help       Show this help message

Commands:
  check          Checks GitHub for latest Lucide release and compares
                 with current version. Outputs JSON result.

  update         Downloads latest Lucide release if newer than current,
                 regenerates icons, updates changelog, and updates version file.
                 Outputs JSON result for CI consumption.

  generate       Regenerates icons.go from icon files in the lucide-icons
                 directory without downloading or updating anything.

Global Flags:
  --dry-run      Preview changes without writing (update command only)
`)
}

func runCheck() error {
	ctx := context.Background()

	currentTag, err := lucide.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Current version: %s\n", currentTag)

	client := lucide.NewClient(os.Getenv("GITHUB_TOKEN"))
	release, err := client.GetLatestRelease(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Latest version: %s\n", release.TagName)

	result := UpdateResult{
		HasUpdates: currentTag != release.TagName,
		CurrentTag: currentTag,
		LatestTag:  release.TagName,
		ReleaseURL: release.URL,
	}

	if result.HasUpdates {
		fmt.Fprintf(os.Stderr, "✓ Update available: %s → %s\n", currentTag, release.TagName)
	} else {
		fmt.Fprintf(os.Stderr, "✓ Already up to date\n")
	}

	return outputJSON(result)
}

func runUpdate() error {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	dryRun := fs.Bool("dry-run", false, "Preview changes without writing")
	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	ctx := context.Background()

	currentTag, err := lucide.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Current version: %s\n", currentTag)

	client := lucide.NewClient(os.Getenv("GITHUB_TOKEN"))
	release, err := client.GetLatestRelease(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Latest version: %s\n", release.TagName)

	result := UpdateResult{
		CurrentTag:   currentTag,
		LatestTag:    release.TagName,
		ReleaseURL:   release.URL,
		ReleaseNotes: release.Body,
	}

	if currentTag == release.TagName {
		fmt.Fprintf(os.Stderr, "✓ Already up to date\n")
		result.HasUpdates = false
		return outputJSON(result)
	}

	result.HasUpdates = true
	fmt.Fprintf(os.Stderr, "\nUpdating from %s to %s...\n", currentTag, release.TagName)

	if *dryRun {
		fmt.Fprintf(os.Stderr, "DRY RUN: Would download and update to %s\n", release.TagName)
		return outputJSON(result)
	}

	asset, err := release.FindIconsAsset()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Downloading %s...\n", asset.GetName())

	if err := lucide.DownloadAndExtract(ctx, client, asset, iconsDir); err != nil {
		return fmt.Errorf("failed to download icons: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Regenerating icons...\n")
	gen := generator.New(iconsDir, outputFile)
	genResult, err := gen.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate icons: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Generated %d icons\n", genResult.IconsGenerated)

	added, removed, err := countIconChanges()
	if err != nil {
		return fmt.Errorf("failed to count changes: %w", err)
	}
	result.IconsAdded = added
	result.IconsRemoved = removed
	fmt.Fprintf(os.Stderr, "Icons added: %d, removed: %d\n", added, removed)

	changelogMgr := changelog.New("CHANGELOG.md")
	entry := changelog.Entry{
		Date:         time.Now(),
		CurrentTag:   currentTag,
		NewTag:       release.TagName,
		IconsAdded:   added,
		IconsRemoved: removed,
	}
	if err := changelogMgr.AddEntry(entry); err != nil {
		return fmt.Errorf("failed to update changelog: %w", err)
	}
	result.ChangelogPath = "CHANGELOG.md"
	fmt.Fprintf(os.Stderr, "Updated CHANGELOG.md\n")

	if err := lucide.SetCurrentVersion(release.TagName); err != nil {
		return fmt.Errorf("failed to update version file: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Updated .lucide-version\n")

	fmt.Fprintf(os.Stderr, "\n✓ Update complete!\n")
	return outputJSON(result)
}

func runGenerate() error {
	fmt.Fprintf(os.Stderr, "Regenerating icons from %s...\n", iconsDir)

	gen := generator.New(iconsDir, outputFile)
	result, err := gen.Generate()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "✓ Successfully generated %d icons to %s\n", result.IconsGenerated, outputFile)
	return nil
}

func countIconChanges() (added, removed int, err error) {
	cmd := exec.Command("git", "diff", "--name-only", outputFile)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("git diff failed: %w", err)
	}

	if strings.TrimSpace(string(output)) == "" {
		return 0, 0, nil
	}

	cmd = exec.Command("git", "diff", outputFile)
	output, err = cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("git diff failed: %w", err)
	}

	diff := string(output)
	addedRe := regexp.MustCompile(`(?m)^\+func [A-Z]`)
	removedRe := regexp.MustCompile(`(?m)^-func [A-Z]`)

	added = len(addedRe.FindAllString(diff, -1))
	removed = len(removedRe.FindAllString(diff, -1))

	return added, removed, nil
}

func outputJSON(result any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}
