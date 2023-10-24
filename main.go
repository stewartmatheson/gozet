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
		os.Exit(1)
	}

	if args[0] == "render" && len(args) > 1 {
		note := read(args[1])
		fmt.Print(string(render(note)))
		os.Exit(0)
	}

	if args[0] == "create" {
		createTestNote()
		os.Exit(0)
	}

	log.Fatalln("Unknown Param " + args[0])
	os.Exit(1)
}
