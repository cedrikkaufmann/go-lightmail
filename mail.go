package lightmail

import (
	"bytes"
	"crypto/tls"
	"github.com/oxtoacart/bpool"
)

const (
	mimePlainText                       = "text/plain"
	mimeHTML							= "text/html"
	headerFrom                          = "From"
	headerTo                            = "To"
	headerSubject                       = "Subject"
	headerMessageId						= "Message-Id"
)

type MailService struct {
	config *SMTPConfig
	bpool *bpool.BufferPool
}

func NewMailService(config *SMTPConfig) *MailService{
	return &MailService{
		config: config,
		bpool: bpool.NewBufferPool(48),
	}
}

type Mail interface {
	Send(from string, to []string, subject string) bool
	ExecuteTemplate(buffer *bytes.Buffer, data interface{}) error
}

type SMTPConfig struct {
	Server    string
	Port      int
	Username  string
	Password  string
	TLSConfig *tls.Config
}
