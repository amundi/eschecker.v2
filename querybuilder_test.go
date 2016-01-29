package main

import (
	"github.com/amundi/escheck.v2/config"
	"github.com/amundi/escheck.v2/eslog"
	"github.com/stretchr/testify/assert"
	elastic "gopkg.in/olivere/elastic.v3"
	"testing"
)

func Test_getRangeQuery(t *testing.T) {
	eslog.InitSilent()
	var result elastic.Query
	var err error

	test := []interface{}{"errorcode", "gte", "500"}
	real := elastic.NewRangeQuery("errorcode").Gte(500)
	result, err = getRangeQuery(test)
	assert.Equal(t, err, nil)
	assert.Equal(t, real, result)

	test = []interface{}{"timestamp", "lt", "now-1h"}
	real = elastic.NewRangeQuery("timestamp").Lt("now-1h")
	result, err = getRangeQuery(test)
	assert.Equal(t, err, nil, "Should be nil")
	assert.Equal(t, real, result)

	test = []interface{}{"code", "lt", "800", "gte", "500"}
	real = elastic.NewRangeQuery("code").Lt(800).Gte(500)
	result, err = getRangeQuery(test)
	assert.Equal(t, err, nil, "Should be nil")
	assert.Equal(t, real, result)

	test = []interface{}{"timestamp", "lt", "now-1h", "gte", "2d"}
	real = elastic.NewRangeQuery("timestamp").Lt("now-1h").Gte("2d")
	result, err = getRangeQuery(test)
	assert.Equal(t, err, nil, "Should be nil")
	assert.Equal(t, real, result)

	test = []interface{}{"code", "lt", "800", "pouet", "500"}
	result, err = getRangeQuery(test)
	assert.NotNil(t, err)

	test = []interface{}{"code", "lt", "800", "gte"}
	result, err = getRangeQuery(test)
	assert.NotNil(t, err)

	test = []interface{}{"test", "lt", "800", "tge", "500"}
	result, err = getRangeQuery(test)
	assert.NotNil(t, err)

	test = []interface{}{"testshort"}
	result, err = getRangeQuery(test)
	assert.NotNil(t, err)

	test = []interface{}{}
	result, err = getRangeQuery(test)
	assert.NotNil(t, err)
}

func Test_getQueries(t *testing.T) {
	eslog.InitSilent()

	// term filters
	var filters []interface{} = []interface{}{
		map[interface{}]interface{}{"term": []interface{}{"test", "yes"}},
		map[interface{}]interface{}{"term": []interface{}{"required", true}},
	}

	realfilters := []elastic.Query{
		elastic.NewTermQuery("test", "yes"),
		elastic.NewTermQuery("required", true),
	}

	testfilters, err := getQueries(filters)
	assert.Equal(t, err, nil, "Should be nil")
	assert.Equal(t, realfilters, testfilters)

	realfilters = []elastic.Query{
		elastic.NewTermQuery("test", "yes"),
		elastic.NewTermQuery("required", false),
	}
	assert.NotEqual(t, realfilters, testfilters)

	//range filters
	filters = []interface{}{
		map[interface{}]interface{}{"term": []interface{}{"value", 146}},
		map[interface{}]interface{}{"term": []interface{}{"othervalue", "testTest"}},
		map[interface{}]interface{}{"range": []interface{}{"Timestamp", "gte", "now-1h"}},
	}

	realfilters = []elastic.Query{
		elastic.NewTermQuery("value", 146),
		elastic.NewTermQuery("othervalue", "testTest"),
		elastic.NewRangeQuery("Timestamp").Gte("now-1h"),
	}

	testfilters, err = getQueries(filters)
	assert.Equal(t, err, nil, "Should be nil")
	assert.Equal(t, realfilters, testfilters)

	realfilters = []elastic.Query{
		elastic.NewTermQuery("value", 146),
		elastic.NewTermQuery("othervalue", "testTest"),
		elastic.NewRangeQuery("Timestamp").Lt("now-1h"),
	}

	assert.NotEqual(t, realfilters, testfilters)

	filters = []interface{}{
		map[interface{}]interface{}{"term": []interface{}{"value", 146}},
		map[interface{}]interface{}{"term": []interface{}{"othervalue", "testTest"}},
		map[interface{}]interface{}{"range": []interface{}{"Timestamp", "lt", "now-1h"}},
	}
	testfilters, err = getQueries(filters)
	assert.Equal(t, err, nil, "Should be nil")
	assert.Equal(t, realfilters, testfilters)

	//non valid fields
	filters = []interface{}{
		map[interface{}]interface{}{"term": []interface{}{145, "yes"}},
		map[interface{}]interface{}{"term": []interface{}{"required", "yes"}},
	}
	testfilters, err = getQueries(filters)
	assert.NotNil(t, err)

	filters = []interface{}{
		map[interface{}]interface{}{"term": []interface{}{"yes"}},
		map[interface{}]interface{}{"term": []interface{}{"required", 28}},
	}
	testfilters, err = getQueries(filters)
	assert.NotNil(t, err)

	filters = []interface{}{
		map[interface{}]interface{}{"range": []interface{}{"Timestamp", "gte"}},
		map[interface{}]interface{}{"term": []interface{}{"required", 28}},
	}

	testfilters, err = getQueries(filters)
	assert.NotNil(t, err)

	filters = []interface{}{
		map[interface{}]interface{}{"range": []interface{}{42, "gte", "now-30m"}},
		map[interface{}]interface{}{"term": []interface{}{"required", 28}},
	}

	testfilters, err = getQueries(filters)
	assert.NotNil(t, err)
}

