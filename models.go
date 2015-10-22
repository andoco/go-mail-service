package main

import (
	"time"
)

type MailMessageResource struct {
	To      string    `param:"to"`
	From    string    `param:"from"`
	Message string    `param:"message"`
	Time    time.Time `param:"time"`
}
