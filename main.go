package main

import (
	"bitbucket.org/andoco/gomailservice/api"
	"bitbucket.org/andoco/gomailservice/delivery"
	"bitbucket.org/andoco/gomailservice/job"
	"bitbucket.org/andoco/gomailservice/queue"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
)

func main() {
	customFields := map[string]interface{}{
		"siteUrl": "http://localhost",
	}

	steps := []job.Step{
		job.SetFieldStep{Fields: customFields},
		job.RenderStep{TemplateId: "default"},
		job.BuildMessageStep{},
		job.DeliveryStep{},
	}

	pipeline := job.NewPipeline("default", steps)
	job.AddPipeline(pipeline)

	queue.Start()

	goji.Post("/mail", api.PostMail)
	goji.Post("/job", api.PostJob)

	graceful.PostHook(func() {
		queue.Stop()
	})

	go func() {
		sender := delivery.SmtpMailSender{}
		for msg := range queue.Listen() {
			//rendered, _ := template.Render(msg)
			//log.Printf("Rendered message: %s", rendered)
			sender.Send(msg)
		}
	}()

	goji.Serve()

	graceful.Wait()
}
