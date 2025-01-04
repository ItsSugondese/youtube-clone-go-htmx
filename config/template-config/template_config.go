package template_config

import (
	"html/template"
	"os"
	"path/filepath"
)

func LoadTemplates(templatesDir string) *template.Template {
	tmpl := template.New("")
	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err := tmpl.ParseFiles(path)
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return tmpl
}
