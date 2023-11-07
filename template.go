package main

import (
	"bytes"
	"text/template"
)

type ZetTemplate struct {
	Body     string
	Keywords []string
	Title    string
}

func (zetTemplate ZetTemplate) render() ([]byte, error) {
	templateFile := getConfiguration().Home + "/templates/template.html"
	template, err := template.New("template.html").ParseFiles(templateFile)

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = template.Execute(buf, zetTemplate)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
