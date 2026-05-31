package mailbox

import (
	"fmt"
	htmltemplate "html/template"
	"net/smtp"
	texttemplate "text/template"

	"github.com/alvii147/lobo/errfmt"
)

// smtpClient implements a Client that sends email through SMTP server.
// This should typically be used in production.
type smtpClient struct {
	hostname     string
	addr         string
	username     string
	password     string
	timeProvider TimeProvider
}

// NewSMTPClient returns a new smtpClient.
func NewSMTPClient(
	hostname string,
	port int,
	username string,
	password string,
	timeProvider TimeProvider,
) *smtpClient {
	return &smtpClient{
		hostname:     hostname,
		addr:         fmt.Sprintf("%s:%d", hostname, port),
		username:     username,
		password:     password,
		timeProvider: timeProvider,
	}
}

// Send sends an email through SMTP server.
func (client *smtpClient) Send(
	to []string,
	subject string,
	textTmpl *texttemplate.Template,
	htmlTmpl *htmltemplate.Template,
	tmplData any,
) error {
	msg, err := BuildMail(
		client.username,
		to,
		subject,
		textTmpl,
		htmlTmpl,
		tmplData,
		client.timeProvider,
	)
	if err != nil {
		return errfmt.Error(err)
	}

	auth := smtp.PlainAuth("", client.username, client.password, client.hostname)
	err = smtp.SendMail(client.addr, auth, client.username, to, msg)
	if err != nil {
		return errfmt.Error(err, "smtp.SendMail failed")
	}

	return nil
}
