package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func write() {
	writeNotes()
	writeIndex()
}

type IndexTemplateData struct {
	Notes []Note
}

func writeIndex() {
	indexTemplateFile := getConfiguration().Home + "/templates/index.html"
	indexTemplate, err := template.New("index.html").ParseFiles(indexTemplateFile)

	indexTemplateData := IndexTemplateData{
		Notes: allNotes(),
	}

	if err != nil {
		panic(err)
	}

	renderedPageBuffer := new(bytes.Buffer)
	err = indexTemplate.Execute(renderedPageBuffer, indexTemplateData)

	renderedTemplate, templateRenderErr := ZetTemplate{
		Title:    "Index",
		Body:     string(renderedPageBuffer.Bytes()),
		Keywords: []string{},
	}.render()

	if templateRenderErr != nil {
		fmt.Println("Can't render template for index")
		return
	}

	outPath := "/index.html"
	writeError := writeFile(outPath, renderedTemplate)

	if writeError != nil {
		fmt.Println(writeError)
	}
}

func writeNotes() {
	for _, noteFile := range allNoteFiles() {
		note, err := read(noteFile)
		if err != nil {
			fmt.Println("Can't parse note: " + noteFile)
			continue
		}

		renderedNote, noteRenderErr := note.render()

		if noteRenderErr != nil {
			fmt.Println("Can't render note: " + noteFile)
			fmt.Println(noteRenderErr)
			continue
		}

		outPath := "/notes" +
			"/" + note.fileDatePrefix() +
			"/" +
			note.slug() +
			".html"

		noteTemplate, templateRenderErr := ZetTemplate{
			Body:     string(renderedNote),
			Keywords: note.Meta.Tags,
			Title:    note.Meta.Title,
		}.render()

		if templateRenderErr != nil {
			fmt.Println("Can't render template: " + noteFile)
			fmt.Println(err)
			continue
		}

		writeError := writeFile(outPath, noteTemplate)

		if writeError != nil {
			fmt.Println(writeError)
		}
	}
}

func writeFile(outPath string, content []byte) error {
	pathToWrite := getConfiguration().Home + "/build" + outPath
	os.MkdirAll(filepath.Dir(pathToWrite), os.FileMode(int(0755)))
	outFile, createFileWriterErr := os.Create(pathToWrite)

	if createFileWriterErr != nil {
		return createFileWriterErr
	}

	defer outFile.Close()

	_, writeError := outFile.Write(content)

	if writeError != nil {
		return writeError
	}

	return nil
}
