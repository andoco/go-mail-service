package delivery

import (
  "log"

  "github.com/andoco/mail-service/models"
)

type MailSender interface {
  Send(msg *models.MailMessage)
}

type FakeMailSender struct {

}

func (s FakeMailSender) Send(msg *models.MailMessage) {
  log.Printf("Sending message %s", msg.Id)
}
