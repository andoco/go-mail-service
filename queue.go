package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "os"
)

type MailEnqueuer interface {
  Enqueue(msg *MailMessage)
}

type FileMailEnqueuer struct {
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func (q FileMailEnqueuer) Enqueue(msg *MailMessage) {
  fmt.Printf("queueuing message %s\n", msg.Id)

  filename := fmt.Sprintf("/tmp/mailmsg-%s", msg.Id)

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
