// Package lucide provides Lucide icons for Go html/template.
//
// Use Lucide icons directly in your Go templates without JavaScript dependencies.
// Perfect for server-side rendered applications.
//
// Basic usage:
//
//	tmpl.Funcs(lucide.FuncMap())
//
// Then in your template:
//
//	{{ lucide "circle-x" }}
//	{{ lucide "play" (dict "size" 32 "class" "text-blue-500") }}
//
// For more examples and documentation, visit:
// https://github.com/kaugesaar/lucide-go
package lucide

import (
	"fmt"
	"html/template"
)

// Config configures the template function map returned by FuncMap.
type Config struct {
	// FuncName is the icon function name (default: "lucide")
	FuncName string

	// SkipDict disables the dict helper function (default: false, meaning dict is included)
	// Set to true only if you're using Sprig or providing your own dict function
	SkipDict bool

	// DictName is the dict function name (default: "dict")
	DictName string
}

// Options configures individual icon rendering.
type Options struct {
	// Size sets both width and height in pixels (default: 24)
	Size int

	// Color sets the color of the stroke. (default: currentColor)
	Color string

	// StrokeWidth sets the stroke width (default: 2)
	StrokeWidth int

	// Class sets CSS classes to add to the SVG element
	Class string
}

// iconRegistry maps icon names to their rendering functions.
// This will be populated by the generated icons.go file.
var iconRegistry = make(map[string]func(opts ...Options) template.HTML)

// Icon renders an icon by name with optional configuration.
// This is the main template function.
//
// Usage in templates:
//
//	{{ lucide "circle-x" }}
//	{{ lucide "play" (dict "size" 32) }}
//	{{ lucide "menu" (dict "size" 24 "color" "red" "strokeWidth" 2 "class" "my-icon") }}
func Icon(name string, options ...map[string]any) template.HTML {
	opts := Options{
		Size:        24,
		Color:       "currentColor",
		StrokeWidth: 2,
		Class:       "",
	}

	if len(options) > 0 {
		if size, ok := options[0]["size"].(int); ok {
			opts.Size = size
		}
		if color, ok := options[0]["color"].(string); ok {
			opts.Color = color
		}
		if strokeWidth, ok := options[0]["strokeWidth"].(int); ok {
			opts.StrokeWidth = strokeWidth
		}
		if class, ok := options[0]["class"].(string); ok {
			opts.Class = class
		}
	}

	iconFn, ok := iconRegistry[name]
	if !ok {
		return template.HTML("")
	}

	return iconFn(opts)
}

// FuncMap returns a template.FuncMap with icon functions registered.
// By default, includes both the "lucide" icon function and "dict" helper.
//
// Usage:
//
//	// Simple - includes lucide and dict functions
//	tmpl.Funcs(lucide.FuncMap())
//
//	// Custom function names (dict still included)
//	tmpl.Funcs(lucide.FuncMap(&lucide.Config{
//	    FuncName: "icon",
//	    DictName: "opts",
//	}))
//
//	// Disable dict helper (if using Sprig or providing your own)
//	tmpl.Funcs(lucide.FuncMap(&lucide.Config{
//	    SkipDict: true,
//	}))
func FuncMap(cfg ...*Config) template.FuncMap {
	funcName := "lucide"
	skipDict := false
	dictName := "dict"

	if len(cfg) > 0 && cfg[0] != nil {
		if cfg[0].FuncName != "" {
			funcName = cfg[0].FuncName
		}
		skipDict = cfg[0].SkipDict
		if cfg[0].DictName != "" {
			dictName = cfg[0].DictName
		}
	}

	fm := template.FuncMap{
		funcName: Icon,
	}

	if !skipDict {
		fm[dictName] = Dict
	}

	return fm
}

// Dict creates a map from key-value pairs.
// Helper function for building option maps in templates.
//
// If an odd number of arguments is provided, the last key gets an empty string value.
// Non-string keys are converted to strings using fmt.Sprint.
//
// Usage in templates:
//
//	{{ lucide "circle-x" (dict "size" 32 "class" "my-icon") }}
//
// Can also be registered manually:
//
//	tmpl.Funcs(template.FuncMap{
//	    "lucide": lucide.Icon,
//	    "dict":   lucide.Dict,
//	})
func Dict(values ...any) map[string]any {
	d := make(map[string]any, len(values)/2)
	lenv := len(values)

	for i := 0; i < lenv; i += 2 {
		key := fmt.Sprint(values[i])
		if i+1 >= lenv {
			d[key] = ""
			continue
		}
		d[key] = values[i+1]
	}

	return d
}

// buildSVG constructs an SVG string with the given parameters.
// This is a helper function used by generated icon functions.
func buildSVG(paths string, opts Options) template.HTML {
	classAttr := ""
	if opts.Class != "" {
		classAttr = fmt.Sprintf(` class="%s"`, opts.Class)
	}

	svg := fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 24 24" fill="none" stroke="%s" stroke-width="%d" stroke-linecap="round" stroke-linejoin="round"%s>%s</svg>`,
		opts.Size,
		opts.Size,
		opts.Color,
		opts.StrokeWidth,
		classAttr,
		paths,
	)

	return template.HTML(svg)
}

// registerIcon registers an icon function in the global registry.
// This is called by generated code in icons.go.
func registerIcon(name string, fn func(opts ...Options) template.HTML) {
	iconRegistry[name] = fn
}
