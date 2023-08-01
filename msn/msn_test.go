package msn

import (
	"github.com/soderasen-au/go-common/util"
	"testing"
)

var agents = []struct {
	name   string
	config Config
	want   *util.Result
}{
	{
		name: "outlook",
		config: Config{Email: &EmailServerConfig{
			ServerType: "smtp",
			Host:       util.Ptr("smtp.office365.com"),
			Port:       util.Ptr(587),
			Username:   util.Ptr("soderasen.au@outlook.com"),
			Password:   util.Ptr("tyhcpchmicjjzmes"),
			Encryption: util.Ptr(STARTTLS),
		}},
		want: nil,
	},
}

const htmlBody = `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>Test go-ipc msn agent</title>
	</head>
	<body>
		<p>This is the <b>Go IPC MSN Agent</b>.</p>
		<p>Cheers,</p>
		<p>go-ipc</p>
	</body>
</html>`

func TestAgent_Send(t *testing.T) {
	msg := Message{
		From:  "Soderasen AU <soderasen.au@outlook.com>",
		To:    []string{"jim.z.shi@gmail.com", "soderasen.au@gmail.com"},
		Cc:    []string{"jim.z.shi@outlook.com"},
		Bcc:   nil,
		Title: "Test mail from gomail",
		Body:  htmlBody,
	}

	agent, res := NewAgent(agents[0].config)
	if res != nil {
		t.Errorf("NewMailer: %v", res.Error())
	}

	res = agent.Send(msg)
	if res != nil {
		t.Errorf("Send: %v", res.Error())
	}
}
