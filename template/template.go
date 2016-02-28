package template

import (
	"log"

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

func Render(templateId string, fields map[string]string) (string, error) {
	log.Printf("Rendering template %s", templateId)
	return "dummy rendered content", nil
}