func Test_boolquery(t *testing.T) {
	eslog.InitSilent()
	info := &config.QueryInfo{}
	myQuery, err := computeQuery(info)
	assert.NotNil(t, err)

	queryInfo := &config.QueryInfo{
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
	myQuery, err = computeQuery(queryInfo)
	assert.Nil(t, err)
	assert.Equal(t, realQuery, myQuery)

	queryInfo = &config.QueryInfo{
		Type: "boolquery",
		Clauses: map[string]interface{}{
			"must_not": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"Status", "Error"}},
				map[interface{}]interface{}{"range": []interface{}{"Timestamp", "gt", "now-30m"}},
			},
		},
	}

	realQuery = elastic.NewBoolQuery().MustNot(
		elastic.NewTermQuery("Status", "Error"),
	).MustNot(
		elastic.NewRangeQuery("Timestamp").Gt("now-30m"),
	)

	myQuery, err = computeQuery(queryInfo)
	assert.Nil(t, err)
	assert.Equal(t, realQuery, myQuery)

	queryInfo = &config.QueryInfo{
		Type: "boolquery",
		Clauses: map[string]interface{}{
			"should": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"User", "Thomas"}},
			},
			"must": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"Status", "OK"}},
				map[interface{}]interface{}{"range": []interface{}{"Timestamp", "gt", "now-2h"}},
			},
		},
	}

	realQuery = elastic.NewBoolQuery().Should(
		elastic.NewTermQuery("User", "Thomas"),
	).Must(
		elastic.NewTermQuery("Status", "OK"),
		elastic.NewRangeQuery("Timestamp").Gt("now-2h"),
	)

	myQuery, err = computeQuery(queryInfo)
	assert.Nil(t, err)
	assert.Equal(t, realQuery, myQuery)

	realQuery = elastic.NewBoolQuery().Should(
		elastic.NewTermQuery("User", "Tobias"),
	).Must(
		elastic.NewTermQuery("Status", "OK"),
		elastic.NewRangeQuery("Timestamp").Gt("now-2h"),
	)
	assert.NotEqual(t, realQuery, myQuery)

	queryInfo = &config.QueryInfo{
		Type: "boolquery",
		Clauses: map[string]interface{}{
			"should": []interface{}{
				map[interface{}]interface{}{"plop": []interface{}{"Thomas"}},
			},
			"must": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"Status", "OK"}},
				map[interface{}]interface{}{"hihi": []interface{}{"Timestamp", "gt", "now-2h"}},
			},
		},
	}
	myQuery, err = computeQuery(queryInfo)
	assert.NotNil(t, err)

	// test filterQueries
	queryInfo = &config.QueryInfo{
		Type: "boolquery",
		Clauses: map[string]interface{}{
			"should": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"User", "Thomas"}},
			},
			"must": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"Status", "OK"}},
			},
			"filter": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"Method", "XPUT"}},
				map[interface{}]interface{}{"range": []interface{}{"Timestamp", "gte", "now-1h"}},
			},
		},
	}

	realQuery = elastic.NewBoolQuery().Should(
		elastic.NewTermQuery("User", "Thomas"),
	).Must(
		elastic.NewTermQuery("Status", "OK"),
	).Filter(
		elastic.NewTermQuery("Method", "XPUT"),
		elastic.NewRangeQuery("Timestamp").Gte("now-1h"),
	)

	myQuery, err = computeQuery(queryInfo)
	assert.Nil(t, err)
	assert.Equal(t, realQuery, myQuery)

	queryInfo = &config.QueryInfo{
		Type: "boolquery",
		Clauses: map[string]interface{}{
			"must_not": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"Type", "plop"}},
			},
			"filter": []interface{}{
				map[interface{}]interface{}{"term": []interface{}{"Method"}}, //missing 2nd parameter
			},
		},
	}
	myQuery, err = computeQuery(queryInfo)
	assert.NotNil(t, err)
}

