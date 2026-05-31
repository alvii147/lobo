package mailbox

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"io"
	"net/mail"
	"regexp"
	"strings"
	texttemplate "text/template"
	"unicode/utf8"

	"github.com/alvii147/lobo/errfmt"
	"github.com/alvii147/lobo/isotope"
)

// mail client types.
const (
	// ClientTypeSMTP represents mail clients using SMTP protocol to send emails.
	ClientTypeSMTP = "smtp"
	// ClientTypeSMTP represents mail clients using console to print emails.
	ClientTypeConsole = "console"
	// BoundaryLength represents mail boundary length.
	BoundaryLength = 32
)

// BuildMail builds multi-line email body using MIME format.
func BuildMail(
	from string,
	to []string,
	subject string,
	textTmpl *texttemplate.Template,
	htmlTmpl *htmltemplate.Template,
	tmplData any,
	timeProvider TimeProvider,
) ([]byte, error) {
	boundary := isotope.String(BoundaryLength, true, true, true)

	var mailBody bytes.Buffer
	var err error
	_, err = fmt.Fprintf(&mailBody, "Content-Type: multipart/alternative; boundary=\"%s\"\n", boundary)
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprintf failed")
	}

	_, err = fmt.Fprintf(&mailBody, "Content-Type: multipart/alternative; boundary=\"%s\"\n", boundary)
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprintf failed")
	}

	_, err = fmt.Fprint(&mailBody, "MIME-Version: 1.0\n")
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprint failed")
	}

	_, err = fmt.Fprintf(&mailBody, "Subject: %s\n", subject)
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprintf failed")
	}

	_, err = fmt.Fprintf(&mailBody, "From: %s\n", from)
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprintf failed")
	}

	_, err = fmt.Fprintf(&mailBody, "From: %s\n", strings.Join(to, ", "))
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprintf failed")
	}

	_, err = fmt.Fprintf(&mailBody, "Date: %s\n", timeProvider.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700"))
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprintf failed")
	}

	_, err = fmt.Fprint(&mailBody, "\n")
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprint failed")
	}

	_, err = fmt.Fprintf(&mailBody, "--%s\n", boundary)
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprintf failed")
	}

	_, err = fmt.Fprint(&mailBody, "Content-Type: text/plain; charset=\"utf-8\"\n")
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprint failed")
	}

	_, err = fmt.Fprint(&mailBody, "MIME-Version: 1.0\n")
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprint failed")
	}

	_, err = fmt.Fprint(&mailBody, "\n")
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprint failed")
	}

	err = textTmpl.Execute(&mailBody, tmplData)
	if err != nil {
		return nil, errfmt.Error(err, "textTmpl.Execute failed")
	}

	_, err = fmt.Fprint(&mailBody, "\n")
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprint failed")
	}

	_, err = fmt.Fprintf(&mailBody, "--%s\n", boundary)
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprintf failed")
	}

	_, err = fmt.Fprint(&mailBody, "Content-Type: text/html; charset=\"utf-8\"\n")
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprint failed")
	}

	_, err = fmt.Fprint(&mailBody, "MIME-Version: 1.0\n")
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprint failed")
	}

	_, err = fmt.Fprint(&mailBody, "\n")
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprint failed")
	}

	err = htmlTmpl.Execute(&mailBody, tmplData)
	if err != nil {
		return nil, errfmt.Error(err, "htmlTmpl.Execute failed")
	}

	_, err = fmt.Fprint(&mailBody, "\n")
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprint failed")
	}

	_, err = fmt.Fprintf(&mailBody, "--%s--\n", boundary)
	if err != nil {
		return nil, errfmt.Error(err, "fmt.Fprintf failed")
	}

	msg := mailBody.Bytes()

	return msg, nil
}

// ParseMail parses a given mail body.
func ParseMail(msg string) (string, string, error) {
	mailMsg, err := mail.ReadMessage(strings.NewReader(msg))
	if err != nil {
		return "", "", errfmt.Error(err, "mail.ReadMessage failed")
	}

	r := regexp.MustCompile(`^multipart/alternative;\s*boundary="(\S+)"$`)
	matches := r.FindStringSubmatch(mailMsg.Header.Get("Content-Type"))
	if len(matches) != 2 {
		return "", "", errfmt.Error(nil, "MustParseMailMessage failed to find content type")
	}

	boundary := matches[1]

	msgBytes, err := io.ReadAll(mailMsg.Body)
	if err != nil {
		return "", "", errfmt.Error(err, "io.ReadAll failed")
	}

	r = regexp.MustCompile(`--+` + boundary + `-*`)
	mailSections := r.Split(string(msgBytes), -1)

	nonEmptyMailSections := make([]string, 0)
	for _, sec := range mailSections {
		if utf8.RuneCountInString(strings.TrimSpace(sec)) != 0 {
			nonEmptyMailSections = append(nonEmptyMailSections, sec)
		}
	}

	if len(nonEmptyMailSections) != 2 {
		return "", "", errfmt.Error(nil, "parsing text and html sections failed")
	}

	textMsg := nonEmptyMailSections[0]
	htmlMsg := nonEmptyMailSections[1]

	return textMsg, htmlMsg, nil
}
