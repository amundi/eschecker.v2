package main

import (
	"github.com/amundi/escheck.v2/config"
	"github.com/amundi/escheck.v2/eslog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/olivere/elastic.v3"
	"testing"
)

func TestAutoQuery_SetQueryConfig(t *testing.T) {
	eslog.InitSilent()
	query := new(autoQuery)
	query.name = "test"
	c := config.ManualQueryList{}
	config.G_Config.Config = &config.Config{
		QueryList: map[string]config.Query{
			"test": config.Query{
				Query: config.QueryInfo{
					Index: "test*",
					Limit: 50,
				},
				Actions: config.Actions{
					List: []string{"email"},
					Email: config.Email{
						To:    []string{"tester1@test.com", "maurice@email.org"},
						Title: "Es gibt ein Problem",
						Text:  "Huge problem in your cluster",
					},
				},
			},
			"badtest": config.Query{},
		},
	}
	err := query.SetQueryConfig(c)
	assert.Equal(t, false, err)
	assert.Equal(t, 50, query.limit)
	assert.Equal(t, []string{"tester1@test.com", "maurice@email.org"}, query.mail.AlertMail.GetRecipients())
	assert.Equal(t, []string{"tester1@test.com", "maurice@email.org"}, query.mail.EndAlertMail.GetRecipients())
	assert.Equal(t, "Alert Elastic: Es gibt ein Problem", query.mail.AlertMail.GetSubject())
	assert.Equal(t, "Huge problem in your cluster", query.mail.body)
	assert.NotNil(t, query.queryInfo)
	query.name = "badtest"
	err = query.SetQueryConfig(c)
	assert.Equal(t, true, err)
	query.name = ""
	err = query.SetQueryConfig(c)
	assert.Equal(t, true, err)
}

func TestAutoQuery_BuildQuery(t *testing.T) {
	eslog.Init()
	test := new(autoQuery)
	test.queryInfo = &config.QueryInfo{
		Type: "boolquery",
		Clauses: map[string]interface{}{
			"must": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"Value", 146.5}},
				map[interface{}]interface{}{"term": []interface{}{"othervalue", "testTest"}},
				map[interface{}]interface{}{"range": []interface{}{"Timestamp", "lt", "now-1h"}},
			},
			"must_not": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"status", "OK"}},
			},
		},
	}
	realQuery := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("Value", 146.5),
		elastic.NewTermQuery("othervalue", "testTest"),
		elastic.NewRangeQuery("Timestamp").Lt("now-1h"),
	).MustNot(
		elastic.NewTermQuery("status", "OK"),
	)

	myQuery, err := test.BuildQuery()
	assert.Equal(t, nil, err)
	assert.Equal(t, realQuery, myQuery)

	test.queryInfo = &config.QueryInfo{
		Type: "boolfilter",
		Clauses: map[string]interface{}{
			"muts": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"Value", 146.5}},
				map[interface{}]interface{}{"term": []interface{}{"othervalue", "testTest"}},
			},
		},
	}
	myQuery, err = test.BuildQuery()
	assert.NotNil(t, err)
}

func TestAutoQuery_CheckCondition(t *testing.T) {
	test := new(autoQuery)
	test.limit = 10
	search := &elastic.SearchResult{
		Hits: &elastic.SearchHits{
			TotalHits: 100,
		},
	}
	assert.Equal(t, true, test.CheckCondition(search))
	test.limit = 400
	assert.Equal(t, false, test.CheckCondition(search))
}
