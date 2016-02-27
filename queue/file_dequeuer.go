package queue

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"bitbucket.org/andoco/gomailservice/models"
	"github.com/kelseyhightower/envconfig"
)

var spec FileMailEnqueuerSpec

func init() {
	err := envconfig.Process("ANDOCO_MAILSERVICE", &spec)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Drop path is set to: %v\n", spec.DropFolder)
}

type FileMailDequeuerSpec struct {
	DropFolder string
}

type FileMailDequeuer struct {
}

func (dequeuer FileMailDequeuer) Dequeue() *models.MailMessage {
	pattern := path.Join(spec.DropFolder, "mailmsg-*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}

	if matches != nil {
		for _, fn := range matches {
			log.Printf("reading %s", fn)
			dat, err := ioutil.ReadFile(fn)
			if err != nil {
				log.Fatal(err)
			}

			var msg *models.MailMessage

			if err := json.Unmarshal(dat, &msg); err != nil {
				log.Fatal(err)
			}

			os.Remove(fn)

			return msg
		}
	}

	return nil
}
