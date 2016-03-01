package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/andoco/gomailservice/delivery"
	"bitbucket.org/andoco/gomailservice/job"
	"bitbucket.org/andoco/gomailservice/queue"
	"github.com/goji/param"
	"github.com/satori/go.uuid"
	"github.com/zenazn/goji/web"
)

type MailMessageResource struct {
	Msg *delivery.MailMessage
}

func Hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello asdf, %s!", c.URLParams["name"])
}

func PostMail(c web.C, w http.ResponseWriter, r *http.Request) {
	var msg delivery.MailMessage

	r.ParseForm()
	err := param.Parse(r.Form, &msg)

	if err != nil || len(msg.Message) > 140 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg.Id = uuid.NewV1().String()

	queue.Enqueue(&msg)

	resource := MailMessageResource{Msg: &msg}

	encoder := json.NewEncoder(w)
	encoder.Encode(resource)
}

func PostJob(c web.C, w http.ResponseWriter, r *http.Request) {
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
