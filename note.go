package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/yaml.v3"
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

func (note Note) relatedNotes() []Note {
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

func allNoteFiles() []string {
	notesPath := getConfiguration().Home + "/notes"
	var files []string
	err := filepath.Walk(notesPath,
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
	files := allNoteFiles()
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

// This method should just panic
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
		"/notes/" +
		note.fileDatePrefix() +
		"/" +
		note.slug() +
		".md"
}

func (note Note) LinkTo() string {
	return "/notes/" +
		note.fileDatePrefix() +
		"/" +
		note.slug() +
		".html"
}

func (note Note) buildPath() string {
	return getConfiguration().Home + "/build" + note.LinkTo()
}

func (note Note) templatePaths() []string {
	return []string{
		getConfiguration().Home + "/templates/base.html",
		getConfiguration().Home + "/templates/note.html",
	}
}

type NoteTemplateData struct {
	Body    string
	Related []Note
}

func (note Note) templateData() BuildableTemplateData {
	return BuildableTemplateData{
		Title:    note.Meta.Title,
		Keywords: note.Meta.Tags,
		Data: NoteTemplateData{
			Body:    string(note.renderBody()),
			Related: note.relatedNotes(),
		},
	}
}

func (note Note) renderBody() []byte {
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(note.Body))
	opts := html.RendererOptions{}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
