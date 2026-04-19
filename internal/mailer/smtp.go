package mail

import (
	"back/internal/config"

	"gopkg.in/gomail.v2"
)

type SMTPMailer struct {
	host     string
	port     int
	user     string
	password string
}

func NewSMTPMailer(cfg *config.Config) *SMTPMailer {
	return &SMTPMailer{
		host:     cfg.MailerHost,
		port:     465, 
		user:     cfg.MailerUser,
		password: cfg.MailerPassword,
	}
}

func (m *SMTPMailer) Send(to, subject, code string) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", m.user)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	html := buildVerifyEmailHTML(code)

	msg.SetBody("text/html", html)

	dialer := gomail.NewDialer(m.host, m.port, m.user, m.password)

	return dialer.DialAndSend(msg)
}