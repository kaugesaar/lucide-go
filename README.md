# Lucide Go

Lucide icons for Go. Each icon renders as an inline SVG and can be used with any template system (`html/template`, `templ`, etc.) or directly in Go code. 

## Installation

```bash
go get github.com/kaugesaar/lucide-go
```

## Usage

### With html/template

Register the icon function in your template:

```go
package main

import (
    "html/template"
    "github.com/kaugesaar/lucide-go"
)

func main() {
    tmpl := template.New("index.html")

    tmpl.Funcs(lucide.FuncMap())

    // Parse and execute your templates...
}
```

Then use icons in your templates:

```html
<!-- Simple usage with defaults (24px, stroke-width 2) -->
{{ lucide "circle-x" }}

<!-- With custom size -->
{{ lucide "chevron-down" (dict "size" 32) }}

<!-- With custom size, color, stroke width, and CSS class -->
{{ lucide "play" (dict "size" 48 "color" "red" "strokeWidth" 3 "class" "hover:text-red-800") }}

<!-- Multiple classes (Tailwind example) -->
{{ lucide "menu" (dict "class" "w-6 h-6 text-gray-700 hover:text-gray-900") }}
```

**Available options:**

| name          | type     | default      |
| ------------- | -------- | ------------ |
| `size`        | *int*    | 24           |
| `color`       | *string* | currentColor |
| `strokeWidth` | *int*    | 2            |
| `class`       | *string* |              |

#### Configuration Options

**Customize function names** to avoid conflicts:

```go
tmpl.Funcs(lucide.FuncMap(&lucide.Config{
    FuncName: "icon",  // Use {{ icon "..." }} instead of {{ lucide "..." }}
    DictName: "opts",  // Use {{ opts ... }} instead of {{ dict ... }}
}))
```

**Disable dict helper** (e.g. using [Sprig](https://masterminds.github.io/sprig/) or providing your own):

```go
tmpl.Funcs(lucide.FuncMap(&lucide.Config{
    SkipDict: true,
}))
```

**Manual registration**:

```go
tmpl.Funcs(template.FuncMap{
    "lucide": lucide.Icon,
    "dict":   lucide.Dict,
})
```

### With Other Template Systems

For template systems other than `html/template`, call `Icon()` directly from your Go code:

```go
import "github.com/kaugesaar/lucide-go"

// Dynamic icon lookup by name
icon := lucide.Icon("circle-x", map[string]any{
    "size":  32,
    "class": "logo",
})
// Returns template.HTML - integrate with your template system

// Or use individual functions when icon is known at compile time
icon := lucide.Menu(lucide.Options{
    Size:  24,
    Class: "menu-icon",
})
```

**Integration examples:**
- **templ**: `@templ.Raw(icon)`
- **Direct rendering**: `fmt.Fprintf(w, "%s", icon)`
- **Other systems**: Use `template.HTML` value as needed

### Direct Go Usage

Use icons directly in Go code with two approaches:

**Dynamic lookup by name**:

```go
package main

import (
    "fmt"
    "github.com/kaugesaar/lucide-go"
)

func main() {
    // Look up icon by name - great for dynamic scenarios
    icon := lucide.Icon("circle-x", map[string]any{
        "size":        32,
        "strokeWidth": 2,
        "class":       "icon-danger",
    })

    fmt.Println(icon)
}
```

**Type-safe individual functions**:

```go
package main

import (
    "fmt"
    "github.com/kaugesaar/lucide-go"
)

func main() {
    // Using individual icon functions
    icon := lucide.CircleX(lucide.Options{
        Size:        32,
        StrokeWidth: 2,
        Class:       "icon-danger",
    })

    fmt.Println(icon)

    // Or use defaults (24px, stroke-width 2)
    icon = lucide.Play()
}
```

## API

### `Icon(name string, options ...map[string]interface{}) template.HTML`

The main template function. Returns an SVG as `template.HTML`.

**Parameters:**
- `name`: Icon name (e.g., `"circle-x"`, `"chevron-down"`)
- `options`: Optional map with:
  - `size` (int): Width and height in pixels (default: 24)
  - `color` (string): Stroke color (default: currentColor)
  - `strokeWidth` (int): Stroke width (default: 2)
  - `class` (string): CSS classes to add

### `FuncMap(cfg ...*Config) template.FuncMap`

Returns a `template.FuncMap` for registering with templates. By default includes both the icon function and dict helper. Accepts optional configuration.

### `Config` struct

```go
type Config struct {
    FuncName string // Icon function name (default: "lucide")
    SkipDict bool   // Disable dict helper (default: false)
    DictName string // Dict function name (default: "dict")
}
```

**Fields:**
- `FuncName`: Customize the icon function name. Default is `"lucide"`.
- `SkipDict`: Set to `true` to disable the dict helper. Default is `false` (dict is included).
- `DictName`: Customize the dict function name. Default is `"dict"`.

### Individual Icon Functions

Every icon has a corresponding Go function:

```go
func CircleX(opts ...Options) template.HTML
func ChevronDown(opts ...Options) template.HTML
func Play(opts ...Options) template.HTML
// ... ~1600 more
```

### `Options` struct

```go
type Options struct {
    Size        int    // Width/height in pixels (default: 24)
    Color       string // Stroke color (default: currentColor)
    StrokeWidth int    // Stroke width (default: 2)
    Class       string // CSS classes
}
```

## License

This project is licensed under the MIT License.

Lucide icons are licensed under the ISC License.

See [LICENSE](./LICENSE) for details.
