package msn

import (
	"crypto/tls"
	"github.com/rs/zerolog"
	"github.com/soderasen-au/go-common/loggers"
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
	ServerType       EmailServerType        `json:"server_type,omitempty" yaml:"server_type,omitempty"`
	From             *string                `json:"from,omitempty" yaml:"from,omitempty"`
	Host             *string                `json:"host,omitempty" yaml:"host,omitempty"`
	Port             *int                   `json:"port,omitempty" yaml:"port,omitempty"`
	Username         *string                `json:"username,omitempty" yaml:"username,omitempty"`
	Password         *string                `json:"password,omitempty" yaml:"password,omitempty"`
	Encryption       *EmailServerEncryption `json:"encryption,omitempty" yaml:"encryption,omitempty"`
	KeepAlive        *bool                  `json:"keep_alive,omitempty" yaml:"keep_alive,omitempty"`
	KeepAliveTimeout *int                   `json:"keep_alive_timeout,omitempty" yaml:"keep_alive_timeout"`
	NoopTimeout      *int                   `json:"noop_timeout,omitempty" yaml:"noop_timeout,omitempty"`
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
	server     *gomail.SMTPServer
	client     *gomail.SMTPClient
	inbox      chan Message
	outbox     chan *util.Result
	stopChan   chan bool
	lastCreate time.Time
	lastSent   time.Time
	running    bool

	Config EmailServerConfig
	Logger *zerolog.Logger
}

type MailerOption func(*Mailer) *Mailer

func WithLogger(logger *zerolog.Logger) MailerOption {
	return func(mailer *Mailer) *Mailer {
		mailer.Logger = logger
		return mailer
	}
}

func NewMailer(cfg EmailServerConfig, opts ...MailerOption) (*Mailer, *util.Result) {
	res := cfg.Validate()
	if res != nil {
		return nil, res.With("Invalid server config")
	}

	mailer := &Mailer{
		server:   nil,
		client:   nil,
		inbox:    make(chan Message),
		outbox:   make(chan *util.Result),
		stopChan: make(chan bool),
		lastSent: time.Now(),
		running:  false,

		Config: cfg,
		Logger: loggers.NullLogger,
	}

	server := gomail.NewSMTPClient()
	server.Host = util.MaybeNil(cfg.Host)
	server.Port = util.MaybeNil(cfg.Port)
	server.Username = util.MaybeNil(cfg.Username)
	server.Password = util.MaybeNil(cfg.Password)
	server.Encryption = getEncryptionMethod(util.MaybeNil(cfg.Encryption))
	server.KeepAlive = util.MaybeNil(cfg.KeepAlive)
	server.ConnectTimeout = 30 * time.Second
	server.SendTimeout = 30 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	mailer.server = server

	if res := mailer.reconnect(); res != nil {
		return nil, res.With("mailer.reconnect")
	}

	if opts != nil {
		for _, opt := range opts {
			_ = opt(mailer)
		}
	}

	if mailer.IsKeepAlive() {
		go mailer.loop()
	} else {
		mailer.running = true
	}

	return mailer, nil
}

func (m *Mailer) IsKeepAlive() bool {
	return util.MaybeNil(m.Config.KeepAlive)
}

func (m *Mailer) reconnect() *util.Result {
	logger := m.Logger.With().Str("mod", "mailer").Str("func", "reconnect").Logger()
	if m.client != nil {
		err := m.client.Close()
		if err != nil {
			util.LogError(&logger, "client.Close()", err)
		}
	}
	var err error
	m.client, err = m.server.Connect()
	if err != nil {
		return util.LogError(&logger, "server.Connect()", err)
	}
	m.lastCreate = time.Now()
	logger.Info().Msg("reconnected")
	return nil
}

