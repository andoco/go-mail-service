package main

import (
  "encoding/json"
	"fmt"
	"net/http"

  "github.com/goji/param"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello asdf, %s!", c.URLParams["name"])
}

func postMail(c web.C, w http.ResponseWriter, r *http.Request) {
	var msg MailMessageResource

	r.ParseForm()
	err := param.Parse(r.Form, &msg)

	if err != nil || len(msg.Message) > 140 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

  encoder := json.NewEncoder(w)
  encoder.Encode(msg)
}

func main() {
	goji.Get("/hello/:name", hello)
	goji.Post("/mail", postMail)
	goji.Serve()
}
