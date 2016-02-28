package delivery

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"bitbucket.org/andoco/gomailservice/models"
	"github.com/kelseyhightower/envconfig"
)

var sender MailSender

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

	body := []byte("To: " + strings.Join(msg.To, ", ") + "\r\nSubject: " + msg.Subject + "\r\n\r\n" + msg.Message)

	err = smtp.SendMail(
		address,
		auth,
		spec.User,
		msg.To,
		body,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func Deliver(msg models.MailMessage) error {
	sender.Send(&msg)
	// TODO: error handling
	return nil
}

func init() {
	sender = SmtpMailSender{}
}
