package main

import (
	"encoding/json"
	"fmt"
	"github.com/amundi/escheck.v2/worker"
	"net/http"
	"sync"
)

type queryStats struct {
	IsUp        bool
	AlertStatus bool
	Tries       int
	NbAlerts    int
	LastAlert   string
}

//request to update the globalstats struct
type queryStatsRequest struct {
	queryName string
	stats     queryStats
}

//request do create the stats page and display it
type displayStatsRequest struct {
	w http.ResponseWriter
	r *http.Request
	c chan struct{}
}

type globalStats struct {
	statsMap map[string]queryStats
	sync.RWMutex
}

var stats globalStats

func initStats() {
	stats.statsMap = make(map[string]queryStats)
	for k, _ := range g_queryList {
		stats.statsMap[k] = queryStats{true, false, 0, 0, "None"}
	}
}

// collector for stats update in launchQuery
func collectorUpdate(r queryStats, name string) {
	worker.G_WorkQueue <- queryStatsRequest{name, r}
}

// collector for updating page in server
// chan is used to wait and be sure that something is written in responsewriter
func collectorDisplay(w http.ResponseWriter, r *http.Request) {
	req := displayStatsRequest{w, r, make(chan struct{}, 1)}
	worker.G_WorkQueue <- req
	<-req.c
}

func (q queryStatsRequest) DoRequest() {
	stats.Lock()
	defer stats.Unlock()
	if _, exists := stats.statsMap[q.queryName]; exists {
		stats.statsMap[q.queryName] = q.stats
	}
}

func (q displayStatsRequest) DoRequest() {
	fmt.Fprintf(q.w, formatQueriesForLayout())
	q.c <- struct{}{}
}

func formatQueriesForLayout() string {
	var list []byte

	stats.RLock()
	list, err := json.MarshalIndent(stats.statsMap, "", "\t")
	stats.RUnlock()
	if err != nil {
		return "Error while getting the stats"
	} else {
		return string(list)
	}
}
