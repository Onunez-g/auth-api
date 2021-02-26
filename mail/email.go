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

type EmailServer struct {
	User      string
	Pass      string
	From      string
	Smtp      string
	Port      int
	To        string
	TlsConfig *tls.Config
}

func (g EmailServer) Send(subject string, body string) error {
	addr := fmt.Sprintf("%s:%d", g.Smtp, g.Port)
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
	// conn, err := tls.Dial("tcp", addr, g.TlsConfig)
	// if err != nil {
	// 	fmt.Println("tls Dial")
	// 	return err
	// }
	// client, err := smtp.NewClient(conn, g.Smtp)
	// if err != nil {
	// 	fmt.Println("New Client")
	// 	return err
	// }
	// if err = client.Auth(auth); err != nil {
	// 	fmt.Println("Auth")
	// 	return err
	// }
	// if err = client.Mail(g.From); err != nil {
	// 	fmt.Println("Mail")
	// 	return err
	// }
	// recievers := []string{g.To}

	// for _, k := range recievers {
	// 	fmt.Println("Sending to: " + k)
	// 	if err = client.Rcpt(k); err != nil {
	// 		return err
	// 	}
	// }
	// w, err := client.Data()
	// if err != nil {
	// 	return err
	// }
	// _, err = w.Write([]byte(msg))
	// if err != nil {
	// 	return err
	// }
	// err = w.Close()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (g EmailServer) Auth() smtp.Auth {
	return smtp.PlainAuth("", g.User, g.Pass, g.Smtp)
}
