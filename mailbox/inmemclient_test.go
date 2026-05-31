package mailbox_test

import (
	htmltemplate "html/template"
	"testing"
	texttemplate "text/template"
	"time"

	"github.com/alvii147/lobo/isotope"
	"github.com/alvii147/lobo/mailbox"
	"github.com/alvii147/lobo/timekeeper"
	"github.com/stretchr/testify/require"
)

func TestInMemMailClientSendSuccess(t *testing.T) {
	t.Parallel()

	username := isotope.Email()
	timeProvider := timekeeper.NewFrozenProvider()
	client := mailbox.NewInMemMailClient(username, timeProvider)

	mailCount := len(client.Logs)

	to := isotope.Email()
	subject := isotope.String(12, true, true, true)
	textTmpl, err := texttemplate.New("textTmpl").Parse("Test Template Content: {{ .Value }}")
	require.NoError(t, err)
	htmlTmpl, err := htmltemplate.New("htmlTmpl").Parse("<div>Test Template Content: {{ .Value }}</div>")
	require.NoError(t, err)
	tmplData := map[string]int{
		"Value": 42,
	}

	err = client.Send([]string{to}, subject, textTmpl, htmlTmpl, tmplData)
	require.NoError(t, err)

	require.Len(t, client.Logs, mailCount+1)

	lastMail := client.Logs[len(client.Logs)-1]
	require.Equal(t, username, lastMail.From)
	require.Equal(t, []string{to}, lastMail.To)
	require.Equal(t, subject, lastMail.Subject)
	require.WithinDuration(t, timeProvider.Now(), lastMail.SentAt, time.Minute)

	textMsg, htmlMsg, err := mailbox.ParseMail(string(lastMail.Message))
	require.NoError(t, err)

	require.Regexp(t, `Content-Type:\s*text\/plain;\s*charset\s*=\s*"utf-8"`, textMsg)
	require.Regexp(t, `MIME-Version:\s*1.0`, textMsg)
	require.Contains(t, textMsg, "Test Template Content: 42")

	require.Regexp(t, `Content-Type:\s*text\/html;\s*charset\s*=\s*"utf-8"`, htmlMsg)
	require.Regexp(t, `MIME-Version:\s*1.0`, htmlMsg)
	require.Contains(t, htmlMsg, "Test Template Content: 42")
}

func TestInMemMailClientSendError(t *testing.T) {
	t.Parallel()

	username := isotope.Email()
	timeProvider := timekeeper.NewFrozenProvider()
	client := mailbox.NewInMemMailClient(username, timeProvider)

	mailCount := len(client.Logs)

	to := isotope.Email()
	subject := isotope.String(12, true, true, true)
	textTmpl, err := texttemplate.New("textTmpl").Parse("Test Template Content: {{ .Value }}")
	require.NoError(t, err)
	htmlTmpl, err := htmltemplate.New("htmlTmpl").Parse("<div>Test Template Content: {{ .Value }}</div>")
	require.NoError(t, err)
	tmplData := map[string]int{
		"Value": 42,
	}

	textTmpl.Tree = nil

	err = client.Send([]string{to}, subject, textTmpl, htmlTmpl, tmplData)
	require.Error(t, err)
	require.Len(t, client.Logs, mailCount)
}
