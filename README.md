MODE D'EMPLOI - INSTRUCTIONS - GEBRAUCHSANDLEITUNG
==================================================

#### This version only works with ElasticSearch 2.X. For ES 1.X, go check the v1

## What is it ?

eschecker is a program monitoring an ElasticSearch cluster. It can send alert
by email or slack when some data reach a critical level. The program works with
a yaml file where the configuration of your cluster and other informations
like the email server etc. The file also contains a list of queries that will
be triggered by the program at regular intervals.

## Queries

A query is like an SQL request but for an ES cluster. The program allows the
creation of queries, and to do actions depending of the results. Two types of
queries are available.

**Autoqueries**

Autoqueries are queries generated from the yaml file. No need to code anything or
recompile the binary. You just need to fill some information.

The queries can be of [boolquery](https://www.elastic.co/guide/en/elasticsearch/reference/2.0/query-dsl-bool-query.html) type.
For example :

```
examplequery:                   #name of query
  schedule: 200s                #launch the query every X
  alert_onlyonce: true          #do the action only once as long as the query is in alert status
  timeout: 30s                  #timeout for every query
  alert_endmsg: true            #send a message when alert ends
  query:                        #query details
    index: myindex*             #the index of the query
    sortby: timestamp           #sort the query by a particular term
    sortorder: ASC              #sort order
    nbdocs: 10                  #max document numbers the query will send back
    limit: 5                    #the number of results that will trigger the alert
    type: boolquery            #query type
    clauses:                    #Query clauses : must, must_not and should
      must:
        - term: ["status", "Error"]                 #for example, the document must have "Error" in its status field
        - range: ["status", "gte", "300"]           #another example, the document must have a value "status" greater or equal to 300
      must_not:
        - range: ["timestamp", "gte", "now-1h"]     #the document must not be older than an hour
      filter:
        - term: ["Method", "GET"]
  actions:                                          
    list: [email, slack]                            #actions list
    email:
      to: ["myfriend@example.com"]                  #array containing the recipients
      title: Errors in my index                     #email title
      text: "myindex reports a problem !"         #the body of email. It will contain also a list of results in json format
    slack:
      channel: "@myfriend"                          #Slack channel or user that will receive message
      text: myindex reports a problem !             #message text
```


It is also possible to create a [querystring](https://www.elastic.co/guide/en/elasticsearch/reference/2.0/query-dsl-query-string-query.html#query-dsl-query-string-query).
It's a query with a simpler syntax that fits in one string :


```
examplequery:   
  schedule: 200s
  alert_onlyonce: true  
  timeout: 30s
  alert_endmsg: true
  query:     
    index: myindex*
    sortby: timestamp
    sortorder: ASC
    nbdocs: 10
    limit: 0
    type: query_string
    clauses:
      query: "(status:Error OR status:>=300) AND timestamp>now-1h"      #the query string
      analyze_wildcard: false                                           #analyze wildcards or not
  actions:                                          
    list: [email, slack]
    email:
      to: ["one@example.com", "two@example.com"]
      title: Errors in my index
      text: "myindex reports a problem !"
    slack:
      channel: "@admin"
      text: myindex reports a problem !
```

The other fields to fill are in the yaml, read the comments.


 *How to be sure that my query is right ?*

Just launch the program with the -c option. It will simply initiate the queries,
print the potential errors and quit. Beware, it doesn't mean that your query will be
accepted by ElasticSearch ! ES is very case sensitive. If you send a query with
a "timestamp" field and the right field is "Timestamp", error is guaranteed.
Beware also of the "sortby" field for the same reasons.

don't forget the brackets around the arrays in the yaml, and if you want to
specify a slack channel, put it into quotes (for example, "#mychannel"), or else
it will be parsed as a comment.

**Hand-made queries**

If the request you need cannot be made via the yaml, you have the possibility
to code your request by yourself. Knowing Golang language and having some
basing knowledge of ElasticSearch could be useful here.

First step : in the queries folder, create a file, for example myquery.go. In this
file, you shall implement all the methods of the Query interface, as defined
in queries.go. The example file examplequery.go could help you.

A handmade query allows to be more precise with the query of course, but also with
the alert condition, the action(s) to realize or what to do when the alert ends.

Some basic information is still needed in the yaml :

```
myquery:                        #request name
  schedule: 30m
  alert_onlyonce: true
  timeout: 20s
  alert_endmsg: true
  query:
    index: myindex*
    sortby: timestamp
    sortorder: ASC
    nbdocs: 10
    type: manual                #IMPORTANT : be sure to put "manual" here
```

Once the query created, it is possible to define configuration information for your
query. In config.go, you can create a structure that will be used to hold values
defined in the yaml. The values are then retrieved in the SetQueryConfig method.

Example :

In config.go

```
//create your own struct
type MyManualQuery struct {
	MyString   string
	MyBool     true
	MyInt      int
}

//and put it in the ManualQueryList struct
type ManualQueryList struct {
	MyManualQuery MyManualQuery
}
```

In config.yml (or whatever name you give it)

```
mymanualquery:
  schedule: 30m
  mystring: "salut les copains"   #user variable
  alert_onlyonce: true
  timeout: 20s
  mybool: false                   #user variable
  myint: 42                       #user variable
  alert_endmsg: true
  query:
    index: myindex*
    sortby: timestamp
    sortorder: ASC
    nbdocs: 10
    limit: 0
    type: manual
```

And finally, in the file of your query :

```
func (m *MyManualQuery) SetQueryConfig(c config.ManualQueryList) bool {
	m.myString = c.MyManualQuery.MyString
  m.myBool = c.MyManualQuery.MyBool
  m.myInt = c.MyManualQuery.MyInt
	return false
}
```

To write your query, it is necessary that you use the v3 package of elastic
for Golang.

```
import 	"gopkg.in/olivere/elastic.v3"
```

## The server

In the yaml you can choose to start a server that will display a page with the
queries' state in json format (is it up, who many times it has been triggered etc).
You can configure the path, the port, or totally deactivate it.
The page will be displayed at http://{adress-of-your-machine}:{port}/{path}.
The page can be accessed from your local network. It is possible to protect the
page with a basic HTTP authentication.

```
server_mode: true
server_port: 4242
server_path: "/escheck"
server_login: roger       #don't fill the field if you don't want a HTTP auth
server_password: rabbit   #don't fill the field if you don't want a HTTP auth
```

## rotating log

You can log the output of escheck in a rotating log. Example configuration :

```
log: true             #activate the rotating log
log_path: /var/log    #path of the log files
log_name: test        #basename of the files. They will rotate as test.1, test.2, etc.
rotate_every: 2048    #rotate every X bytes written in the file
number_of_files: 10   #number of files to keep
```

## CREDITS

escheck uses the Oliver Eilhard's Elastic package in Golang for sending
queries. It uses also the Canonical Inc. yaml package.

## TROUBLESHOOTING

A bug ? A question ? Something is not right ? Open an issue !
