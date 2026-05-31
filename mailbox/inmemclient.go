package mailbox

import (
	htmltemplate "html/template"
	texttemplate "text/template"
	"time"

	"github.com/alvii147/lobo/errfmt"
)

// inMemMailLogEntry represents an in-memory entry of an email event.
type inMemMailLogEntry struct {
	From    string
	To      []string
	Subject string
	Message []byte
	SentAt  time.Time
}

// inMemMailClient implements a Client that saves email data in local memory.
// This should typically be used in unit tests.
type inMemMailClient struct {
	username     string
	timeProvider TimeProvider
	Logs         []inMemMailLogEntry
}

// NewInMemMailClient returns a new inMemMailClient.
func NewInMemMailClient(username string, timeProvider TimeProvider) *inMemMailClient {
	return &inMemMailClient{
		username:     username,
		timeProvider: timeProvider,
		Logs:         make([]inMemMailLogEntry, 0),
	}
}

// Send adds an email event to in-memory storage.
func (client *inMemMailClient) Send(
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

	client.Logs = append(
		client.Logs,
		inMemMailLogEntry{
			From:    client.username,
			To:      to,
			Subject: subject,
			Message: msg,
			SentAt:  client.timeProvider.Now(),
		},
	)

	return nil
}
