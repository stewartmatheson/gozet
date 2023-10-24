package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	/*
			reader := bufio.NewReader(os.Stdin)
			fileName, err := reader.ReadString('\n')

			if err != nil {
				panic(err)
			}

			read(fileName)
		  createTestNote()
	*/

	reader := bufio.NewReader(os.Stdin)
	fileName, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	note := read(strings.TrimSuffix(fileName, "\n"))
	fmt.Print(render(note))
}
