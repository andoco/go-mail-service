package main

import (
	"time"
)

type MailMessage struct {
	Id			string
	To      string    `param:"to"`
	From    string    `param:"from"`
	Message string    `param:"message"`
	Time    time.Time `param:"time"`
}

type MailMessageResource struct {
	Msg *MailMessage
}
