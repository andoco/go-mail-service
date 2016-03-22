package api

import (
	"encoding/json"
	"log"
	"net/http"

	"bitbucket.org/andoco/gomailservice/job"
	"github.com/zenazn/goji/web"
)

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
