package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"text/template"
)

type BuildableTemplateData struct {
	Title    string
	Keywords []string
	Data     any
}

type Buildable interface {
	templatePaths() []string
	buildPath() string
	templateData() BuildableTemplateData
}

func shouldBuild(buildable Buildable) bool {
	buildFile, existingOutFileInfoErr := os.Stat(buildable.buildPath())
	if errors.Is(existingOutFileInfoErr, os.ErrNotExist) {
		return true
	}

	shouldBuild := false
	for _, templatePath := range buildable.templatePaths() {
		templateFile, templateFileErr := os.Stat(templatePath)
		if templateFileErr != nil {
			panic(templateFileErr)
		}

		if buildFile.ModTime().Before(templateFile.ModTime()) {
			shouldBuild = true
		}
	}
	return shouldBuild
}

func build(buildable Buildable) {
	/*
		if !shouldBuild(buildable) {
			return
		}
	*/

	t := template.Must(template.ParseFiles(buildable.templatePaths()...))
	templateBuffer := new(bytes.Buffer)
	templateExecuteErr := t.Execute(templateBuffer, buildable.templateData())
	if templateExecuteErr != nil {
		panic(templateExecuteErr)
	}

	buildPath := buildable.buildPath()
	os.MkdirAll(filepath.Dir(buildPath), os.FileMode(int(0755)))
	buildFile, buildFileErr := os.Create(buildPath)
	if buildFileErr != nil {
		panic(buildFileErr)
	}

	_, writeErr := buildFile.Write(templateBuffer.Bytes())
	if writeErr != nil {
		panic(writeErr)
	}
}
