package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	stmpAuthAddress   = "smtp.gmail.com"
	stmpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attchFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}

}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attchFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attchFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("could not attach file: %v", err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, stmpAuthAddress)
	return e.Send(stmpServerAddress, smtpAuth)
}
