package main

import (
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
	tests := []struct {
		name       string
		filename   string
		wantName   string
		wantPascal string
	}{
		{
			name:       "circle-x icon",
			filename:   "../../lucide/icons/circle-x.svg",
			wantName:   "circle-x",
			wantPascal: "CircleX",
		},
		{
			name:       "chevron-down icon",
			filename:   "../../lucide/icons/chevron-down.svg",
			wantName:   "chevron-down",
			wantPascal: "ChevronDown",
		},
		{
			name:       "a-arrow-down icon",
			filename:   "../../lucide/icons/a-arrow-down.svg",
			wantName:   "a-arrow-down",
			wantPascal: "AArrowDown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			icon, err := processIcon(tt.filename)
			if err != nil {
				t.Fatalf("processIcon(%q) failed: %v", tt.filename, err)
			}

			if icon.Name != tt.wantName {
				t.Errorf("processIcon().Name = %q, want %q", icon.Name, tt.wantName)
			}

			if icon.PascalName != tt.wantPascal {
				t.Errorf("processIcon().PascalName = %q, want %q", icon.PascalName, tt.wantPascal)
			}

			if icon.Paths == "" {
				t.Error("processIcon().Paths is empty")
			}

			if len(icon.Paths) > 0 && (icon.Paths[0] == '<' && icon.Paths[1] == 's') {
				t.Error("processIcon().Paths should not start with <svg")
			}
		})
	}
}
