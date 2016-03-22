package job

import (
	"fmt"
	"log"
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
}

type Step interface {
	Process(state *JobState) error
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
