package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

// implementing interface for what email sender should have
// subject, content, to, from, cc, bcc, attachment
type EmailSender interface {
	SendEmail(subject, content string, to, cc, bcc []string, attachment []string, smtpAuthAddress, smtpServerAddress string) error
}

// typing struct for sender. then ex. platform is Gmail
// name, email, password, host, port, sslmode
type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddress, fromEmailPassword string) *GmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

// for default smtp server gmail
// "smtp.gmail.com", "smtp.gmail.com:587"
func (gSender *GmailSender) SendGmail(subject, content string, to, cc, bcc []string, attachment []string, smtpAuthAddress, smtpServerAddress string) error {
	e := email.NewEmail()

	e.From = gSender.name + "<" + gSender.fromEmailAddress + ">"
	e.Subject = subject
	//this content can be replaced with html using backtick
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, file := range attachment {
		_, err := e.AttachFile(file)
		if err != nil {
			return fmt.Errorf("error attaching file: %v", err)
		}
	}

	smtpAuth := smtp.PlainAuth("", gSender.fromEmailAddress, gSender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}
