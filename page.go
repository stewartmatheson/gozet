package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func render(note Note) []byte {
	p := parser.New()
	doc := p.Parse([]byte(note.Body))
	opts := html.RendererOptions{}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
