package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"os"
	"path"
	"strings"
)

type Template struct {
	templates map[string]*template.Template
	paths     []string
}

func NewTemplateRenderer(paths ...string) *Template {
	return &Template{
		templates: map[string]*template.Template{},
		paths:     paths,
	}
}

func (t *Template) findTemplates(templateList ...string) ([]string, error) {
	var foundTemplates []string
	var err error
	for _, templateFile := range templateList {
		fullFilename := ""
		for _, templatePath := range t.paths {
			fullFilename = path.Join(templatePath, templateFile)
			_, err = os.Stat(fullFilename)
			if err == nil {
				break
			}
		}
		if err != nil {
			return foundTemplates, fmt.Errorf("couldn't find template %s in path: %s", templateFile, strings.Join(t.paths, ", "))
		}
		foundTemplates = append(foundTemplates, fullFilename)
	}
	return foundTemplates, nil
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	templateFileList := strings.Split(name, ",")
	baseTemplateName := templateFileList[0]
	//templateName := templateFileList[len(templateFileList)-1]
	var pageTemplate *template.Template
	pageTemplate, ok := t.templates[name]
	if !ok {
		templateList, err := t.findTemplates(templateFileList...)
		if err != nil {
			return fmt.Errorf("could not find templates for %s, %s", name, err)
		}
		pageTemplate, err = template.New(baseTemplateName).ParseFiles(templateList...)
		if err != nil {
			return fmt.Errorf("could not load templates for %s, %s", strings.Join(templateList, ", "), err)
		}
		//t.templates[name] = pageTemplate
	}

	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
		viewContext["csrf"] = c.Get("csrf")
	}

	return pageTemplate.ExecuteTemplate(w, baseTemplateName, data)
}
