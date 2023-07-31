package msn

import (
	"github.com/soderasen-au/go-common/util"
	"strings"
)

type EmailServerType string

type EmailServerEncryption string

const (
	SMTP EmailServerType = "smtp"
	//mailgun, mailchimp,sendgrid,sparkpost etc.

	NoEncryption EmailServerEncryption = "none"
	SSLTLS       EmailServerEncryption = "ssl/tls"
	STARTTLS     EmailServerEncryption = "starttls"
)

type EmailServerConfig struct {
	ServerType EmailServerType        `json:"server_type,omitempty" yaml:"server_type"`
	Host       *string                `json:"host,omitempty" yaml:"host"`
	Port       *int                   `json:"port,omitempty" yaml:"port"`
	Username   *string                `json:"username,omitempty" yaml:"username"`
	Password   *string                `json:"password,omitempty" yaml:"password"`
	Encryption *EmailServerEncryption `json:"encryption,omitempty" yaml:"encryption"`
}

func (s EmailServerConfig) Validate() *util.Result {
	if EmailServerType(strings.ToLower(string(s.ServerType))) != SMTP {
		return util.MsgError("EmailServerConfigCheck", "invalid server type: "+string(s.ServerType))
	}
	if s.Host == nil {
		return util.MsgError("EmailServerConfigCheck", "nil host")
	}
	if s.Port == nil {
		s.Port = util.Ptr(587)
	}
	if s.Encryption == nil {
		s.Encryption = util.Ptr(NoEncryption)
	}
	if util.MaybeNil(s.Encryption) != NoEncryption {
		if s.Username == nil || s.Password == nil {
			return util.MsgError("EmailServerConfigCheck", "invalid username or passowrd")
		}
	}
	return nil
}
