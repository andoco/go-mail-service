package main

import (
	"bitbucket.org/andoco/gomailservice/api"
	"bitbucket.org/andoco/gomailservice/delivery"
	"bitbucket.org/andoco/gomailservice/queue"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
)

func main() {
	queue.Start()

	goji.Get("/hello/:name", api.Hello)
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
