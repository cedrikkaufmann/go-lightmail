package lightmail

import (
	"crypto/tls"
	"github.com/oxtoacart/bpool"
	"net/mail"
)

const (
	mimePlainText   = "text/plain"
	mimeHTML        = "text/html"
	headerFrom      = "From"
	headerTo        = "To"
	headerSubject   = "Subject"
	headerMessageId = "Message-Id"
)

type MailService struct {
	config *SMTPConfig
	bpool  *bpool.BufferPool
}

func NewMailService(config *SMTPConfig) *MailService {
	return &MailService{
		config: config,
		bpool:  bpool.NewBufferPool(48),
	}
}

type Mail interface {
	Send(from *mail.Address, to []*mail.Address, subject string) error
	ExecuteTemplate(data interface{}) error
}

type SMTPConfig struct {
	Server    string
	Port      int
	Username  string
	Password  string
	TLSConfig *tls.Config
}
