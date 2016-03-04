package job

import (
	"fmt"
	"log"

	"bitbucket.org/andoco/gomailservice/delivery"
	"bitbucket.org/andoco/gomailservice/template"
)

var pipelines = make(map[string]Pipeline)

type Job struct {
	TemplateId string                 `json:"templateId"`
	Fields     map[string]interface{} `json:"fields"`
	To         []string               `json:"to"`
	Cc         []string               `json:"cc"`
	Bcc        []string               `json:"bcc"`
	From       string                 `json:"from"`
}

type Pipeline struct {
	Id    string
	Steps []Step
}

type JobState struct {
	job Job
	msg delivery.MailMessage
}

type Step interface {
	Process(state JobState) error
}

/*
 * RenderStep
 */

type RenderStep struct {
	TemplateId string
}

func (rs RenderStep) Process(state JobState) error {
	return nil
}

func AddPipeline(pipeline Pipeline) {
	pipelines[pipeline.Id] = pipeline
}

func Process(job Job) error {
	log.Printf("Processing job")

	// render
	rendered, err := template.Render(job.TemplateId, job.Fields)
	if err != nil {
		return fmt.Errorf("error processing job; %v", err)
	}

	// Construct MailMessage
	msg := delivery.MailMessage{
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
