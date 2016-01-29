package main

import (
	"github.com/amundi/escheck.v2/eslog"
	"github.com/amundi/escheck.v2/worker"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initStatsForTests1() {
	stats.statsMap = make(map[string]queryStats)
	stats.statsMap["Test"] = queryStats{true, false, 3, 0, "Yesterday"}
	stats.statsMap["Test"] = queryStats{true, false, 3, 0, "Yesterday"}
}

func initStatsForTests2() {
	stats.statsMap = make(map[string]queryStats)
	stats.statsMap["Test"] = queryStats{true, true, 3, 0, "Now"}
}

func Test_DisplayPage(t *testing.T) {
	eslog.InitSilent()
	ts := httptest.NewServer(http.HandlerFunc(collectorDisplay))
	defer ts.Close()
	initStatsForTests1()
	worker.StartDispatcher(32)

	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)
	expected := formatQueriesForLayout()
	assert.Equal(t, expected, string(page))
	worker.StopAllWorkers(32)
}
