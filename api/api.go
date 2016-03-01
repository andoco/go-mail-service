package api

import "bitbucket.org/andoco/gomailservice/delivery"

type MailMessageResource struct {
	Msg *delivery.MailMessage
}
