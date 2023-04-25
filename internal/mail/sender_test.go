package mail

import (
	"testing"

	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/stretchr/testify/require"
)

func TestGmailSender_SendGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("app", "../..")
	require.NoError(t, err)

	type fields struct {
		name              string
		fromEmailAddress  string
		fromEmailPassword string
	}
	type args struct {
		subject           string
		content           string
		to                []string
		cc                []string
		bcc               []string
		attachment        []string
		smtpAuthAddress   string
		smtpServerAddress string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: 200 ok
		{
			name: "200 ok",
			fields: fields{
				name:              config.EmailSenderName,
				fromEmailAddress:  config.EmailSenderAddress,
				fromEmailPassword: config.EmailSenderPassword,
			},
			args: args{
				subject:           "Test Email",
				content:           "This is a test email",
				to:                []string{"wahyuajisulaiman011@gmail.com"},
				cc:                nil,
				bcc:               nil,
				attachment:        []string{"../../readme.md"},
				smtpAuthAddress:   "smtp.gmail.com",
				smtpServerAddress: "smtp.gmail.com:587",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gSender := &GmailSender{
				name:              tt.fields.name,
				fromEmailAddress:  tt.fields.fromEmailAddress,
				fromEmailPassword: tt.fields.fromEmailPassword,
			}
			err := gSender.SendGmail(tt.args.subject, tt.args.content, tt.args.to, tt.args.cc, tt.args.bcc, tt.args.attachment, tt.args.smtpAuthAddress, tt.args.smtpServerAddress)
			require.NoError(t, err)
		})
	}
}
