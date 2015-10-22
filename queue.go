package main

import (
  "fmt"
)

type MailEnqueuer interface {
  Enqueue(msg *MailMessage)
}

type FileMailEnqueuer struct {
}

func (q FileMailEnqueuer) Enqueue(msg *MailMessage) {
  fmt.Printf("queueuing message %s\n", msg.Id)
}
