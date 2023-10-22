package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Note struct {
	Title     string
	Body      string
	Tags      []string
	CreatedAt time.Time
}

func (note Note) create() (path string, err error) {
	os.MkdirAll(
		filepath.Dir(note.fileName()),
		os.FileMode(int(0755)),
	)
	file, err := os.Create(note.fileName())

	if err != nil {
		return "", err
	}

	writeErr := note.write(file)

	if writeErr != nil {
		return "", err
	}

	return note.fileName(), nil
}

func (note Note) write(writer io.Writer) error {
	yamlData, err := yaml.Marshal(&note)

	if err != nil {
		return err
	}

	writer.Write([]byte("# " + note.Title + "\n"))
	writer.Write([]byte("\n"))
	writer.Write([]byte(note.Body + "\n"))
	writer.Write([]byte("\n"))
	writer.Write([]byte("---" + "\n"))
	writer.Write([]byte(yamlData))
	return nil
}

func (note Note) slug() string {
	return strings.ToLower(strings.ReplaceAll(note.Title, " ", "-"))
}

func (note Note) fileName() string {
	datePrefix := note.CreatedAt.Format("2006/01/02")

	return os.Getenv("ZET_HOME") +
		"/content/" +
		datePrefix +
		"/" +
		note.slug() +
		".md"
}

func (note Note) save() {
	file, err := os.Create(note.fileName())
	if err != nil {
		panic(err)
	}

	defer file.Close()
}
