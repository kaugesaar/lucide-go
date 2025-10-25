package lucide

import (
	"html/template"
	"strings"
	"testing"
)

func TestIcon(t *testing.T) {
	tests := []struct {
		name    string
		icon    string
		opts    map[string]any
		wantErr bool
		want    string
	}{
		{
			name: "basic icon",
			icon: "circle-x",
			want: "<svg",
		},
		{
			name: "icon with size",
			icon: "circle-x",
			opts: map[string]any{"size": 32},
			want: `width="32" height="32"`,
		},
		{
			name: "icon with class",
			icon: "circle-x",
			opts: map[string]any{"class": "my-icon"},
			want: `class="my-icon"`,
		},
		{
			name: "icon with stroke width",
			icon: "circle-x",
			opts: map[string]any{"strokeWidth": 3},
			want: `stroke-width="3"`,
		},
		{
			name: "non-existent icon",
			icon: "doesnt-exist",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result template.HTML
			if tt.opts != nil {
				result = Icon(tt.icon, tt.opts)
			} else {
				result = Icon(tt.icon)
			}

			got := string(result)
			if !strings.Contains(got, tt.want) {
				t.Errorf("Icon() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestFuncMap(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   []string
	}{
		{
			name:   "default config includes dict",
			config: nil,
			want:   []string{"lucide", "dict"},
		},
		{
			name: "skip dict",
			config: &Config{
				SkipDict: true,
			},
			want: []string{"lucide"},
		},
		{
			name: "custom names",
			config: &Config{
				FuncName: "icon",
				DictName: "opts",
			},
			want: []string{"icon", "opts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fm template.FuncMap
			if tt.config != nil {
				fm = FuncMap(tt.config)
			} else {
				fm = FuncMap()
			}

			for _, key := range tt.want {
				if _, ok := fm[key]; !ok {
					t.Errorf("FuncMap() missing key %q", key)
				}
			}
		})
	}
}

func TestDict(t *testing.T) {
	result := Dict("size", 32, "class", "my-icon")

	if result["size"] != 32 {
		t.Errorf("dict[size] = %v, want 32", result["size"])
	}

	if result["class"] != "my-icon" {
		t.Errorf("dict[class] = %v, want my-icon", result["class"])
	}
}

func TestIndividualIconFunction(t *testing.T) {
	icon := CircleX()
	got := string(icon)

	if !strings.Contains(got, "<svg") {
		t.Errorf("CircleX() didn't return SVG, got: %v", got)
	}

	if !strings.Contains(got, `width="24"`) {
		t.Errorf("CircleX() didn't use default size 24")
	}

	icon = CircleX(Options{Size: 48, Class: "test"})
	got = string(icon)

	if !strings.Contains(got, `width="48"`) {
		t.Errorf("CircleX() didn't respect custom size")
	}

	if !strings.Contains(got, `class="test"`) {
		t.Errorf("CircleX() didn't respect custom class")
	}
}
