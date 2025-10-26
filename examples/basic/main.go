package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kaugesaar/lucide-go"
)

func main() {
	tmpl := template.Must(template.New("index").Funcs(lucide.FuncMap()).Parse(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Lucide Go Example</title>
    <style>
        body {
            font-family: system-ui, -apple-system, sans-serif;
            max-width: 800px;
            margin: 2rem auto;
            padding: 0 1rem;
        }
        h1 {
            color: #333;
        }
        .examples {
            display: grid;
            gap: 2rem;
            margin-top: 2rem;
        }
        .example {
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 1.5rem;
        }
        .example h3 {
            margin-top: 0;
            color: #555;
        }
        .icon-row {
            display: flex;
            gap: 1rem;
            align-items: center;
            margin: 0.5rem 0;
        }
        .icon-row code {
            background: #f5f5f5;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            font-size: 0.9em;
        }
        .red { color: #ef4444; }
        .green { color: #22c55e; }
        .blue { color: #3b82f6; }
        .gray { color: #6b7280; }
    </style>
</head>
<body>
    <h1>{{ lucide "sparkles" (dict "size" 32) }} Lucide Go Examples</h1>

    <div class="examples">
        <div class="example">
            <h3>Basic Icons (defaults: 24px, stroke-width 2)</h3>
            <div class="icon-row">
                {{ lucide "circle-x" }}
                {{ lucide "check-circle" }}
                {{ lucide "alert-circle" }}
                {{ lucide "info" }}
                <code>{{ "{{" }} lucide "circle-x" {{ "}}" }}</code>
            </div>
        </div>

        <div class="example">
            <h3>Custom Sizes</h3>
            <div class="icon-row">
                {{ lucide "heart" (dict "size" 16) }}
                {{ lucide "heart" (dict "size" 24) }}
                {{ lucide "heart" (dict "size" 32) }}
                {{ lucide "heart" (dict "size" 48) }}
                <code>{{ "{{" }} lucide "heart" (dict "size" 48) {{ "}}" }}</code>
            </div>
        </div>

        <div class="example">
            <h3>Custom Stroke Width</h3>
            <div class="icon-row">
                {{ lucide "star" (dict "strokeWidth" 1) }}
                {{ lucide "star" (dict "strokeWidth" 2) }}
                {{ lucide "star" (dict "strokeWidth" 3) }}
                {{ lucide "star" (dict "strokeWidth" 4) }}
                <code>{{ "{{" }} lucide "star" (dict "strokeWidth" 3) {{ "}}" }}</code>
            </div>
        </div>

        <div class="example">
            <h3>Custom Stroke Color</h3>
            <div class="icon-row">
                {{ lucide "circle" (dict "color" "red") }}
                {{ lucide "square" (dict "color" "green") }}
                {{ lucide "triangle" (dict "color" "blue") }}
                {{ lucide "hexagon" (dict "color" "gray") }}
                <code>{{ "{{" }} lucide "circle" (dict "color" "red") {{ "}}" }}</code>
            </div>
        </div>

        <div class="example">
            <h3>With CSS Classes</h3>
            <div class="icon-row">
                {{ lucide "circle" (dict "class" "red") }}
                {{ lucide "square" (dict "class" "green") }}
                {{ lucide "triangle" (dict "class" "blue") }}
                {{ lucide "hexagon" (dict "class" "gray") }}
                <code>{{ "{{" }} lucide "circle" (dict "class" "red") {{ "}}" }}</code>
            </div>
        </div>

        <div class="example">
            <h3>Combined Options</h3>
            <div class="icon-row">
                {{ lucide "zap" (dict "size" 32 "strokeWidth" 3 "class" "blue") }}
                {{ lucide "flame" (dict "size" 32 "strokeWidth" 3 "class" "red") }}
                {{ lucide "droplet" (dict "size" 32 "strokeWidth" 3 "class" "blue") }}
                <code>{{ "{{" }} lucide "zap" (dict "size" 32 "strokeWidth" 3 "class" "blue") {{ "}}" }}</code>
            </div>
        </div>

        <div class="example">
            <h3>Common UI Icons</h3>
            <div class="icon-row">
                {{ lucide "menu" }}
                {{ lucide "x" }}
                {{ lucide "search" }}
                {{ lucide "settings" }}
                {{ lucide "user" }}
                {{ lucide "bell" }}
                {{ lucide "calendar" }}
                {{ lucide "home" }}
            </div>
        </div>

        <div class="example">
            <h3>Navigation Icons</h3>
            <div class="icon-row">
                {{ lucide "chevron-left" }}
                {{ lucide "chevron-right" }}
                {{ lucide "chevron-up" }}
                {{ lucide "chevron-down" }}
                {{ lucide "arrow-left" }}
                {{ lucide "arrow-right" }}
                {{ lucide "arrow-up" }}
                {{ lucide "arrow-down" }}
            </div>
        </div>
    </div>
</body>
</html>
`))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	addr := ":8080"
	fmt.Printf("Server running at http://localhost%s\n", addr)
	fmt.Println("Press Ctrl+C to stop")
	log.Fatal(http.ListenAndServe(addr, nil))
}