func TestQueryString(t *testing.T) {

	queryInfo := &config.QueryInfo{
		Type: "query_string",
		Clauses: map[string]interface{}{
			"query": "type:MySQL AND Timestamp [2012-01-01 TO 2012-12-31]",
		},
	}
	myQuery, err := computeQuery(queryInfo)
	assert.Equal(t, err, nil)
	realQuery := elastic.NewQueryStringQuery("type:MySQL AND Timestamp [2012-01-01 TO 2012-12-31]").AnalyzeWildcard(false)
	assert.Equal(t, realQuery, myQuery)

	queryInfo = &config.QueryInfo{
		Type: "querystring",
		Clauses: map[string]interface{}{
			"query":             "this OR (that OR thi*)",
			"analyze_wildcards": true,
		},
	}
	myQuery, err = computeQuery(queryInfo)
	assert.Equal(t, err, nil)
	realQuery = elastic.NewQueryStringQuery("this OR (that OR thi*)").AnalyzeWildcard(true)
	assert.Equal(t, realQuery, myQuery)

	queryInfo = &config.QueryInfo{
		Type: "querystring",
		Clauses: map[string]interface{}{
			"query":             "this OR (that OR this)",
			"analyze_wildcards": 42,
		},
	}
	myQuery, err = computeQuery(queryInfo)
	assert.Equal(t, err, nil)
	realQuery = elastic.NewQueryStringQuery("this OR (that OR this)").AnalyzeWildcard(false)
	assert.Equal(t, realQuery, myQuery)

	//ERRORS
	queryInfo = &config.QueryInfo{
		Type: "querystring",
		Clauses: map[string]interface{}{
			"quer": "type:Error AND method:GET",
		},
	}
	myQuery, err = computeQuery(queryInfo)
	assert.NotNil(t, err)

	queryInfo = &config.QueryInfo{
		Type: "plop",
		Clauses: map[string]interface{}{
			"query": "type:Error AND method:GET",
		},
	}
	myQuery, err = computeQuery(queryInfo)
	assert.NotNil(t, err)

	queryInfo = &config.QueryInfo{
		Type: "querystring",
		Clauses: map[string]interface{}{
			"query": "",
		},
	}
	myQuery, err = computeQuery(queryInfo)
	assert.NotNil(t, err)
}

func TestStringtoNb(t *testing.T) {
	t1 := "40.7"
	t2 := 3.3
	t3 := "100"
	t4 := "nothing to do here"
	t5 := "now-30m"
	t6 := "100mille"

	r1 := stringToNb(t1)
	r2 := stringToNb(t2)
	r3 := stringToNb(t3)
	r4 := stringToNb(t4)
	r5 := stringToNb(t5)
	r6 := stringToNb(t6)
	assert.Equal(t, 40.7, r1)
	assert.Equal(t, 3.3, r2)
	assert.Equal(t, 100, r3)
	assert.Equal(t, "nothing to do here", r4)
	assert.Equal(t, "now-30m", r5)
	assert.Equal(t, "100mille", r6)

	query1 := elastic.NewBoolQuery().Should(
		elastic.NewTermQuery("User", "Thomas"),
	).Must(
		elastic.NewTermQuery("Status", stringToNb("500")),
		elastic.NewRangeQuery("Code").Gt(stringToNb("42.42")),
	)

	queryWrong := elastic.NewBoolQuery().Should(
		elastic.NewTermQuery("User", "Thomas"),
	).Must(
		elastic.NewTermQuery("Status", "500"),
		elastic.NewRangeQuery("Code").Gt("42.42"),
	)

	query2 := elastic.NewBoolQuery().Should(
		elastic.NewTermQuery("User", "Thomas"),
	).Must(
		elastic.NewTermQuery("Status", 500),
		elastic.NewRangeQuery("Code").Gt(42.42),
	)
	assert.NotEqual(t, queryWrong, query2)
	assert.Equal(t, query1, query2)
}
