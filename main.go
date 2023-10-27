package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func createTestNote() {
	note := Note{
		Body: "",
		Meta: Meta{
			Tags:      []string{"Note", "Test", "There", "Are", "Lots"},
			Title:     "This is a sample note",
			CreatedAt: time.Now(),
		},
	}

	fileName, err := note.create()

	if err != nil {
		panic(err)
	}

	fmt.Println(fileName)
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatalln("Must pass at least one arg")
	}

	if args[0] == "render" && len(args) > 1 {
		note, err := read(args[1])
		if err != nil {
			panic(err)
		}

		noteContent, noteRenderErr := render(*note)

		if noteRenderErr != nil {
			panic(err)
		}

		fmt.Print(string(noteContent))
		os.Exit(0)
	}

	if args[0] == "create" {
		createTestNote()
		os.Exit(0)
	}

	if args[0] == "list" {
		fmt.Println(allFiles())
		os.Exit(0)
	}

	if args[0] == "build" {
		build(allFiles())
		os.Exit(0)
	}

	if args[0] == "tags" {
		allNotes := allNotes()
		for _, note := range allNotes {
			fmt.Println(note.Meta.Tags)
		}

		os.Exit(0)
	}

	log.Fatalln("Unknown Param " + args[0])
	os.Exit(1)
}
