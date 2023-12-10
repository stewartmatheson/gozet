package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type IndexPageData struct {
	Notes []Note
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatalln("Must pass at least one arg")
	}

	if args[0] == "touch" {
		fileName, err := createOrGetNote(strings.Join(args[1:], " "))
		if err != nil {
			panic(err)
		}
		fmt.Println(fileName)
		os.Exit(0)
	}

	if args[0] == "ls" {
		for _, noteFileName := range allNoteFiles() {
			fmt.Println(noteFileName)
		}
		os.Exit(0)
	}

	if args[0] == "build" {
		build(Page{
			Name: "index",
			Data: IndexPageData{
				Notes: allNotes(),
			},
		})

		for _, note := range allNotes() {
			build(note)
		}

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