func (m *Mailer) loop() {
	logger := m.Logger.With().Str("mod", "mailer").Str("func", "loop").Logger()
	logger.Info().Msg("start")
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	noopTimeout := util.MaybeDefault(m.Config.NoopTimeout, 30)
	keepaliveTimeout := util.MaybeDefault(m.Config.KeepAliveTimeout, 300)
	m.running = true
	defer func() { m.running = false }()

	for {
		select {
		case _, ok := <-m.stopChan:
			if !ok {
				logger.Warn().Msgf("stopped in invalid status")
			}
			logger.Info().Msg("stopped")
			return
		case t := <-ticker.C:
			logger.Trace().Msg("noop")
			if sinceLastCreate := int(t.Sub(m.lastCreate).Seconds()); sinceLastCreate > keepaliveTimeout {
				logger.Debug().Msgf("connection was last created at: %s, more than %d seconds. reset client", m.lastCreate.Format(time.RFC1123Z), keepaliveTimeout)
				if res := m.reconnect(); res != nil {
					logger.Err(res).Msg("reconnect")
					return
				}
			}
			if sinceLastSec := int(t.Sub(m.lastSent).Seconds()); sinceLastSec > noopTimeout {
				logger.Debug().Msgf("last sent was at: %s, more than %d seconds. send noop", m.lastSent.Format(time.RFC1123Z), noopTimeout)
				err := m.client.Noop()
				if err != nil {
					logger.Err(err).Msg("Noop() failed. try to reconnect to server")
					if res := m.reconnect(); res != nil {
						logger.Err(res).Msg("reconnect")
						return
					}
				}
				m.lastSent = time.Now()
			}
		case msg, ok := <-m.inbox:
			if !ok {
				logger.Warn().Msgf("got a invalid message: %v", msg)
			} else {
				res := m.doSend(msg)
				if res != nil {
					logger.Err(res).Msg("doSend")
				}
				logger.Info().Msg("email sent, getting next message ...")
				time.Sleep(time.Duration(100) * time.Millisecond)
			}
		}
	}
}

func (m *Mailer) outResult(res *util.Result) *util.Result {
	if m.IsKeepAlive() {
		m.outbox <- res
	}
	return res
}

func (m *Mailer) doSend(msg Message) *util.Result {
	logger := m.Logger.With().Str("mod", "mailer").Str("func", "doSend").Logger()
	email := gomail.NewMSG()
	defer func() { m.lastSent = time.Now() }()

	if msg.From == "" {
		return m.outResult(util.LogMsgError(&logger, "MsgTo", "there's no From"))
	}
	if len(msg.To) < 1 {
		return m.outResult(util.LogMsgError(&logger, "MsgTo", "there's no receipts"))
	}
	if msg.Title == "" {
		return m.outResult(util.LogMsgError(&logger, "MsgTo", "there's no Title"))
	}
	if msg.Body == "" {
		return m.outResult(util.LogMsgError(&logger, "MsgTo", "there's no Body"))
	}
	email.SetFrom(msg.From).AddTo(msg.To...).SetSubject(msg.Title).SetBody(gomail.TextHTML, msg.Body)
	logger.Info().Msgf("sending [%s] to [%v]", msg.Title, msg.To)

	if msg.Attachments != nil {
		for _, fp := range msg.Attachments {
			email.Attach(&gomail.File{FilePath: fp})
		}
	}

	if len(msg.Cc) > 0 {
		email.AddCc(msg.Cc...)
	}
	if len(msg.Bcc) > 0 {
		email.AddBcc(msg.Bcc...)
	}

	if err := email.Send(m.client); err != nil {
		return m.outResult(util.LogError(&logger, "SendEmail", err))
	}
	return m.outResult(nil)
}

func (m *Mailer) Send(msg Message) *util.Result {
	if !m.running {
		return util.LogMsgError(m.Logger, "Send", "Mailer is not running, Panic!!")
	}

	if !m.IsKeepAlive() {
		var err error
		m.client, err = m.server.Connect()
		if err != nil {
			return util.Error("Connet", err)
		}
		defer m.client.Close()
		return m.doSend(msg)
	}

	m.inbox <- msg
	select {
	case res := <-m.outbox:
		return res
	case <-time.After(30 * time.Second):
		return util.LogMsgError(m.Logger, "Outbox", "timeout to get sent result")
	}
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
