package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func read(fileName string) (*Note, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	readingMetadata := false

	metadata := []string{}
	content := []string{}

	for scanner.Scan() {
		if scanner.Text() == "---" {
			readingMetadata = true
			continue
		}

		if readingMetadata {
			metadata = append(metadata, scanner.Text())
		}

		if !readingMetadata {
			content = append(content, scanner.Text())
		}
	}

	meta := Meta{}
	yamlErr := yaml.Unmarshal([]byte(strings.Join(metadata, "\n")), &meta)

	if yamlErr != nil {
		return nil, yamlErr
	}

	return &Note{
		Body: strings.Join(content, "\n"),
		Meta: meta,
	}, nil
}

type Meta struct {
	Title     string
	Tags      []string
	CreatedAt time.Time
}

type Note struct {
	Body string
	Meta Meta
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
	yamlData, err := yaml.Marshal(&note.Meta)

	if err != nil {
		return err
	}

	writer.Write([]byte("# " + note.Meta.Title + "\n"))
	writer.Write([]byte("\n"))
	writer.Write([]byte(note.Body + "\n"))
	writer.Write([]byte("\n"))
	writer.Write([]byte("---" + "\n"))
	writer.Write([]byte(yamlData))
	return nil
}

func (note Note) slug() string {
	return strings.ToLower(strings.ReplaceAll(note.Meta.Title, " ", "-"))
}

func (note Note) fileDatePrefix() string {
	return note.Meta.CreatedAt.Format("2006/01/02")
}

func (note Note) fileName() string {
	return getConfiguration().Home +
		"/content/" +
		note.fileDatePrefix() +
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
