package job

import (
	"fmt"
	"log"

	"bitbucket.org/andoco/gomailservice/delivery"
	"bitbucket.org/andoco/gomailservice/template"
)

var pipelines = make(map[string]Pipeline)

type Job struct {
	PipelineId string                 `json:"pipelineId"`
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

func AddPipeline(pipeline Pipeline) {
	pipelines[pipeline.Id] = pipeline
}

func NewPipeline(id string, steps []Step) Pipeline {
	pipeline := Pipeline{Id: id, Steps: steps}
	return pipeline
}

type JobState struct {
	job Job
	msg delivery.MailMessage
}

type Step interface {
	Process(state *JobState) error
}

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

func Process(job Job) error {
	log.Printf("Processing job")

	state := JobState{job: job}

	pipeline, ok := pipelines[job.PipelineId]

	if !ok {
		return fmt.Errorf("could not find a pipeline with the id %s", job.PipelineId)
	}

	for _, step := range pipeline.Steps {
		if err := step.Process(&state); err != nil {
			return fmt.Errorf("error processing pipeline step %v; %v", step, err)
		}
	}

	return nil
}
