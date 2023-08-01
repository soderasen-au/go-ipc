package msn

import (
	"crypto/tls"
	"github.com/soderasen-au/go-common/util"
	gomail "github.com/xhit/go-simple-mail/v2"
	"strings"
	"time"
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

func (s *EmailServerConfig) Validate() *util.Result {
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
	enc := strings.ToLower(string(*s.Encryption))
	s.Encryption = util.Ptr(EmailServerEncryption(enc))
	if util.MaybeNil(s.Encryption) != NoEncryption {
		if s.Username == nil || s.Password == nil {
			return util.MsgError("EmailServerConfigCheck", "invalid username or passowrd")
		}
	}
	return nil
}

type Mailer struct {
	Config EmailServerConfig
	client *gomail.SMTPClient
}

func (m Mailer) Send(msg Message) *util.Result {
	email := gomail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To...).AddCc(msg.Cc...).AddBcc(msg.Bcc...).SetSubject(msg.Title).SetBody(gomail.TextHTML, msg.Body)
	if err := email.Send(m.client); err != nil {
		return util.Error("SendEmail", err)
	}
	return nil
}

func (m Mailer) Name() string {
	return "EMail[" + util.MaybeNil(m.Config.Host) + "]"
}

func getEncryptionMethod(enc EmailServerEncryption) gomail.Encryption {
	switch enc {
	case NoEncryption:
		return gomail.EncryptionNone
	case SSLTLS:
		return gomail.EncryptionSSLTLS
	case STARTTLS:
		return gomail.EncryptionSTARTTLS
	}
	return gomail.EncryptionSTARTTLS
}

func NewMailer(cfg EmailServerConfig) (*Mailer, *util.Result) {
	res := cfg.Validate()
	if res != nil {
		return nil, res.With("Invalid server config")
	}

	mailer := &Mailer{
		Config: cfg,
		client: nil,
	}

	server := gomail.NewSMTPClient()
	server.Host = util.MaybeNil(cfg.Host)
	server.Port = util.MaybeNil(cfg.Port)
	server.Username = util.MaybeNil(cfg.Username)
	server.Password = util.MaybeNil(cfg.Password)
	server.Encryption = getEncryptionMethod(util.MaybeNil(cfg.Encryption))
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	smtpClient, err := server.Connect()
	if err != nil {
		return nil, util.Error("CreateSmtpClient", err)
	}
	mailer.client = smtpClient

	return mailer, nil
}
