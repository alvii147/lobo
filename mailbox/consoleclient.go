package mailbox

import (
	"fmt"
	htmltemplate "html/template"
	"io"
	texttemplate "text/template"

	"github.com/alvii147/lobo/errfmt"
)

// consoleClient implements a Client that prints email contents to the console.
// This should typically be used in local development.
type consoleClient struct {
	username     string
	timeProvider TimeProvider
	writer       io.Writer
}

// NewConsoleClient returns a new consoleClient.
func NewConsoleClient(username string, timeProvider TimeProvider, writer io.Writer) *consoleClient {
	return &consoleClient{
		username:     username,
		timeProvider: timeProvider,
		writer:       writer,
	}
}

// Send prints email body to the console.
func (client *consoleClient) Send(
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

	_, err = fmt.Fprint(client.writer, string(msg))
	if err != nil {
		return errfmt.Error(err, "fmt.Fprint failed")
	}

	return nil
}
