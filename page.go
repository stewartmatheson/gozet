package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func renderBody(note Note) []byte {
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(note.Body))
	opts := html.RendererOptions{}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

type Page struct {
	Body     string
	Keywords string
	Title    string
	Related  []Note
}

func render(note Note) ([]byte, error) {
	page := Page{
		Body:     string(renderBody(note)),
		Keywords: strings.Join(note.Meta.Tags, ", "),
		Title:    note.Meta.Title,
		Related:  findRelatedNotes(note),
	}
	templateFile := getConfiguration().Home + "/templates/index.html"
	template, err := template.New("index.html").ParseFiles(templateFile)

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = template.Execute(buf, page)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
