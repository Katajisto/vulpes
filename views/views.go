package views

// View creating functionality here
// view is a wrapper around a template.Template

import (
	"html/template"
	"io"
	"path/filepath"
)

type View struct {
	Template *template.Template
	Layout   string
}

func componentFiles() []string {
	files, err := filepath.Glob("views/templates/components/*.tmpl")
	if err != nil {
		panic(err)
	}
	return files
}

func NewView(layout string, files ...string) *View {
	files = append(componentFiles(), files...)
	return &View{
		Template: template.Must(template.ParseFiles(files...)),
		Layout:   layout,
	}
}

func (v *View) Render(w io.Writer, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}
