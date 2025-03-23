package mail

import (
	"testing"

	"github.com/liang3030/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test in short mode.")
	}
	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	sender := NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderPassword)

	subject := "Hello"
	content := `
		<h1>Hi there</h1>
		<p>This is a test email from your bank.</p>`
	to := []string{"zhliang0204@gmail.com"}
	attchFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attchFiles)
	require.NoError(t, err)

}
