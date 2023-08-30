package msn

import (
	"github.com/soderasen-au/go-common/util"
	"testing"
)

var servers = []struct {
	name   string
	server EmailServerConfig
	want   *util.Result
}{
	{
		name: "gmail",
		server: EmailServerConfig{
			ServerType: "smtp",
			Host:       util.Ptr("smtp.gmail.com"),
			Port:       util.Ptr(587),
			Username:   util.Ptr("soderasen.au@gmail.com"),
			Password:   util.Ptr("rscmpkbnogheveft"), //util.Ptr("#5gpE9?M5TLNdJP^"),
			Encryption: util.Ptr(STARTTLS),
		},
		want: nil,
	},
	{
		name: "InvalidServerType",
		server: EmailServerConfig{
			ServerType: "smtt",
			Host:       util.Ptr("smtp.gmail.com"),
			Port:       util.Ptr(587),
			Username:   util.Ptr("soderasen.au@gmail.com"),
			Password:   util.Ptr("#5gpE9?M5TLNdJP^"),
			Encryption: util.Ptr(STARTTLS),
		},
		want: util.MsgError("EmailServerConfigCheck", "invalid server type: smtt"),
	},
	{
		name: "outlook",
		server: EmailServerConfig{
			ServerType: "smtp",
			Host:       util.Ptr("smtp.office365.com"),
			Port:       util.Ptr(587),
			Username:   util.Ptr("soderasen.au@outlook.com"),
			Password:   util.Ptr("tyhcpchmicjjzmes"),
			Encryption: util.Ptr(STARTTLS),
		},
		want: util.MsgError("EmailServerConfigCheck", "invalid server type: smtt"),
	},
}

//tyhcpchmicjjzmes

func TestEmailServerConfig_Validate(t *testing.T) {
	type fields struct {
		ServerType EmailServerType
		Host       *string
		Port       *int
		Username   *string
		Password   *string
		Encryption *EmailServerEncryption
	}
	for _, tt := range servers {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.server
			if got := s.Validate(); !util.SameResult(got, tt.want) {
				t.Errorf("Validate() = `%s`, want `%s`", util.JsonStr(got), util.JsonStr(tt.want))
			}
		})
	}
}

func TestMailer_Send(t *testing.T) {
	msg := Message{
		From:  "",
		To:    []string{"jim.z.shi@gmail.com", "soderasen.au@gmail.com"},
		Cc:    []string{"jim.z.shi@outlook.com"},
		Bcc:   nil,
		Title: "Test mail from gomail",
		Body:  "This is test mail body\n\ncheers,\nSA",
	}

	mailer, res := NewMailer(servers[2].server)
	if res != nil {
		t.Errorf("NewMailer: %v", res.Error())
	}

	res = mailer.Send(msg)
	if res != nil {
		t.Errorf("Send: %v", res.Error())
	}
}
