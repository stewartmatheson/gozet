package main

import (
	"fmt"
	"os"
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

func (note Note) print() {
	yamlData, err := yaml.Marshal(&note)

	if err != nil {
		panic(err)
	}

	fmt.Println("---")
	fmt.Println(string(yamlData))
	fmt.Println("---")
	fmt.Println(note.Body)
}

func (note Note) slug() string {
	return strings.ToLower(strings.ReplaceAll(note.Title, " ", "-"))
}

func (note Note) fileName() string {
	datePrefix := note.CreatedAt.Format("2006/01/02")

	return os.Getenv("ZET_HOME") +
		"/" +
		datePrefix +
		"/" +
		note.slug() +
		".md"
}

func main() {
	note := Note{
		Title:     "This is a sample note",
		Tags:      []string{"Note", "Test"},
		Body:      "This is the body of the note.",
		CreatedAt: time.Now(),
	}

	note.print()
	fmt.Println("---")
	fmt.Println(note.fileName())
}
