package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/andoco/gomailservice/api"
	"bitbucket.org/andoco/gomailservice/delivery"
	"bitbucket.org/andoco/gomailservice/job"
	"bitbucket.org/andoco/gomailservice/queue"

	"github.com/goji/param"
	"github.com/satori/go.uuid"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello asdf, %s!", c.URLParams["name"])
}

func postMail(c web.C, w http.ResponseWriter, r *http.Request) {
	var msg delivery.MailMessage

	r.ParseForm()
	err := param.Parse(r.Form, &msg)

	if err != nil || len(msg.Message) > 140 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg.Id = uuid.NewV1().String()

	queue.Enqueue(&msg)

	resource := api.MailMessageResource{Msg: &msg}

	encoder := json.NewEncoder(w)
	encoder.Encode(resource)
}

func postJob(c web.C, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var mailjob job.Job
	if err := decoder.Decode(&mailjob); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := job.Process(mailjob); err != nil {
		log.Printf("error handling postJob; %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	queue.Start()

	goji.Get("/hello/:name", hello)
	goji.Post("/mail", postMail)
	goji.Post("/job", postJob)

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
