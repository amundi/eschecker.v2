package main

import (
	"errors"
	"github.com/amundi/escheck.v2/config"
	"gopkg.in/olivere/elastic.v3"
	"strconv"
)

/*
** The querybuilder parse Elastic queries from a map of interfaces.
 */
func getRangeQuery(v []interface{}) (*elastic.RangeQuery, error) {
	var ret *elastic.RangeQuery

	if size := len(v); size < 3 {
		return ret, errors.New("not enough values for rangeQuery")
	}

	//get name of the field
	if name, ok := v[0].(string); ok {
		ret = elastic.NewRangeQuery(name)
	} else {
		return ret, errors.New("Range Query : first parameter must be a string")
	}
	//get methods to apply
	methods := v[1:]
	if len(methods)%2 != 0 {
		return ret, errors.New("Range Query : Wrong number of parameters")
	}

	for len(methods) > 0 {
		method, ok := methods[0].(string)
		if !ok {
			return ret, errors.New("Range Query : parameter must be a string")
		}
		val := stringToNb(methods[1])
		switch method {
		case "gt":
			ret = ret.Gt(val)
		case "gte":
			ret = ret.Gte(val)
		case "lt":
			ret = ret.Lt(val)
		case "lte":
			ret = ret.Lte(val)
		default:
			return ret, errors.New("method not (yet) supported, only: gt, gte, lt, lte")
		}
		methods = methods[2:]
	}
	return ret, nil
}

func getQueries(queries []interface{}) ([]elastic.Query, error) {
	var ret []elastic.Query

	for i := 0; i < len(queries); i++ {
		term, ok := queries[i].(map[interface{}]interface{})
		if ok {
			for k, values := range term {
				v, ok := values.([]interface{})
				typ, ok2 := k.(string)
				if ok && ok2 {
					switch typ {
					case "term":
						if len(v) < 2 {
							return nil, errors.New("not enough values for termQuery")
						}
						if name, ok := v[0].(string); ok {
							ret = append(ret, elastic.NewTermQuery(name, v[1]))
						} else {
							return nil, errors.New("termQuery: first value of array must be a string")
						}
					case "range":
						rangeQuery, err := getRangeQuery(v)
						if err != nil {
							return nil, err
						}
						ret = append(ret, rangeQuery)
					default:
						return nil, errors.New("Query not (yet) supported")
					}
				} else {
					return nil, errors.New("wrong types for query")
				}
			}
		} else {
			return nil, errors.New("Query badly formatted")
		}
	}
	return ret, nil
}

func boolQuery(clauses map[string]interface{}) (elastic.Query, error) {
	var mustQueries, mustNotQueries, shouldQueries, filterQueries []elastic.Query
	var err error

	//get must, must not, should clauses, if presents
	if mustClauses, ok := clauses["must"].([]interface{}); ok {
		mustQueries, err = getQueries(mustClauses)
		if err != nil {
			return nil, err
		}
	}

	if mustNotClauses, ok := clauses["must_not"].([]interface{}); ok {
		mustNotQueries, err = getQueries(mustNotClauses)
		if err != nil {
			return nil, err
		}
	}

	if shouldClauses, ok := clauses["should"].([]interface{}); ok {
		shouldQueries, err = getQueries(shouldClauses)
		if err != nil {
			return nil, err
		}
	}

	if filterClauses, ok := clauses["filter"].([]interface{}); ok {
		filterQueries, err = getQueries(filterClauses)
		if err != nil {
			return nil, err
		}
	}

	if len(mustQueries) == 0 && len(mustNotQueries) == 0 && len(shouldQueries) == 0 &&
		len(filterQueries) == 0 {
		// empty query, might be an Error
		return nil, errors.New("No query clauses specified, query building failed")
	}

	return elastic.NewBoolQuery().Must(mustQueries...).
		MustNot(mustNotQueries...).
		Should(shouldQueries...).Filter(filterQueries...), nil
}

func queryString(clauses map[string]interface{}) (elastic.Query, error) {
	var ok bool
	analyzeWildcard := false
	var query string

	// do not allow emty strings
	if query, ok = clauses["query"].(string); !ok || query == "" {
		return nil, errors.New("missing query parameter in query string")
	}
	if wildcard, ok := clauses["analyze_wildcards"].(bool); ok {
		analyzeWildcard = wildcard
	}
	return elastic.NewQueryStringQuery(query).AnalyzeWildcard(analyzeWildcard), nil
}

func computeQuery(queryInfo *config.QueryInfo) (elastic.Query, error) {
	if queryInfo == nil || queryInfo.Type == "manual" {
		return nil, errors.New("no query info or query is not an autoquery")
	}
	switch queryInfo.Type {
	case "boolquery":
		return boolQuery(queryInfo.Clauses)
	case "query_string", "querystring":
		return queryString(queryInfo.Clauses)
	}
	return nil, errors.New("type of query unknown")
}

func stringToNb(value interface{}) interface{} {
	switch t := value.(type) {
	case string:
		if ret, err := strconv.Atoi(t); err == nil {
			return ret
		} else if ret, err := strconv.ParseFloat(t, 64); err == nil {
			return ret
		} else {
			return value
		}
	default:
		return value
	}
}
