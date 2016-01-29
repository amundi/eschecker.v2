package main

import (
	"errors"
	"fmt"
	"github.com/amundi/escheck.v2/config"
	"gopkg.in/olivere/elastic.v3"
	"time"
)

const (
	FAIL = "Failed to initialize query information"
)

type sender struct {
	index     string
	sortBy    string
	sortOrder bool
	nbDocs    int
	timeOut   time.Duration
}

func (s *sender) initSender(info *config.Query) error {
	if info == nil {
		return errors.New("no query information")
	}

	s.index = info.Query.Index
	if s.index == "" {
		return errors.New("index cannot be empty")
	}
	s.sortBy = info.Query.SortBy
	if info.Query.SortOrder == "ASC" {
		s.sortOrder = true
	} else {
		s.sortOrder = false
	}
	s.nbDocs = info.Query.NbDocs
	if info.TimeOut == "" {
		s.timeOut = 30 * time.Second
	} else {
		timeOut, err := time.ParseDuration(info.TimeOut)
		if err != nil {
			return err
		} else {
			s.timeOut = timeOut
		}
	}
	return nil
}

//the function sends the request in a goroutine, and sends back either the results,
//either an error via their respective channels. In the meantive, the main goroutine
//waits for the results, and leave if timeout is reached.
func (s *sender) SendRequest(client *elastic.Client, query elastic.Query) (*elastic.SearchResult, error) {
	resultsChan := make(chan *elastic.SearchResult, 1)
	errChan := make(chan error, 1)

	go func(client *elastic.Client, s *sender, query elastic.Query, resultsChan chan *elastic.SearchResult, errChan chan error) {
		searchResults, err := client.Search().
			Index(s.index).              // search in index
			Query(query).                // specify the query
			Sort(s.sortBy, s.sortOrder). // sort by "timestamp" DESC. The field must exist
			From(0).Size(s.nbDocs).      // take documents 0-9
			Pretty(false).               // pretty print request and response JSON
			Do()

		if err != nil {
			errChan <- err
		} else {
			resultsChan <- searchResults
		}
	}(client, s, query, resultsChan, errChan)

	//wait for result, or leave because timeout
	select {
	case ret := <-resultsChan:
		return ret, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(s.timeOut):
		return nil, fmt.Errorf("request in index %s has timeout'ed", s.index)
	}
}
