package job

import (
	"fmt"
	"log"

	"bitbucket.org/andoco/gomailservice/delivery"
	"bitbucket.org/andoco/gomailservice/template"
)

/*
 * RenderStep
 */

type RenderStep struct {
	TemplateId string
}

func (step RenderStep) Process(state *JobState) error {
	// render
	rendered, err := template.Render(step.TemplateId, state.job.Fields)
	if err != nil {
		return fmt.Errorf("error processing job; %v", err)
	}

	state.job.Fields["renderedSubject"] = "TODO: render a subject"
	state.job.Fields["renderedContent"] = rendered

	return nil
}

/*
 * BuildMessageStep
 */

type BuildMessageStep struct {
}

func (step BuildMessageStep) Process(state *JobState) error {
	job := state.job

	subject, ok := job.Fields["renderedSubject"]
	if !ok {
		return fmt.Errorf("error retrieving rendered subject from the job state")
	}

	content, ok := job.Fields["renderedContent"]
	if !ok {
		return fmt.Errorf("error retrieving rendered content from the job state")
	}

	// Construct MailMessage
	state.msg = delivery.MailMessage{
		To:      job.To,
		Subject: subject.(string),
		Message: content.(string),
	}

	log.Printf("built message %v", state.msg)

	return nil
}

/*
 * DeliveryStep
 */

type DeliveryStep struct {
}

func (s DeliveryStep) Process(state *JobState) error {
	return delivery.Deliver(state.msg)
}
