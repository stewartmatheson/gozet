package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func all() []string {
	contentPath := getConfiguration().Home + "/content"
	var files []string
	err := filepath.Walk(contentPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return files
}

func build(files []string) {
	for _, noteFile := range files {
		note, err := read(noteFile)
		if err != nil {
			fmt.Println("Can't parse note: " + noteFile)
			continue
		}

		outPath := getConfiguration().Home +
			"/build" +
			"/" + note.fileDatePrefix() +
			"/" +
			note.slug() +
			".html"

		renderedNote, err := render(*note)

		if err != nil {
			fmt.Println("Can't render note: " + noteFile)
			continue
		}

		os.MkdirAll(filepath.Dir(outPath), os.FileMode(int(0755)))
		outFile, createFileWriterErr := os.Create(outPath)

		if createFileWriterErr != nil {
			fmt.Println("Can't open render page: " + outPath)
			fmt.Println(createFileWriterErr)
			continue
		}

		defer outFile.Close()

		_, writeError := outFile.Write(renderedNote)

		if writeError != nil {
			fmt.Println("Can't write file: " + outPath)
			continue
		}
	}

}
