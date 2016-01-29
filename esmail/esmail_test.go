package esmail

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMail(t *testing.T) {
	test := NewMail()
	assert.NotNil(t, test)
	assert.Equal(t, test.Header.mimeVersion, MIME_VERSION)
	assert.Equal(t, test.Header.contentType, CONTENT_TYPE)
}

func TestSetters(t *testing.T) {
	test := NewMail()
	add := []string{"john@lol.com", "kimiko@lol.com"}
	assert.NotNil(t, test)
	test.SetRecipients(add)
	assert.Equal(t, add, test.to)
	assert.Equal(t, add, test.Header.to)
	test.SetBody("Hello You")
	assert.Equal(t, "Hello You", test.body)
	test.SetBody("Hello %s of %d", "friends", 42)
	assert.Equal(t, "Hello friends of 42", test.body)
	test.AddToBody("\r\nBest regards from %s", "Paris\r\n")
	assert.Equal(t, "Hello friends of 42\r\nBest regards from Paris\r\n", test.body)
	test.ResetBody()
	assert.Equal(t, "", test.body)
	test.SetSubject("Danger !")
	assert.Equal(t, "Alert Elastic: Danger !", test.subject)
	assert.Equal(t, "Alert Elastic: Danger !", test.Header.subject)
}

func TestGetFullHeader(t *testing.T) {
	test := NewMail()
	add := []string{"john@lol.com", "kimiko@lol.com"}
	test.SetRecipients(add)
	test.SetSubject("Hello mate !")
	test.SetFrom("Roberto")

	expected := "From: Roberto" + RN +
		"To: john@lol.com, kimiko@lol.com" + RN +
		"Subject: Alert Elastic: Hello mate !" + RN +
		"MIME-Version: 1.0" + RN +
		"Content-Type: text/html; charset=\"utf-8\"" + RN + RN

	assert.Equal(t, expected, test.Header.getFullHeader())
}
