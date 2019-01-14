package lightmail

import (
	"fmt"
	"github.com/oxtoacart/bpool"
	template2 "html/template"
	mail2 "net/mail"
	"strings"
	"text/template"

	gomail "github.com/go-mail/mail"
	"github.com/google/uuid"
)

// HTMLMail text/plain mail instance
type HTMLMail struct {
	htmlTemplate  *template2.Template
	plainTemplate *template.Template
	config        *SMTPConfig
	htmlBody	  string
	plainBody	  string
	bpool 		  *bpool.BufferPool
}

// NewHTMLMail creates a new instance of a text/html mail
// It returns the instance
func (s *MailService) NewHTMLMail(htmlTemplate *template2.Template, plainTemplate *template.Template) *HTMLMail {
	return &HTMLMail{
		htmlTemplate:  htmlTemplate,
		plainTemplate: plainTemplate,
		config:        s.config,
		bpool: 		   s.bpool,
	}
}

// ExecuteTemplates executes the associated template with the given data
// It returns an error or nil
func (mail *HTMLMail) ExecuteTemplate(data interface{}) error {
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

	// execute template
	buffer.Reset()
	err = mail.htmlTemplate.Execute(buffer, data)

	// check for error
	if err != nil {
		return err
	}

	mail.htmlBody = buffer.String()

	return nil
}

// Send sends the executed template as text/hmtl mail
// It returns an error or nil
func (mail *HTMLMail) Send(from *mail2.Address, to []*mail2.Address, subject string) error {
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

	// set mime multipart/alternative and body
	msg.SetBody(mimePlainText, mail.plainBody)
	msg.AddAlternative(mimeHTML, mail.htmlBody)

	// Setup smtp connection
	d := gomail.NewDialer(mail.config.Server, mail.config.Port, mail.config.Username, mail.config.Password)
	d.TLSConfig = mail.config.TLSConfig

	// Send using connection
	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
