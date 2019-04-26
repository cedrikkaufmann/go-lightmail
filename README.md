# LightMail
LightMail is a simple package to send text/plain or text/html mails via SMTP using the standard html/template and text/template packages.

## Usage
In order to use this package you should be familiar with the following packages `text/template` and `html/template`.

### Mail service instance
First a mail service instance is needed. To establish a connection to the smtp server, you need to specify the smtp settings.
````
config := lightmail.SMTPConfig{
   Server:      "smtp.mailserver.tld",
   Port:        "587",
   User:        "bob@mailserver.tld",
   Password:    "foo",
   TLSConfig:   &tls.Config{
                    InsecureSkipVerify: false,
                    ServerName:         "smtp.mailserver.tld",
                },
}

mailService := lightmail.NewMailService(&config)
````

### Plain text mail instance
To send a plain text mail you need to create a new text mail instance using your text template reference.
````
func (s *MailService) NewTextMail(template *template.Template) *TextMail
````

### HTML mail instance
To send a html mail you need to create a new html mail instance using your html template reference.
````
func (s *MailService) NewHTMLMail(template *template.Template) *HTMLMail
````

### Execute Template
Next you need to execute your template.
````
func ExecuteTemplate(data interface{}) error
````

### Send mail
Now just call the send method on your mail instance to send the mail using the mail service, where mail.Address is from the default go mail package.
````
func Send(from *mail.Address, to []*mail.Address, subject string) error
````

### Dependencies
In order to work properly, this packages needs the follwoing third party packages installed in your GoPath.
- github.com/google/uuid
- github.com/oxtoacart/bpool
- github.com/go-mail/mail


## License
MIT licensed 2019 Cedrik Kaufmann. See the LICENSE file for further details.