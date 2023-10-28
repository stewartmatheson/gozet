package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"slices"
)

func createOrGetNote(title string) (string, error) {
	note := Note{
		Body: "",
		Meta: Meta{
			Tags:      []string{},
			Title:     title,
			CreatedAt: time.Now(),
		},
	}

	found, fileName := findNoteBySlug(note.slug())

	if found {
		return fileName, nil
	}

	fileName, err := note.create()

	if err != nil {
		return "", err
	}

	return fileName, nil
}

func findNoteBySlug(slug string) (bool, string) {
	for _, note := range allNotes() {
		if note.slug() == slug {
			return true, note.fileName()
		}
	}
	return false, ""
}

func findRelatedNotes(note Note) []Note {
	realtedNotes := make(map[string]Note)
	for _, currentNote := range allNotes() {
		if currentNote.slug() == note.slug() {
			continue
		}
		for _, tag := range currentNote.Meta.Tags {
			if slices.Contains(note.Meta.Tags, tag) {
				realtedNotes[currentNote.slug()] = currentNote
			}
		}
	}

	notes := make([]Note, 0, len(realtedNotes))
	for _, relatedNote := range realtedNotes {
		notes = append(notes, relatedNote)
	}
	return notes
}

func allFiles() []string {
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

func allNotes() []Note {
	files := allFiles()
	var notes []Note
	for _, noteFile := range files {
		note, err := read(noteFile)
		if err != nil {
			fmt.Println("Can't parse note: " + noteFile)
			continue
		}
		notes = append(notes, *note)
	}
	return notes
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
			fmt.Println(err)
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
