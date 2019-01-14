package lightmail

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/oxtoacart/bpool"
	mail2 "net/mail"
	"strings"
	"text/template"

	gomail "github.com/go-mail/mail"
)

// TextMail text/plain mail instance
type TextMail struct {
	plainTemplate *template.Template
	config   *SMTPConfig
	plainBody     string
	bpool 	 *bpool.BufferPool
}

// NewTextMail creates a new instance of a text/plain mail
// It returns the instance
func (s *MailService) NewTextMail(template *template.Template) *TextMail {
	return &TextMail{
		plainTemplate: template,
		config:   s.config,
		bpool: 	  s.bpool,
	}
}

// ExecuteTemplates executes the associated template with the given data
// It returns an error or nil
func (mail *TextMail) ExecuteTemplate(data interface{}) error {
	// get buffer
	buffer := mail.bpool.Get()
	defer mail.bpool.Put(buffer)

	// render plain text version
	buffer.Reset()
	err := mail.plainTemplate.Execute(buffer, data)

	// check for error
	if err != nil {
		return err
	}

	mail.plainBody = buffer.String()

	return nil
}

// Send sends the executed template as text/plain mail
// It returns an error or nil
func (mail *TextMail) Send(from *mail2.Address, to []*mail2.Address, subject string) error {
	// get buffer
	buffer := mail.bpool.Get()
	defer mail.bpool.Put(buffer)

	// create new message
	msg := gomail.NewMessage()

	// set sender
	msg.SetHeader(headerFrom, from.String())

	// build message id
	messageUuid := uuid.New()
	mailComponents := strings.Split(from.Address, "@")
	msg.SetHeader(headerMessageId, fmt.Sprintf("<%s@%s>", messageUuid.String(), mailComponents[1]))

	// set addresses
	addrLen := len(to) - 1

	for i, addr := range to {
		buffer.WriteString(addr.String())

		if i != addrLen {
			buffer.WriteString(", ")
		}
	}

	msg.SetHeader(headerTo, buffer.String())

	// set subject
	msg.SetHeader(headerSubject, subject)

	// set mime text/plain and body
	msg.SetBody(mimePlainText, mail.plainBody)

	// Setup smtp connection
	d := gomail.NewDialer(mail.config.Server, mail.config.Port, mail.config.Username, mail.config.Password)
	d.TLSConfig = mail.config.TLSConfig

	// Send using connection
	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
