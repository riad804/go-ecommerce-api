package tests

import (
	"testing"

	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/mail"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := config.LoadConfig("../.")
	require.NoError(t, err)

	sender := mail.NewGmailSender(config.Email.Name, config.Email.Address, config.Email.Password)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="http://google.com">Full stack e-commerce Project</a></p>
	`
	to := []string{"developer.riad@gmail.com"}
	//attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}
