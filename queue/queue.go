package queue

import (
  "bufio"
  "encoding/json"
  "fmt"
  "log"
  "os"

  "github.com/andoco/mail-service/models"

  "github.com/kelseyhightower/envconfig"
)

var queueChannel = make(chan *models.MailMessage)
var done = make(chan bool)
var enqueuer = newEnqueuer()

func Enqueue(msg *models.MailMessage) {
  queueChannel <- msg
}

func Start() {
  go process(queueChannel, done, enqueuer)
}

func Stop() {
  log.Print("Stopping queue")
  log.Print("Closing mail queue channel")
  close(queueChannel)
  <-done
  log.Print("Stop queue complete")
}

func process(c chan *models.MailMessage, done chan bool, enqueuer MailEnqueuer) {
  for msg := range c {
    enqueuer.Enqueue(msg)
  }

  log.Print("Finished processing mail queue channel")
  done <- true
}

func newEnqueuer() MailEnqueuer {
  return FileMailEnqueuer{}
}

type MailEnqueuer interface {
  Enqueue(msg *models.MailMessage)
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

func (q FileMailEnqueuer) Enqueue(msg *models.MailMessage) {
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
