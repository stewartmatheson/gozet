package main

type Page struct {
	Name string
	Data any
}

func (page Page) templateData() BuildableTemplateData {
	return BuildableTemplateData{
		Title:    page.Name,
		Keywords: []string{},
		Data:     page.Data,
	}
}

func (page Page) buildPath() string {
	return getConfiguration().Home + "/build/" + page.Name + ".html"
}

func (page Page) templatePaths() []string {
	return []string{
		getConfiguration().Home + "/templates/base.html",
		getConfiguration().Home + "/templates/" + page.Name + ".html",
	}
}
