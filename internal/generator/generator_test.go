package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple two parts",
			input: "circle-x",
			want:  "CircleX",
		},
		{
			name:  "three parts",
			input: "chevron-down",
			want:  "ChevronDown",
		},
		{
			name:  "single letter prefix",
			input: "a-arrow-down",
			want:  "AArrowDown",
		},
		{
			name:  "single letter",
			input: "x",
			want:  "X",
		},
		{
			name:  "multiple dashes",
			input: "align-horizontal-space-between",
			want:  "AlignHorizontalSpaceBetween",
		},
		{
			name:  "numbers",
			input: "number-1",
			want:  "Number1",
		},
		{
			name:  "already capitalized",
			input: "Circle-X",
			want:  "CircleX",
		},
		{
			name:  "single word",
			input: "menu",
			want:  "Menu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toPascalCase(tt.input)
			if got != tt.want {
				t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractSVGContent(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "simple circle",
			input: `<svg
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <circle cx="12" cy="12" r="10" />
</svg>`,
			want: `<circle cx="12" cy="12" r="10" />`,
		},
		{
			name: "multiple paths",
			input: `<svg
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <path d="m15 9-6 6" />
  <path d="m9 9 6 6" />
</svg>`,
			want: `<path d="m15 9-6 6" /> <path d="m9 9 6 6" />`,
		},
		{
			name: "mixed elements",
			input: `<svg
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <circle cx="12" cy="12" r="10" />
  <path d="m15 9-6 6" />
  <path d="m9 9 6 6" />
</svg>`,
			want: `<circle cx="12" cy="12" r="10" /> <path d="m15 9-6 6" /> <path d="m9 9 6 6" />`,
		},
		{
			name:  "empty svg",
			input: `<svg></svg>`,
			want:  ``,
		},
		{
			name: "self-closing tag",
			input: `<svg
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <rect width="18" height="18" x="3" y="3" rx="2" />
</svg>`,
			want: `<rect width="18" height="18" x="3" y="3" rx="2" />`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSVGContent(tt.input)
			if got != tt.want {
				t.Errorf("extractSVGContent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProcessIcon(t *testing.T) {
	tmpDir := t.TempDir()

	testSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24">
  <circle cx="12" cy="12" r="10" />
</svg>`
	svgPath := filepath.Join(tmpDir, "test-icon.svg")
	if err := os.WriteFile(svgPath, []byte(testSVG), 0o644); err != nil {
		t.Fatalf("failed to create test SVG: %v", err)
	}

	icon, err := processIcon(svgPath, tmpDir)
	if err != nil {
		t.Fatalf("processIcon() failed: %v", err)
	}

	if icon.Name != "test-icon" {
		t.Errorf("processIcon().Name = %q, want %q", icon.Name, "test-icon")
	}

	if icon.PascalName != "TestIcon" {
		t.Errorf("processIcon().PascalName = %q, want %q", icon.PascalName, "TestIcon")
	}

	if icon.Paths == "" {
		t.Error("processIcon().Paths is empty")
	}

	wantPaths := `<circle cx="12" cy="12" r="10" />`
	if icon.Paths != wantPaths {
		t.Errorf("processIcon().Paths = %q, want %q", icon.Paths, wantPaths)
	}

	if len(icon.Aliases) != 0 {
		t.Errorf("processIcon().Aliases length = %d, want 0", len(icon.Aliases))
	}
}

func TestProcessIconWithAliases(t *testing.T) {
	tmpDir := t.TempDir()

	testSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24">
  <circle cx="12" cy="12" r="10" />
</svg>`
	svgPath := filepath.Join(tmpDir, "test-icon.svg")
	if err := os.WriteFile(svgPath, []byte(testSVG), 0o644); err != nil {
		t.Fatalf("failed to create test SVG: %v", err)
	}

	testMetadata := `{
  "aliases": [
    {
      "name": "old-name",
      "deprecationReason": "Renamed for clarity",
      "deprecated": true
    },
    {
      "name": "another-alias",
      "deprecated": false
    }
  ]
}`
	metadataPath := filepath.Join(tmpDir, "test-icon.json")
	if err := os.WriteFile(metadataPath, []byte(testMetadata), 0o644); err != nil {
		t.Fatalf("failed to create test metadata: %v", err)
	}

	icon, err := processIcon(svgPath, tmpDir)
	if err != nil {
		t.Fatalf("processIcon() failed: %v", err)
	}

	if len(icon.Aliases) != 2 {
		t.Fatalf("processIcon().Aliases length = %d, want 2", len(icon.Aliases))
	}

	alias1 := icon.Aliases[0]
	if alias1.Name != "old-name" {
		t.Errorf("alias[0].Name = %q, want %q", alias1.Name, "old-name")
	}
	if alias1.PascalName != "OldName" {
		t.Errorf("alias[0].PascalName = %q, want %q", alias1.PascalName, "OldName")
	}
	if alias1.TargetName != "test-icon" {
		t.Errorf("alias[0].TargetName = %q, want %q", alias1.TargetName, "test-icon")
	}
	if alias1.TargetPascalName != "TestIcon" {
		t.Errorf("alias[0].TargetPascalName = %q, want %q", alias1.TargetPascalName, "TestIcon")
	}
	if !alias1.Deprecated {
		t.Error("alias[0].Deprecated = false, want true")
	}
	if alias1.DeprecationReason != "Renamed for clarity" {
		t.Errorf("alias[0].DeprecationReason = %q, want %q", alias1.DeprecationReason, "Renamed for clarity")
	}

	alias2 := icon.Aliases[1]
	if alias2.Name != "another-alias" {
		t.Errorf("alias[1].Name = %q, want %q", alias2.Name, "another-alias")
	}
	if alias2.PascalName != "AnotherAlias" {
		t.Errorf("alias[1].PascalName = %q, want %q", alias2.PascalName, "AnotherAlias")
	}
	if alias2.Deprecated {
		t.Error("alias[1].Deprecated = true, want false")
	}
}

func TestReadMetadata(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		wantAliases int
		wantErr     bool
	}{
		{
			name: "valid metadata with aliases",
			content: `{
  "aliases": [
    {
      "name": "test-alias",
      "deprecated": true
    }
  ]
}`,
			wantAliases: 1,
			wantErr:     false,
		},
		{
			name:        "valid metadata without aliases",
			content:     `{"tags": ["test"]}`,
			wantAliases: 0,
			wantErr:     false,
		},
		{
			name:        "invalid JSON",
			content:     `{invalid json}`,
			wantAliases: 0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, "test.json")
			if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			metadata, err := readMetadata(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("readMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && len(metadata.Aliases) != tt.wantAliases {
				t.Errorf("readMetadata() aliases count = %d, want %d", len(metadata.Aliases), tt.wantAliases)
			}

			_ = os.Remove(path)
		})
	}
}

func TestReadMetadataFileNotFound(t *testing.T) {
	_, err := readMetadata("/nonexistent/file.json")
	if err == nil {
		t.Error("readMetadata() should return error for nonexistent file")
	}
}
