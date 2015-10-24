package main

import (
  "encoding/json"
	"fmt"
  "log"
  "net/http"
  "github.com/andoco/mail-service/queue"
  "github.com/andoco/mail-service/models"

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
	var msg models.MailMessage

	r.ParseForm()
	err := param.Parse(r.Form, &msg)

	if err != nil || len(msg.Message) > 140 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

  msg.Id = uuid.NewV1().String()

  queue.Enqueue(&msg)

  resource := models.MailMessageResource{Msg: &msg}

  encoder := json.NewEncoder(w)
  encoder.Encode(resource)
}

func main() {
  queue.Start()

	goji.Get("/hello/:name", hello)
	goji.Post("/mail", postMail)

  graceful.PostHook(func () {
    queue.Stop()
  })

  go func() {
    for msg := range queue.Listen() {
      log.Printf("Will send msg %s", msg.Id)
    }
  }()

	goji.Serve()

  graceful.Wait()
}
