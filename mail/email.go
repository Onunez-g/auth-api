package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type EmailSender interface {
	Send(body string) error
	Auth() smtp.Auth
}

type MailSettings struct {
	User      string
	Pass      string
	From      string
	Smtp      string
	Port      string
	To        string
	TlsConfig *tls.Config
}

func (g MailSettings) Send(subject string, body string) error {
	addr := fmt.Sprintf("%s:%s", g.Smtp, g.Port)
	msg := "From: " + g.From + "\n" +
		"To: " + g.To + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	g.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         g.Smtp,
	}
	auth := g.Auth()
	err := smtp.SendMail(addr, auth, g.From, []string{g.To}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func (g MailSettings) Auth() smtp.Auth {
	return smtp.PlainAuth("", g.User, g.Pass, g.Smtp)
}
