package main

import (
	"fmt"
	"time"
)

func main() {
	note := Note{
		Title:     "This is a sample note",
		Tags:      []string{"Note", "Test", "There", "Are", "Lots"},
		Body:      "",
		CreatedAt: time.Now(),
	}

	fileName, err := note.create()

	if err != nil {
		panic(err)
	}

	fmt.Println(fileName)
}
