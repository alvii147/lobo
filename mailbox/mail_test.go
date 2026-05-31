package mailbox_test

import (
	htmltemplate "html/template"
	"testing"
	texttemplate "text/template"

	"github.com/alvii147/lobo/isotope"
	"github.com/alvii147/lobo/mailbox"
	"github.com/alvii147/lobo/timekeeper"
	"github.com/stretchr/testify/require"
)

func TestBuildMailSuccess(t *testing.T) {
	t.Parallel()

	from := isotope.Email()
	to := isotope.Email()
	subject := isotope.String(12, true, true, true)
	textTmpl, err := texttemplate.New("textTmpl").Parse("Test Template Content: {{ .Value }}")
	require.NoError(t, err)
	htmlTmpl, err := htmltemplate.New("htmlTmpl").Parse("<div>Test Template Content: {{ .Value }}</div>")
	require.NoError(t, err)
	tmplData := map[string]int{
		"Value": 42,
	}
	timeProvider := timekeeper.NewFrozenProvider()

	msg, err := mailbox.BuildMail(
		from,
		[]string{to},
		subject,
		textTmpl,
		htmlTmpl,
		tmplData,
		timeProvider,
	)
	require.NoError(t, err)

	textMsg, htmlMsg, err := mailbox.ParseMail(string(msg))
	require.NoError(t, err)

	require.Regexp(t, `Content-Type:\s*text\/plain;\s*charset\s*=\s*"utf-8"`, textMsg)
	require.Regexp(t, `MIME-Version:\s*1.0`, textMsg)
	require.Contains(t, textMsg, "Test Template Content: 42")

	require.Regexp(t, `Content-Type:\s*text\/html;\s*charset\s*=\s*"utf-8"`, htmlMsg)
	require.Regexp(t, `MIME-Version:\s*1.0`, htmlMsg)
	require.Contains(t, htmlMsg, "Test Template Content: 42")
}

func TestBuildMailTextTemplateError(t *testing.T) {
	t.Parallel()

	from := isotope.Email()
	to := isotope.Email()
	subject := isotope.String(12, true, true, true)
	textTmpl, err := texttemplate.New("textTmpl").Parse("Test Template Content: {{ .Value }}")
	require.NoError(t, err)
	htmlTmpl, err := htmltemplate.New("htmlTmpl").Parse("<div>Test Template Content: {{ .Value }}</div>")
	require.NoError(t, err)
	tmplData := map[string]int{
		"Value": 42,
	}
	timeProvider := timekeeper.NewFrozenProvider()

	textTmpl.Tree = nil

	_, err = mailbox.BuildMail(
		from,
		[]string{to},
		subject,
		textTmpl,
		htmlTmpl,
		tmplData,
		timeProvider,
	)
	require.Error(t, err)
}

func TestBuildMailHTMLTemplateError(t *testing.T) {
	t.Parallel()

	from := isotope.Email()
	to := isotope.Email()
	subject := isotope.String(12, true, true, true)
	textTmpl, err := texttemplate.New("textTmpl").Parse("Test Template Content: {{ .Value }}")
	require.NoError(t, err)
	htmlTmpl, err := htmltemplate.New("htmlTmpl").Parse("<div>Test Template Content: {{ .Value }}</div>")
	require.NoError(t, err)
	tmplData := map[string]int{
		"Value": 42,
	}
	timeProvider := timekeeper.NewFrozenProvider()

	htmlTmpl.Tree = nil

	_, err = mailbox.BuildMail(
		from,
		[]string{to},
		subject,
		textTmpl,
		htmlTmpl,
		tmplData,
		timeProvider,
	)
	require.Error(t, err)
}

func TestMustParseMailMessageSuccess(t *testing.T) {
	t.Parallel()

	msg := `Content-Type: multipart/alternative; boundary="deadbeef"
MIME-Version: 1.0
Subject: 1L9cBLMEBzSn
From: vfgtd7ujt535@ucgufkizok.bih
From: y32y4v6iyx5i@lyijjmvasg.tcn
Date: Thu, 25 Jan 2024 14:11:10 +0000

--deadbeef
Content-Type: text/plain; charset="utf-8"
MIME-Version: 1.0

Text Message
--deadbeef
Content-Type: text/html; charset="utf-8"
MIME-Version: 1.0

HTML Message
--deadbeef--
	`

	textMsg, htmlMsg, err := mailbox.ParseMail(msg)
	require.NoError(t, err)

	require.Contains(t, textMsg, "Text Message")
	require.NotContains(t, textMsg, "HTML Message")
	require.Contains(t, htmlMsg, "HTML Message")
	require.NotContains(t, htmlMsg, "Text Message")
}

func TestMustParseMailMessageError(t *testing.T) {
	t.Parallel()

	invalidMsg := "1nv4l1d m3554g3"
	msgWithInvalidContentType := `Content-Type: invalid/type; boundary="deadbeef"
MIME-Version: 1.0
Subject: 1L9cBLMEBzSn
From: vfgtd7ujt535@ucgufkizok.bih
From: y32y4v6iyx5i@lyijjmvasg.tcn
Date: Thu, 25 Jan 2024 14:11:10 +0000

--deadbeef
Content-Type: text/plain; charset="utf-8"
MIME-Version: 1.0

Text Message
--deadbeef
Content-Type: text/html; charset="utf-8"
MIME-Version: 1.0

HTML Message
--deadbeef--
	`
	msgWithOneSection := `Content-Type: multipart/alternative; boundary="deadbeef"
MIME-Version: 1.0
Subject: 1L9cBLMEBzSn
From: vfgtd7ujt535@ucgufkizok.bih
From: y32y4v6iyx5i@lyijjmvasg.tcn
Date: Thu, 25 Jan 2024 14:11:10 +0000

--deadbeef
Content-Type: text/plain; charset="utf-8"
MIME-Version: 1.0

Text Message
--deadbeef
	`

	testcases := map[string]struct {
		msg string
	}{
		"Invalid message": {
			msg: invalidMsg,
		},
		"Message with invalid content type": {
			msg: msgWithInvalidContentType,
		},
		"Message with one section": {
			msg: msgWithOneSection,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, _, err := mailbox.ParseMail(testcase.msg)
			require.Error(t, err)
		})
	}
}
