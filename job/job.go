package job

import (
	"log"

	"bitbucket.org/andoco/gomailservice/delivery"
	"bitbucket.org/andoco/gomailservice/models"
	"bitbucket.org/andoco/gomailservice/template"
)

type Job struct {
	TemplateId string `json:"templateId"`
	Fields     map[string]string
	To         []string `json:"to"`
	Cc         string
	Bcc        string
	From       string
}

func Process(job Job) error {
	log.Printf("Processing job")

	// render
	rendered, _ := template.Render(job.TemplateId, job.Fields)

	// Construct MailMessage
	msg := models.MailMessage{
		To:      job.To,
		Subject: "test msg",
		Message: rendered,
	}

	// deliver
	log.Printf("Rendered:\n%s", rendered)
	delivery.Deliver(msg)

	return nil
}
