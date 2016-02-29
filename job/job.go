package job

import (
	"fmt"
	"log"

	"bitbucket.org/andoco/gomailservice/delivery"
	"bitbucket.org/andoco/gomailservice/models"
	"bitbucket.org/andoco/gomailservice/template"
)

type Job struct {
	TemplateId string                 `json:"templateId"`
	Fields     map[string]interface{} `json:"fields"`
	To         []string               `json:"to"`
	Cc         []string               `json:"cc"`
	Bcc        []string               `json:"bcc"`
	From       string                 `json:"from"`
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

	// Deliver
	if err := delivery.Deliver(msg); err != nil {
		return fmt.Errorf("error processing job; %v", err)
	}

	return nil
}
