package main

// View creating functionality here
// view is a wrapper around a template.Template

import (
	"html/template"
	"path/filepath"
)

type view struct {
	Template *template.Template
	Layout   string
}

var homeView view

func init() {
	homeView = newView("main", "templates/home.tmpl")
}

func componentFiles() []string {
	files, err := filepath.Glob("templates/components/*.tmpl")
	if err != nil {
		panic(err)
	}
	return files
}

func newView(layout string, files ...string) view {
	files = append(componentFiles(), files...)
	return view{
		Template: template.Must(template.ParseFiles(files...)),
		Layout:   layout,
	}
}
