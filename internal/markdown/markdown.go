package markdown

import (
	"bytes"
	"fmt"
	"html"
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

var mermaidBlock = regexp.MustCompile(`(?s)<pre><code class="language-mermaid">(.*?)</code></pre>`)

// Document is a Markdown file split into optional front matter and body.
type Document struct {
	FrontMatter map[string]any
	Body        []byte
}

// ParseDocument splits YAML front matter from Markdown body.
func ParseDocument(raw []byte) (Document, error) {
	if !bytes.HasPrefix(raw, []byte("---\n")) {
		return Document{Body: raw}, nil
	}

	parts := bytes.SplitN(raw, []byte("\n---\n"), 2)
	if len(parts) != 2 {
		return Document{}, fmt.Errorf("invalid front matter delimiter")
	}

	var fm map[string]any
	if err := yaml.Unmarshal(parts[0][4:], &fm); err != nil {
		return Document{}, fmt.Errorf("parse front matter: %w", err)
	}

	return Document{
		FrontMatter: fm,
		Body:        parts[1],
	}, nil
}

// Converter renders Markdown to HTML with Mermaid support.
type Converter struct {
	engine goldmark.Markdown
}

// NewConverter returns a GFM-compatible Markdown converter.
func NewConverter() *Converter {
	engine := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			gmhtml.WithHardWraps(),
			gmhtml.WithXHTML(),
			gmhtml.WithUnsafe(),
		),
	)
	return &Converter{engine: engine}
}

// ToHTML converts Markdown to HTML and rewrites Mermaid fenced blocks.
func (c *Converter) ToHTML(source []byte) (string, error) {
	var buf bytes.Buffer
	if err := c.engine.Convert(source, &buf); err != nil {
		return "", fmt.Errorf("convert markdown: %w", err)
	}

	return string(transformMermaid(buf.Bytes())), nil
}

func transformMermaid(htmlDoc []byte) []byte {
	return mermaidBlock.ReplaceAllFunc(htmlDoc, func(match []byte) []byte {
		sub := mermaidBlock.FindSubmatch(match)
		if len(sub) != 2 {
			return match
		}
		content := html.UnescapeString(string(sub[1]))
		return []byte(`<pre class="mermaid">` + content + `</pre>`)
	})
}
