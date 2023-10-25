package main

import (
	"bytes"
	"log"
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
	Body string
}

func render(note Note) []byte {
	page := Page{Body: string(renderBody(note))}
	templateFile := getConfiguration().Home + "/templates/index.html"
	template, err := template.New("index.html").ParseFiles(templateFile)

	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = template.Execute(buf, page)

	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}
