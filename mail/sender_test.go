package mail

import (
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	context := `
	<h1> Hello World </h1>
	<p> This is a test email from <a href="https://congjourney.com">Cong Journey</a></p>
	`
	to := []string{"ltcong1411@gmail.com"}
	attachFiles := []string{"../doc/swagger/simple_bank.swagger.json"}

	err = sender.SendEmail(subject, context, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
