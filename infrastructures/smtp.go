package infrastructures

import (
	"github.com/InkForge/Blog_Website/domain"
	"gopkg.in/gomail.v2"
)

type ISMTPDialer interface {
	DialAndSend(...*gomail.Message) error
}

type SMTPService struct {
	dialer    ISMTPDialer
	EmailFrom string
}

// NetSMTPService expects configuration settings needed for settings
// up SMTP services and returns a reference to SMTPService object
func NewSMTPService(SMTPHost string, SMTPPort int, SMTPUsername, SMTPPassowrd string, EmailFrom string) domain.INotificationService {
	d := gomail.NewDialer(SMTPHost, SMTPPort, SMTPUsername, SMTPPassowrd)
	return &SMTPService{
		dialer:    d,
		EmailFrom: EmailFrom,
	}
}

// SendEmail method sends email to the specified user with subject and body.
// to -> the receiver email address
// subject -> subject of the email
// body -> html body of the email content
func (s *SMTPService) SendEmail(to, subject, body string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", s.EmailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}
