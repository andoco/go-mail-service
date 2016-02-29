package template

import (
	"fmt"
	"log"

	"github.com/cbroglie/mustache"

	"bitbucket.org/andoco/gomailservice/models"
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

type Renderer interface {
	Render(msg *models.MailMessage) (string, error)
}

var testTmpl = `Hello {{fullname}},

You are {{age}} years old.

Items:
{{#items}}
- {{key}}
{{/items}}
{{#subscribe}}You have chosen to subscribe.{{/subscribe}}

Bye.`

func loadTemplate(templateId string) (string, error) {
	return testTmpl, nil
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
