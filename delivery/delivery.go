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
	Send(msg *models.MailMessage) error
}

type FakeMailSender struct {
}

func (s FakeMailSender) Send(msg *models.MailMessage) error {
	log.Printf("Sending message %s", msg.Id)
	return nil
}

type SmtpMailSender struct {
}

func (s SmtpMailSender) Send(msg *models.MailMessage) error {
	log.Printf("Sending message %s", msg.Id)

	var spec SmtpMailSenderSpec

	if err := envconfig.Process("ANDOCO_MAILSERVICE_SMTP", &spec); err != nil {
		return fmt.Errorf("could not process envconfig; %v", err)
	}

	auth := smtp.PlainAuth("", spec.User, spec.Pwd, spec.Server)

	address := fmt.Sprintf("%v:%v", spec.Server, 587)

	body := []byte("To: " + strings.Join(msg.To, ", ") + "\r\nSubject: " + msg.Subject + "\r\n\r\n" + msg.Message)

	if err := smtp.SendMail(
		address,
		auth,
		spec.User,
		msg.To,
		body,
	); err != nil {
		return fmt.Errorf("error sending mail; %v", err)
	}

	log.Printf("Sent message %s", msg.Id)

	return nil
}

func Deliver(msg models.MailMessage) error {
	return sender.Send(&msg)
}

func init() {
	sender = SmtpMailSender{}
}
