package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

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

	if args[0] == "note" {
		fileName, err := createOrGetNote(strings.Join(args[1:], " "))
		if err != nil {
			panic(err)
		}
		fmt.Println(fileName)
		os.Exit(0)
	}

	if args[0] == "list" {
		for _, noteFileName := range allFiles() {
			fmt.Println(noteFileName)
		}
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
}
