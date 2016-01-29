package queries

import (
//"github.com/amundi/escheck.v2/config"
//"github.com/amundi/escheck.v2/esmail"
//"github.com/amundi/escheck.v2/esslack"
//"github.com/amundi/escheck.v2/eslog"
//"gopkg.in/olivere/elastic.v2"
//"time"
)

/*
  This is a query example to show how to use the interface for manual queries.
  First, define the	structures you will use to do the query. The members of your
  structure can be values that you will set in the yml, for example.

type ExampleQuery struct {
	errorStatus int
  mailList []string
}

  The other structure corresponds to the mapping of the results (or "Hits") of your query.
  For example, a structure gathering the status code and timestamp fields of an
  elasticsearch document. See the golang eslatic package to learn how to get results
  from an *elastic.SearchResult structs.

type queryResult struct {
	Timestamp time.Time
	Status    int
}

	Now you shall implement the 6 methods that will satisfy the Query interface.
	The first is the setQueryConfig method that initialize the variables from your
	struct.


func (e *ExampleQuery) SetQueryConfig(c config.ManualQueryList) bool {
	e.errorStatus = c.ExampleQuery.ErrorStatus
  e.mailList = c.ExampleQuery.MailList
	return false
}

	The second method requires you to build the query. Use the Go elastic package
	documentation to do it. This query will look for statuses >= e.errorStatus in the whole
	cluster during the last hour. Beware, the query can be long and send back a lot
  of results. Try to put as much filters as you can.

func (e *ExampleQuery) BuildQuery() (elastic.Query, error) {
	query := elastic.NewBoolFilter().
		Must(elastic.NewRangeFilter("status").Gte(e.errorStatus)).
		Must(elastic.NewRangeFilter("Timestamp").Gt("now-1h"))
	return query, nil
}

	Is your critical condition reached ? This method is used to check that.
	Depending on the number of hits, or on some info into the hits, you can return
	true to trigger the action, or false if there is no emergency.

func (e *ExampleQuery) CheckCondition(search *elastic.SearchResult) bool {
	//here if there are more than 50 errors since last hour, the condition is met
	return search.Hits.TotalHits > 50
}

	To finish, if the condition is met, define what your action shall be. You can
	send a mail with the esmail package, for example.

func (e *ExampleQuery) DoAction(search *elastic.SearchResult) error {
  //say something on stdout
	eslog.Alert("There are a lot of errors in your cluster !!!")

	//or send a mail
	mail := esmail.NewMail()
  mail.SetRecipients(e.mailList)
  mail.SetSubject("houston we have a problem")
  mail.SetBody("Do something plz !!!")
  mail.Send()

  //or send a slack message
  slack := esslack.NewSlackMsg("Houston we have a problem", "Example Query", "@channelname")
  slack.Send()

	return nil
}

  In this last and optional method, you can choose to do something when an alert
  ends. You can send a new mail or slack message, or simply do nothing, like in
  the example.

func (e *ExampleQuery) OnAlertEnd() error {
  return nil
}

Don't forget to add your query in the queryList map in main.go, and your config for yml
parsing in config.go. Compile the binary... and you are done !
*/
