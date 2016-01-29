package esmail

import (
	"encoding/json"
	"fmt"
	"github.com/amundi/escheck.v2/config"
	"github.com/amundi/escheck.v2/eslog"
	"github.com/amundi/escheck.v2/worker"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

const (
	TITLE          = "Alert Elastic: "
	RN             = "\r\n"
	BR             = "<br />"
	MIME_VERSION   = "1.0"
	CONTENT_TYPE   = "text/html; charset=\"utf-8\""
	EMAIL_TEMPLATE = `<html><head><title>%s</title></head><body><center><h2>%s</h2></center>%s</body></html>`
)

var g_servinfo = struct {
	server   string
	port     int
	username string
	password string
}{}

type Mail struct {
	subject string
	body    string
	to      []string
	Header  Header
}

type Header struct {
	from        string
	to          []string
	subject     string
	mimeVersion string
	contentType string
}

func Init() {
	g_servinfo.server = config.G_Config.Config.Server
	g_servinfo.port = config.G_Config.Config.Port
	g_servinfo.username = config.G_Config.Config.Username
	g_servinfo.password = config.G_Config.Config.Password
}

func NewMail() (ret *Mail) {
	ret = new(Mail)
	ret.SetFrom(g_servinfo.username)
	ret.Header.mimeVersion = MIME_VERSION
	ret.Header.contentType = CONTENT_TYPE
	return (ret)
}

func (m *Mail) SetFrom(from string) {
	m.Header.from = from
}

func (m *Mail) SetRecipients(to []string) {
	m.Header.to = to
	m.to = to
}

func (m *Mail) GetRecipients() []string {
	return m.to
}

func (m *Mail) SetSubject(subject string) {
	m.Header.subject = TITLE + subject
	m.subject = TITLE + subject
}

func (m *Mail) GetSubject() string {
	return m.subject
}

func (m *Mail) SetBody(format string, v ...interface{}) {
	m.body = fmt.Sprintf(format, v...)
}

func (m *Mail) GetBody() string {
	return m.body
}

func (m *Mail) AddToBody(format string, v ...interface{}) {
	m.body += fmt.Sprintf(format, v...)
}

func (m *Mail) ResetBody() {
	m.body = ""
}

func (m *Mail) Send() {
	collectorMail(*m)
}

func (m Mail) DoRequest() {
	var err error

	if isAuth() {
		err = smtp.SendMail(
			g_servinfo.server+":"+strconv.Itoa(g_servinfo.port),
			getAuth(),
			g_servinfo.username,
			m.to,
			[]byte(m.Header.getFullHeader()+fmt.Sprintf(EMAIL_TEMPLATE, m.subject, m.subject, m.body)+RN),
		)
	} else {
		err = smtp.SendMail(
			g_servinfo.server+":"+strconv.Itoa(g_servinfo.port),
			nil,
			g_servinfo.username,
			m.to,
			[]byte(m.Header.getFullHeader()+fmt.Sprintf(EMAIL_TEMPLATE, m.subject, m.subject, m.body)+RN),
		)
	}

	if err != nil {
		eslog.Error("%s : error sending mail : "+err.Error(), os.Args[0])
	}
}

func collectorMail(m Mail) {
	worker.G_WorkQueue <- m
}

// function to format results "pretty-print"-style
func FormatResults(res interface{}) string {
	pretty, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(pretty)
}

//the same for emails
func FormatResultsHTML(res interface{}) string {
	pretty, err := json.MarshalIndent(res, BR, "&emsp;")
	if err != nil {
		return err.Error()
	}
	return string(pretty)
}

func (h Header) getFullHeader() string {
	if h.subject == "" {
		h.subject = "Alert ElasticSearch"
	}
	return "From: " + h.from + RN +
		"To: " + strings.Join(h.to, ", ") + RN +
		"Subject: " + h.subject + RN +
		"MIME-Version: " + h.mimeVersion + RN +
		"Content-Type: " + h.contentType + RN + RN
}

func getAuth() smtp.Auth {
	return smtp.PlainAuth("",
		g_servinfo.username,
		g_servinfo.password,
		g_servinfo.server,
	)
}

func isAuth() bool {
	return len(g_servinfo.username) > 0 && len(g_servinfo.password) > 0
}
