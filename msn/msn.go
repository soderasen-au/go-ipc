package msn

type Config struct {
	Email *EmailServerConfig `json:"email" yaml:"email"`
}

type Message struct {
	From  string
	To    []string
	Cc    []string
	Bcc   []string
	Title string
	Body  string
}
