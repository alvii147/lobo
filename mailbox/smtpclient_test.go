package mailbox_test

import (
	"testing"

	"github.com/alvii147/lobo/isotope"
	"github.com/alvii147/lobo/mailbox"
	"github.com/alvii147/lobo/timekeeper"
)

func TestNewSMTPMailClient(t *testing.T) {
	t.Parallel()

	hostname := isotope.String(12, true, true, true)
	port := 587
	username := isotope.Email()
	password := isotope.Password()
	timeProvider := timekeeper.NewFrozenProvider()

	mailbox.NewSMTPClient(hostname, port, username, password, timeProvider)
}
