package main

import (
	"encoding/hex"
	"github.com/amundi/escheck.v2/config"
	"github.com/amundi/escheck.v2/eslog"
	"github.com/amundi/escheck.v2/queries"
	"github.com/stretchr/testify/assert"
	"gopkg.in/olivere/elastic.v3"
	"strings"
	"testing"
)

func TestDummyYaml1(t *testing.T) {
	e := new(Env)
	check := false
	silent := true
	e.flagcheck = &check
	e.flagsilent = &silent
	eslog.InitSilent()
	dummy, err := hex.DecodeString("2364756d6d792079616d6c2066696c6520666f72" +
		"2074657374732c2069676e6f7265206d65" +
		"20666f72657665720d0a232071756572796c6973743a0d0a232020206e616d656f6671756" +
		"572793a0d0a232020202020746167676c653a20747275650d0a2320202020207363686564" +
		"756c653a203330730d0a23202020202071756572793a0d0a23202020202020696e6465783" +
		"a207061636b6574626561742a0d0a23202020202020736f727462793a2074696d65737461" +
		"6d700d0a23202020202020736f72746f726465723a204153430d0a232020202020206e626" +
		"46f63733a2031300d0a232020202020206c696d69743a2031300d0a232020202020207479" +
		"70653a20626f6f6c66696c7465720d0a23202020202020636c61757365733a0d0a2320202" +
		"020202020206d7573743a0d0a232020202020202020202d207465726d3a205b226d657468" +
		"6f64222c202265786368616e67652e6465636c617265225d0d0a232020202020202020202" +
		"d207465726d3a205b22737461747573222c20224f4b225d0d0a232020202020202020202d" +
		"2072616e67653a205b2274696d657374616d70222c2022677465222c20226e6f772d31682" +
		"25d0d0a230d0a230d0a230d0a230d0a23207468652061646472657373206f662074686520" +
		"455320636c757374657220746f206d6f6e69746f720d0a636c75737465725f616464723a2" +
		"0687474703a2f2f7777772e6d79636c75737465722e636f6d3a393230302f0d0a0d0a2320" +
		"646f20796f752077616e7420616e20696e666f20776562706167652c20776865726520616" +
		"e64206f6e20776869636820706f72742e20506f7274206d757374206265203e2031303234" +
		"0d0a7365727665725f6d6f64653a20747275650d0a7365727665725f706f72743a2034323" +
		"4320d0a7365727665725f706174683a20222f776f6f686f6f220d0a0d0a2320646f20796f" +
		"752077616e74206120726f746174696e67206c6f6720616e642077686572650d0a6c6f673" +
		"a20747275650d0a6c6f675f706174683a202f7661722f6c6f670d0a6c6f675f6e616d653a" +
		"206d6f6e69746f720d0a726f746174655f65766572793a20313032340d0a6e756d6265725" +
		"f6f665f66696c65733a20370d0a0d0a23206e756d626572206f6620726571756573742061" +
		"7474656d707473206265666f726520636c6f73696e67206120676f726f7574696e652e205" +
		"07574202d3120696620796f752077616e740d0a2320746f2074727920746f20646f207265" +
		"71756573747320666f726576657220286e6f74206120676f6f642069646561290d0a6d617" +
		"85f726574726965733a20330d0a0d0a23206e756d626572206f6620776f726b6572732069" +
		"6e20746865207461736b2071756575652e205468697320616666656374732074686520737" +
		"0656564206174207768696368207461736b73206c696b650d0a232073656e64696e672065" +
		"6d61696c732f736c61636b206d65737361676573206172652070726f6365737365642e204" +
		"d6f6469667920746869732076616c756520696620796f7520686176650d0a232061206c6f" +
		"74206f66207175657269657320616e64207468652070726f6772616d207374727567676c6" +
		"520746f2068616e646c6520746865206368617267652e0d0a776f726b6572733a20313238" +
		"0d0a0d0a2320656d61696c2073657276657220696e666f726d6174696f6e2e20466f72207" +
		"3656e64696e6720656d61696c732e0d0a6d61696c696e666f3a0d0a20207365727665723a" +
		"206d797365727665720d0a2020706f72743a20343234320d0a2020757365726e616d653a2" +
		"0726f624068656c6c6f2e636f6d0d0a202070617373776f72643a207365637265747a0d0a" +
		"0d0a736c61636b696e666f3a0d0a2020746f6b656e3a20226b696b6f6f6c65746f6b656e3" +
		"432220d0a0d0a71756572796c6973743a0d0a202074657374313a0d0a20202020616c6572" +
		"745f6f6e6c796f6e63653a20747275650d0a202020207363686564756c653a2033336d0d0" +
		"a20202020616c6572745f656e646d73673a20747275650d0a20202020616374696f6e733a" +
		"0d0a2020202020206c6973743a205b656d61696c2c20736c61636b5d0d0a2020202020206" +
		"56d61696c3a0d0a2020202020202020746f3a205b226a65616e2d6d696368406578616d70" +
		"6c652e636f6d222c2022676572617264406578616d706c652e636f6d225d0d0a202020202" +
		"02020207469746c653a2054657374207469747265202121210d0a20202020202020207465" +
		"78743a2059206120756e2070726f626c656d65206d65630d0a202020202020736c61636b3" +
		"a0d0a2020202020202020746578743a2059206120756e2070726f626c656d65206d65630d" +
		"0a20202020202020206368616e6e656c3a20272367656e6572616c270d0a2020202020202" +
		"020757365723a204368696368617269746f0d0a2020202071756572793a0d0a2020202020" +
		"20696e6465783a20746573742a0d0a202020202020736f72746f726465723a204153430d0" +
		"a202020202020736f727462793a20636f64650d0a2020202020206e62646f63733a20330d" +
		"0a2020202020206c696d69743a2034300d0a202020202020747970653a20626f6f6c71756" +
		"572790d0a202020202020636c61757365733a0d0a20202020202020206d7573745f6e6f74" +
		"3a0d0a202020202020202020202d207465726d3a205b22737461747573222c20226f6b225" +
		"d0d0a202020202020202020202d2072616e67653a205b22636f6465222c20226c74222c20" +
		"3530305d0d0a202020202020202066696c7465723a0d0a202020202020202020202d20746" +
		"5726d3a205b226d6574686f64222c2022474554225d0d0a202074657374323a0d0a202020" +
		"20616c6572745f6f6e6c796f6e63653a2066616c73650d0a2020202072656d696e6465723" +
		"a20333030306d0d0a202020207363686564756c653a203330680d0a20202020616c657274" +
		"5f656e646d73673a2066616c73650d0a20202020616374696f6e733a0d0a2020202020206" +
		"c6973743a205b656d61696c5d0d0a202020202020656d61696c3a0d0a2020202020202020" +
		"746f3a205b6d6f6940686f746d61696c2e636f6d5d0d0a2020202071756572793a0d0a202" +
		"020202020696e6465783a2074657374322a0d0a202020202020736f72746f726465723a20" +
		"444553430d0a202020202020736f727462793a2074696d657374616d700d0a20202020202" +
		"06e62646f63733a203138300d0a2020202020206c696d69743a20300d0a20202020202074" +
		"7970653a2071756572795f737472696e670d0a202020202020636c61757365733a0d0a202" +
		"020202020202071756572793a20227374617475733a4572726f72204f5220737461747573" +
		"3a3e3d333030220d0a2020202020202020616e616c797a655f77696c64636172643a20666" +
		"16c73650d0a")
	assert.Nil(t, err)

	err = e.setConfig(dummy)
	assert.Nil(t, err)
	e.initIntegrations()
	c := config.G_Config.Config
	nb := e.parseQueries()
	assert.Equal(t, 2, nb)

	//cluster and server info
	assert.Equal(t, "http://www.mycluster.com:9200/", c.Cluster_addr)
	assert.Equal(t, true, c.Server_mode)
	assert.Equal(t, "4242", c.Server_port)
	assert.Equal(t, "/woohoo", c.Server_path)
	assert.Equal(t, 3, c.Max_retries)
	assert.Equal(t, 3, getMaxRetries())
	assert.Equal(t, 128, c.Workers)
	assert.Equal(t, 128, getNbWorkers())
	assert.Equal(t, true, c.Log)
	assert.Equal(t, 1024, c.Rotate_every)
	assert.Equal(t, 7, c.Number_of_files)
	assert.Equal(t, "/var/log", c.Log_path)
	assert.Equal(t, "monitor", c.Log_name)

	//test1
	test1, ok := g_queryList["test1"]
	assert.Equal(t, true, ok)
	_, ok = test1.(queries.Query)
	assert.Equal(t, true, ok, "test1 should implement queries.Query")
	errconfig := test1.SetQueryConfig(config.G_Config.ManualConfig.List)
	assert.Equal(t, false, errconfig)
	autoTest1, ok := test1.(*autoQuery)
	assert.Equal(t, true, ok, "Test1 should be an Autoquery")
	assert.Equal(t, 40, autoTest1.limit)
	realQuery := elastic.NewBoolQuery().
		MustNot(elastic.NewTermQuery("status", "ok")).
		MustNot(elastic.NewRangeQuery("code").Lt(500)).
		Filter(elastic.NewTermQuery("method", "GET"))
	myQuery, err := test1.BuildQuery()
	assert.Nil(t, err)
	assert.Equal(t, realQuery, myQuery)
	assert.Equal(t, 40, autoTest1.limit)
	assert.Equal(t, []string{"jean-mich@example.com", "gerard@example.com"}, autoTest1.mail.EndAlertMail.GetRecipients())

	//slack
	msg := "https://slack.com/api/chat.postMessage?username=Chicharito&token=" + "kikooletoken42" + "&channel=%23general&pretty=1&text=Y+a+un+probleme+mec"
	assert.Equal(t, msg, autoTest1.slack.msg.GetSlackRequest())
	msgEnd := "https://slack.com/api/chat.postMessage?username=Chicharito&token=" + "kikooletoken42" + "&channel=%23general&pretty=1&text=End+of+alert+for+test1"
	assert.Equal(t, msgEnd, autoTest1.slack.endMsg.GetSlackRequest())
	assert.Equal(t, "Alert Elastic: Test titre !!!", autoTest1.mail.AlertMail.GetSubject())
	assert.Equal(t, "Y a un probleme mec", autoTest1.mail.body)

	s := new(scheduler)
	schedInfo, ok := e.queries[strings.ToLower(autoTest1.name)]
	assert.Equal(t, true, ok)
	s.initScheduler(&schedInfo)
	assert.Equal(t, true, s.isAlertOnlyOnce)
	assert.Equal(t, "33m0s", s.waitSchedule.String())

	sd := new(sender)
	notGood := sd.initSender(&schedInfo)
	assert.Nil(t, notGood)
	assert.Equal(t, "test*", sd.index)
	assert.Equal(t, 3, sd.nbDocs)
	assert.Equal(t, true, sd.sortOrder)
	assert.Equal(t, "code", sd.sortBy)

	//test2
	test2, ok := g_queryList["test2"]
	assert.Equal(t, true, ok)
	_, ok = test2.(queries.Query)
	assert.Equal(t, true, ok, "test2 should implement queries.Query")
	errconfig = test2.SetQueryConfig(config.G_Config.ManualConfig.List)
	assert.Equal(t, false, errconfig)
	autoTest2, ok := test2.(*autoQuery)
	assert.Equal(t, true, ok, "Test2 should be an Autoquery")
	assert.Equal(t, 0, autoTest2.limit)
	realQuery2 := elastic.NewQueryStringQuery("status:Error OR status:>=300").AnalyzeWildcard(false)
	myQuery, err = test2.BuildQuery()
	assert.Nil(t, err)
	assert.Equal(t, realQuery2, myQuery)
	assert.Equal(t, []string{"moi@hotmail.com"}, autoTest2.mail.EndAlertMail.GetRecipients())
	assert.Equal(t, "Alert Elastic: ", autoTest2.mail.AlertMail.GetSubject())
	assert.Equal(t, "", autoTest2.mail.body)

	s = new(scheduler)
	schedInfo, ok = e.queries[strings.ToLower(autoTest2.name)]
	assert.Equal(t, true, ok, "Should be ok")
	s.initScheduler(&schedInfo)
	assert.Equal(t, false, s.isAlertOnlyOnce)
	assert.Equal(t, "30h0m0s", s.waitSchedule.String())

	sd = new(sender)
	notGood2 := sd.initSender(&schedInfo)
	assert.Nil(t, notGood2)
	assert.Equal(t, "test2*", sd.index)
	assert.Equal(t, 180, sd.nbDocs)
	assert.Equal(t, false, sd.sortOrder)
	assert.Equal(t, "timestamp", sd.sortBy)
}
