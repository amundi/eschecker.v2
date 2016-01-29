package main

import (
	"github.com/amundi/escheck.v2/config"
	"github.com/amundi/escheck.v2/eslog"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_SenderInit(t *testing.T) {
	s := new(sender)
	eslog.InitSilent()

	info := &config.Query{
		TimeOut: "15s",
		Query: config.QueryInfo{
			Index:     "testindex*",
			SortBy:    "Timestamp",
			SortOrder: "ASC",
			Limit:     42,
			NbDocs:    10,
		},
	}

	err := s.initSender(info)
	assert.Nil(t, err)
	assert.Equal(t, s.index, "testindex*")
	assert.Equal(t, s.nbDocs, 10)
	assert.Equal(t, s.sortOrder, true)
	assert.Equal(t, s.sortBy, "Timestamp")
	testtimeout, _ := time.ParseDuration("15s")
	assert.Equal(t, s.timeOut, testtimeout)

	info = &config.Query{
		Query: config.QueryInfo{
			Index:     "aaaargh",
			SortBy:    "Severity",
			SortOrder: "DESC",
			Limit:     10,
			NbDocs:    53,
		},
	}

	err = s.initSender(info)
	assert.Nil(t, err)
	assert.Equal(t, s.index, "aaaargh")
	assert.Equal(t, s.nbDocs, 53)
	assert.Equal(t, s.sortOrder, false)
	assert.Equal(t, s.sortBy, "Severity")
	testtimeout, _ = time.ParseDuration("30s")
	assert.Equal(t, s.timeOut, testtimeout)

	info = &config.Query{
		Query: config.QueryInfo{},
	}
	err = s.initSender(info)
	assert.NotNil(t, err)

	info = &config.Query{}
	err = s.initSender(info)
	assert.NotNil(t, err)

	info = &config.Query{
		TimeOut: "23s",
		Query: config.QueryInfo{
			Index:  "aaaargh",
			SortBy: "Severity",
			Limit:  10,
			NbDocs: 53,
		},
	}
	err = s.initSender(info)
	assert.Nil(t, err)
	assert.Equal(t, s.index, "aaaargh")
	assert.Equal(t, s.nbDocs, 53)
	assert.Equal(t, s.sortOrder, false)
	assert.Equal(t, s.sortBy, "Severity")
	testtimeout, _ = time.ParseDuration("23s")
	assert.Equal(t, s.timeOut, testtimeout)

	info = &config.Query{
		TimeOut: "23znfi9bgf",
		Query: config.QueryInfo{
			Index:  "aaaargh",
			SortBy: "Severity",
			Limit:  10,
			NbDocs: 53,
		},
	}
	err = s.initSender(info)
	assert.NotNil(t, err)
}
