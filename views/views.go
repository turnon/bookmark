package views

import (
	"bytes"
	"embed"
	"errors"
	"text/template"
)

var (
	//go:embed *
	views embed.FS

	templates = make(map[string]*template.Template)
)

func init() {
	files, _ := views.ReadDir(".")
	for _, f := range files {
		t, _ := template.ParseFS(views, f.Name())
		templates[f.Name()] = t
	}
}

func Render(templateName string, obj interface{}) ([]byte, error) {
	t, ok := templates[templateName]
	if !ok {
		return nil, errors.New("template not found: " + templateName)
	}
	var buffer bytes.Buffer
	t.Execute(&buffer, obj)
	return buffer.Bytes(), nil
}
