package queue

import (
  "bufio"
  "encoding/json"
  "fmt"
  "log"
  "os"

  "bitbucket.org/andoco/gomailservice/common"
  "bitbucket.org/andoco/gomailservice/models"
  "github.com/kelseyhightower/envconfig"
)

type FileMailEnqueuer struct {
}

type FileMailEnqueuerSpec struct {
  DropFolder string
}


func (q FileMailEnqueuer) Enqueue(msg *models.MailMessage) {
  fmt.Printf("queueuing message %s\n", msg.Id)

  var spec FileMailEnqueuerSpec

  err := envconfig.Process("ANDOCO_MAILSERVICE", &spec)
  if err != nil {
    log.Fatal(err)
  }

  log.Printf("Drop path is set to: %v\n", spec.DropFolder)

  filename := fmt.Sprintf("%s/mailmsg-%s", spec.DropFolder, msg.Id)

  if _, err := os.Stat(filename); os.IsNotExist(err) {

    f, err := os.Create(filename)
    common.Check(err)
    defer f.Close()

    w := bufio.NewWriter(f)

    encoder := json.NewEncoder(w)
    encoder.Encode(msg)

    w.Flush()

    /*bytes := []byte(msg.Message)
    err := ioutil.WriteFile(filename, bytes, 0644)
    check(err)*/
  }
}
