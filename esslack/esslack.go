package esslack

import (
	"github.com/amundi/escheck.v2/config"
	"github.com/amundi/escheck.v2/eslog"
	"github.com/amundi/escheck.v2/worker"
	"net/http"
	"net/url"
	"os"
)

const (
	SLACKADDR = "https://slack.com/api/chat.postMessage"
)

var g_slack = struct {
	token    string
	client   *http.Client
	proxyUrl *url.URL
}{}

type SlackMsg struct {
	text    string
	user    string
	channel string
}

func Init() {
	var err error

	g_slack.token = config.G_Config.Config.Token
	g_slack.proxyUrl, err = url.Parse(getProxy())
	if err != nil {
		eslog.Error("%s : "+err.Error(), os.Args[0])
	}
	g_slack.client = &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(g_slack.proxyUrl)},
	}
}

func NewSlackMsg(text string, user string, channel string) (ret *SlackMsg) {
	ret = new(SlackMsg)
	ret.text = text
	ret.user = user
	ret.channel = channel
	return ret
}

func (s *SlackMsg) SetText(text string) {
	s.text = text
}

func (s *SlackMsg) SetUser(user string) {
	s.user = user
}

func (s *SlackMsg) SetChannel(channel string) {
	s.channel = channel
}

func (s *SlackMsg) Send() {
	collectorSlack(*s)
}

func collectorSlack(s SlackMsg) {
	worker.G_WorkQueue <- s
}

func (s SlackMsg) GetSlackRequest() string {
	if s.user == "" {
		s.user = "Elastic-Alert"
	}
	return SLACKADDR +
		"?username=" + url.QueryEscape(s.user) +
		"&token=" + getSlackToken() +
		"&channel=" + url.QueryEscape(s.channel) +
		"&pretty=1&text=" + url.QueryEscape(s.text)
}

func (s SlackMsg) DoRequest() {
	_, err := http.Get(s.GetSlackRequest())
	if err != nil {
		eslog.Error(err.Error())
	}
}

func getSlackToken() string {
	return g_slack.token
}

func getProxy() string {
	// TODO
	return ""
}
