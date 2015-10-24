package queue

import (
  "log"
  "time"

  "github.com/andoco/mail-service/models"
)

var queueChannel chan *models.MailMessage
var done chan bool
var enqueuer MailEnqueuer
var dequeuer MailDequeuer
var listeners []chan *models.MailMessage

func Enqueue(msg *models.MailMessage) {
  queueChannel <- msg
}

func Listen() chan *models.MailMessage {
  c := make(chan *models.MailMessage)
  listeners = append(listeners, c)

  return c
}

func Start() {
  queueChannel = make(chan *models.MailMessage)
  done = make(chan bool)
  enqueuer = newEnqueuer()
  dequeuer = newDequeuer()
  listeners = []chan *models.MailMessage{}
  go process(queueChannel, done, enqueuer)
  go processDequeue()
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

func processDequeue() {
  log.Print("Processing queue")
  for {
    msg := dequeuer.Dequeue()
    if msg != nil {
      log.Printf("Processing message %s", msg.Id)
    }

    time.Sleep(500 * time.Millisecond)
  }
}

func newEnqueuer() MailEnqueuer {
  return FileMailEnqueuer{}
}

func newDequeuer() MailDequeuer {
  return FileMailDequeuer{}
}
type MailEnqueuer interface {
  Enqueue(msg *models.MailMessage)
}

type MailDequeuer interface {
  Dequeue() *models.MailMessage
}
