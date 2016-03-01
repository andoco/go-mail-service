package template

import (
	"fmt"
	"log"

	"github.com/cbroglie/mustache"
)

type Engine struct {
	Id string
}

type Template struct {
	Id            string
	EngineId      string
	Content       string
	DefaultFields map[string]interface{}
}

func loadTemplate(templateId string) (string, error) {
	log.Printf("Loading template with id %s", templateId)

	data, err := Asset(fmt.Sprintf("%s.mustache", templateId))
	if err != nil {
		return "", fmt.Errorf("error loading template with id %s; %v", templateId, err)
	}

	log.Printf("Loaded template %s:\n%s", templateId, string(data))

	return string(data), nil
}

func Render(templateId string, fields map[string]interface{}) (string, error) {
	log.Printf("Rendering template %s", templateId)

	tmpl, err := loadTemplate(templateId)
	if err != nil {
		return "", fmt.Errorf("error loading template %s; %v", templateId, err)
	}

	rendered, err := mustache.Render(tmpl, fields)
	if err != nil {
		return "", fmt.Errorf("error rendering template; %v", err)
	}

	return rendered, nil
}
