package main

import (
	"encoding/json"
	"fmt"
	"github.com/amundi/escheck.v2/config"
	"github.com/amundi/escheck.v2/eslog"
	"github.com/amundi/escheck.v2/esmail"
	"github.com/amundi/escheck.v2/esslack"
	"gopkg.in/olivere/elastic.v3"
)

type mailer struct {
	body         string
	AlertMail    *esmail.Mail
	EndAlertMail *esmail.Mail
}

type slacker struct {
	msg    *esslack.SlackMsg
	endMsg *esslack.SlackMsg
}

type autoQuery struct {
	name      string //the name of the query, to get from the yml
	limit     int    //the limit for checkcondition
	queryInfo *config.QueryInfo
	//integrations
	actionList []string
	mail       *mailer //pointer rather than a struct in case of action doesn't exist
	slack      *slacker
}

func (a *autoQuery) SetQueryConfig(c config.ManualQueryList) bool {
	//autoqueries don't need the manualquery list
	info, ok := config.G_Config.Config.QueryList[a.name]
	if !ok {
		eslog.Error("%s : failed to get query configuration", a.name)
		return true
	}
	a.actionList = info.Actions.List
	if len(a.actionList) == 0 {
		eslog.Error("%s : No action defined", a.name)
		return true
	}
	for _, val := range a.actionList {

		switch val {
		case "email":
			if len(info.Actions.Email.To) == 0 {
				eslog.Error("%s : No recipients defined for email action", a.name)
				return true
			}
			a.initMailerForAutoQuery(&info)
		case "slack":
			if len(info.Actions.Slack.Channel) == 0 {
				eslog.Error("%s : No channels defined for slack action", a.name)
				return true
			}
			a.initSlackForAutoQuery(info.Actions.Slack)
		}
	}
	a.limit = info.Query.Limit
	a.queryInfo = &info.Query
	return false
}

func (a *autoQuery) BuildQuery() (elastic.Query, error) {
	return computeQuery(a.queryInfo)
}

func (a *autoQuery) CheckCondition(search *elastic.SearchResult) bool {
	return search.Hits.TotalHits >= int64(a.limit)
}

func (a *autoQuery) DoAction(search *elastic.SearchResult) error {

	for i := 0; i < len(a.actionList); i++ {
		switch a.actionList[i] {
		case "email":
			//pretty-format the results if they exists, set them in the body with base text
			if size := len(search.Hits.Hits); size > 0 {
				res := make([]*json.RawMessage, size)
				for i, hit := range search.Hits.Hits {
					res[i] = hit.Source
				}
				pretty := esmail.FormatResultsHTML(res)

				a.mail.AlertMail.SetBody("<p>%s</p><p>Here is an excerpt of results : %s</p>", a.mail.body, pretty)
			} else {
				//no results to add, just send the m.text
				a.mail.AlertMail.SetBody("<p>%s</p>", a.mail.body)
			}
			a.mail.AlertMail.Send()
			a.mail.AlertMail.ResetBody()
		case "slack":
			a.slack.msg.Send()
		}
	}
	return nil
}

func (a *autoQuery) OnAlertEnd() error {
	for i := 0; i < len(a.actionList); i++ {
		switch a.actionList[i] {
		case "email":
			a.mail.EndAlertMail.Send()
		case "slack":
			a.slack.endMsg.Send()
		}
	}
	return nil
}

// getAutoQueryList get the list of autoqueries (not manual queries)
func getAutoQueryList(list map[string]config.Query) (ret []string) {
	ret = []string{}

	for k, v := range list {
		if v.Query.Type != "manual" {
			ret = append(ret, k)
		}
	}
	return ret
}

//this func add the autoqueries and their names to the queryList map
func initAutoQueries(queries []string) map[string]*autoQuery {
	ret := map[string]*autoQuery{}
	for _, v := range queries {
		ret[v] = new(autoQuery)
		ret[v].name = v
	}
	return ret
}

func (a *autoQuery) initMailerForAutoQuery(info *config.Query) {
	a.mail = new(mailer)
	a.mail.AlertMail = esmail.NewMail()
	a.mail.EndAlertMail = esmail.NewMail()
	a.mail.body = info.Actions.Email.Text
	a.mail.AlertMail.SetSubject(info.Actions.Email.Title)
	a.mail.AlertMail.SetRecipients(info.Actions.Email.To)
	a.mail.AlertMail.SetFrom(a.name)
	a.mail.EndAlertMail.SetSubject("End of alert")
	a.mail.EndAlertMail.SetRecipients(info.Actions.Email.To)
	a.mail.EndAlertMail.SetFrom(a.name)
	a.mail.EndAlertMail.SetBody("End of alert for query %s", a.name)
}

func (a *autoQuery) initSlackForAutoQuery(info config.Slack) {
	a.slack = new(slacker)
	a.slack.msg = esslack.NewSlackMsg(info.Text, info.User, info.Channel)
	a.slack.endMsg = esslack.NewSlackMsg(fmt.Sprintf("End of alert for %s", a.name), info.User, info.Channel)
}
