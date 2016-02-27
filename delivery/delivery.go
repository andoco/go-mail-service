package delivery

import (
	"fmt"
	"log"
	"net/smtp"

	"bitbucket.org/andoco/gomailservice/models"
	"github.com/kelseyhightower/envconfig"
)

type SmtpMailSenderSpec struct {
	Server string
	User   string
	Pwd    string
}

type MailSender interface {
	Send(msg *models.MailMessage)
}

type FakeMailSender struct {
}

func (s FakeMailSender) Send(msg *models.MailMessage) {
	log.Printf("Sending message %s", msg.Id)
}

type SmtpMailSender struct {
}

func (s SmtpMailSender) Send(msg *models.MailMessage) {
  log.Printf("Sending message %s", msg.Id)

	var spec SmtpMailSenderSpec

	err := envconfig.Process("ANDOCO_MAILSERVICE_SMTP", &spec)
	if err != nil {
		log.Fatal(err)
	}

  auth := smtp.PlainAuth("", spec.User, spec.Pwd, spec.Server)

  address := fmt.Sprintf("%v:%v", spec.Server, 587)

  body := []byte("To: " + msg.To + "\r\nSubject: " + msg.Subject + "\r\n\r\n" + msg.Message)

  err = smtp.SendMail(
		address,
		auth,
		spec.User,
    []string{msg.To},
		body,
	)
	if err != nil {
    log.Fatal(err)
	}
}
