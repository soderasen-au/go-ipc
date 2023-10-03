package msn

import (
	"github.com/soderasen-au/go-common/util"
	"sync"
)

type Config struct {
	Email          *EmailServerConfig `json:"email" yaml:"email"`
	DefaultMessage *Message           `json:"default_message,omitempty" yaml:"default_message"`
}

type Message struct {
	From  string
	To    []string
	Cc    []string
	Bcc   []string
	Title string
	Body  string
}

type Sender interface {
	Send(m Message) *util.Result
	Name() string
}

type Agent struct {
	mu      sync.Mutex
	Config  Config
	Senders []Sender
}

func (a Agent) Send(m Message) *util.Result {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, s := range a.Senders {
		if res := s.Send(m); res != nil {
			return res.With("Sender: " + s.Name())
		}
	}
	return nil
}

func NewAgent(cfg Config) (*Agent, *util.Result) {
	agent := &Agent{
		Config:  cfg,
		Senders: make([]Sender, 0),
	}
	if cfg.Email != nil {
		m, res := NewMailer(*cfg.Email)
		if res != nil {
			return nil, res.With("NewMailer")
		}
		agent.Senders = append(agent.Senders, m)
	}
	return agent, nil
}
