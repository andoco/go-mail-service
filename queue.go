package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "log"
  "os"

  "github.com/kelseyhightower/envconfig"
)

type MailEnqueuer interface {
  Enqueue(msg *MailMessage)
}

type FileMailEnqueuer struct {
}

type FileMailEnqueuerSpec struct {
  DropFolder string
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func (q FileMailEnqueuer) Enqueue(msg *MailMessage) {
  fmt.Printf("queueuing message %s\n", msg.Id)

  var spec FileMailEnqueuerSpec

  err := envconfig.Process("ANDOCO_MAILSERVICE", &spec)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("Drop path is set to: %v\n", spec.DropFolder)

  filename := fmt.Sprintf("%s/mailmsg-%s", spec.DropFolder, msg.Id)

  if _, err := os.Stat(filename); os.IsNotExist(err) {

    f, err := os.Create(filename)
    check(err)
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
