package queue

import (
  "encoding/json"
  "io/ioutil"
  "path/filepath"
  "log"
  "os"

  "github.com/andoco/mail-service/models"
)

type FileMailDequeuer struct {

}

func (dequeuer FileMailDequeuer) Dequeue() *models.MailMessage {
  matches, err := filepath.Glob("/tmp/mailmsg-*")
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
