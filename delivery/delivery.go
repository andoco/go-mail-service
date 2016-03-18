package delivery

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var sender MailSender

type MailMessage struct {
	Id      string
	To      []string  `param:"to"`
	From    string    `param:"from"`
	Message string    `param:"message"`
	Subject string    `param:"subject"`
	Time    time.Time `param:"time"`
}

type SmtpMailSenderSpec struct {
	Server string
	User   string
	Pwd    string
}

type MailSender interface {
	Send(msg *MailMessage) error
}

type FakeMailSender struct {
}

func (s FakeMailSender) Send(msg *MailMessage) error {
	log.Printf("Sending message %s", msg.Id)
	return nil
}

type SmtpMailSender struct {
}

func (s SmtpMailSender) Send(msg *MailMessage) error {
	log.Printf("Sending message %v", msg)

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

func Deliver(msg MailMessage) error {
	return sender.Send(&msg)
}

func init() {
	sender = SmtpMailSender{}
}
