package ha

import (
	"crypto/tls"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/soderasen-au/go-common/util"
	"io"
	"net/http"
	"strings"
	"time"
)

type RunMode int

const (
	WAITING RunMode = 0
	SERVING RunMode = 1

	Healthy   string = "ok"
	UnHealthy string = "error"
)

func (m RunMode) String() string {
	switch m {
	case WAITING:
		return "WAITING"
	case SERVING:
		return "SERVING"
	}
	return "unknown"
}

type Config struct {
	PeerEndpoint   string `json:"peer_endpoint,omitempty" yaml:"peer_endpoint,omitempty"`
	SkipVerifyCert bool   `json:"skip_verify_cert,omitempty" yaml:"skip_verify_cert,omitempty"`
	PeriodInSec    int    `json:"period_in_sec,omitempty" yaml:"period_in_sec,omitempty"`
	TimeoutInMs    int    `json:"timeout_in_ms,omitempty" yaml:"timeout_in_ms,omitempty"`
	Retries        int    `json:"retries,omitempty" yaml:"retries"`
}

type Response struct {
	Status    string    `json:"status,omitempty"`
	Mode      RunMode   `json:"mode,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func (r Response) IsServing() bool {
	return strings.ToLower(r.Status) == Healthy && r.Mode == SERVING
}

type Agent struct {
	Mode     RunMode
	Config   Config
	client   *http.Client
	stopChan chan bool
	Logger   *zerolog.Logger
}

func NewAgent(cfg Config) (*Agent, *util.Result) {
	clnt := &Agent{Config: cfg}

	tlsConfig := &tls.Config{}
	if cfg.SkipVerifyCert {
		tlsConfig.InsecureSkipVerify = true
	}
	var transport *http.Transport
	transport = &http.Transport{TLSClientConfig: tlsConfig}
	clnt.client = &http.Client{
		Transport: transport,
		Timeout:   time.Duration(time.Duration(cfg.TimeoutInMs) * time.Millisecond),
	}

	return clnt, nil
}

func (a Agent) ping() (*Response, *util.Result) {
	logger := a.Logger.With().Str("HeartBeat", a.Config.PeerEndpoint).Logger()
	logger.Debug().Msg("start")
	req, err := http.NewRequest(http.MethodGet, a.Config.PeerEndpoint, nil)
	if err != nil {
		return nil, util.Error("NewHttpRequest", err)
	}

	logger.Trace().Msgf("req => %v", req)
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, util.Error("DoRequest", err)
	}
	defer resp.Body.Close()
	logger.Trace().Msgf("resp <= %v", resp)

	if resp.StatusCode != 200 {
		return nil, util.Errorf("Response(%d): %s", resp.StatusCode, resp.Status).With("DoRequest")
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, util.Error("ReadRespBody", err)
	}

	var ret Response
	err = json.Unmarshal(buf, &ret)
	if err != nil {
		return nil, util.Error("ParseResponse", err)
	}

	logger.Debug().Msgf("peer status: %v", ret)
	ret.Status = strings.ToLower(ret.Status)
	return &ret, nil
}

func (a Agent) getResponse() (*Response, *util.Result) {
	for i := 0; i < a.Config.Retries; i++ {
		resp, res := a.ping()
		if res == nil {
			return resp, nil
		}
		time.Sleep(500 * time.Millisecond)
		a.Logger.Warn().Err(res).Msgf("try[%d] failed: %s", i+1)
	}
	return nil, util.MsgError("GetResponse", "failed after retries")
}

func (a *Agent) loop() {
	logger := a.Logger.With().Str("HA", "Loop").Logger()
	logger.Info().Msg("start")
	ticker := time.NewTicker(time.Duration(a.Config.PeriodInSec) * time.Second)
	a.stopChan = make(chan bool)

	for {
		select {
		case _, ok := <-a.stopChan:
			if !ok {
				logger.Warn().Msgf("stopped in invalid status")
			}
			logger.Info().Msg("stopped")
			return
		case t := <-ticker.C:
			logger.Debug().Msgf("check heart beat at %v", t)
			resp, res := a.getResponse()
			a.handleResponse(resp, res)
		}
	}
}

func (a *Agent) handleResponse(resp *Response, res *util.Result) {
	logger := a.Logger.With().Str("HA", "handleResponse").Logger()
	logger.Debug().Msgf(" - Mode[%s] is handling %s with %s", a.Mode.String(), util.JsonStr(resp), util.JsonStr(res))
	oldMode := a.Mode
	if res != nil || resp == nil {
		a.Mode = SERVING
	} else if resp.IsServing() {
		a.Mode = WAITING
	} else {
		a.Mode = SERVING
	}
	logger.Debug().Msgf(" - New mode is %s", a.Mode.String())
	if a.Mode != oldMode {
		logger.Info().Msgf("   - Mode %s => %s", oldMode.String(), a.Mode.String())
	}
}

func (a *Agent) Start() {
	go a.loop()
}

func (a *Agent) Stop() {
	a.stopChan <- true
}
