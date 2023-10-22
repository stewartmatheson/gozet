package main

import (
	"testing"
)

func TestNoteSlug(t *testing.T) {
	note := Note{
		Title: "File Name",
	}

	got := note.fileName()
	want := "file-name.md"

	if got != want {
		t.Errorf("Got %q wanted %q", got, want)
	}
}
