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
			Password:   util.Ptr("#5gpE9?M5TLNdJP^"),
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
}

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
