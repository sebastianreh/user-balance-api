package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailService interface {
	SendEmail(to []string, subject, body string) error
}

type smtpEmailService struct {
	username string
	password string
	from     string
	to       string
	host     string
	port     string
}

func NewSMTPEmailService(username, password, from, to, host, port string) EmailService {
	return &smtpEmailService{
		from:     from,
		username: username,
		password: password,
		to:       to,
		host:     host,
		port:     port,
	}
}

func (s *smtpEmailService) SendEmail(to []string, subject, body string) error {
	if len(to) == 0 {
		to = append(to, s.to)
	}

	auth := smtp.PlainAuth("apikey", s.username, s.password, s.host)

	headers := make(map[string]string)
	headers["From"] = s.from
	headers["To"] = strings.Join(to, ",")
	headers["Subject"] = subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	return smtp.SendMail(fmt.Sprintf("%s:%s", s.host, s.port), auth, s.from, to, []byte(message))
}
