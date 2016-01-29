package main

import (
	"flag"
	"github.com/amundi/escheck.v2/config"
	"github.com/amundi/escheck.v2/eslog"
	"github.com/amundi/escheck.v2/esmail"
	"github.com/amundi/escheck.v2/esslack"
	"github.com/amundi/escheck.v2/queries"
	"github.com/amundi/escheck.v2/worker"
	"gopkg.in/olivere/elastic.v3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	TIMELAYOUT = "Jan 2 15:04:05"
	NBWORKERS  = 64
)

//manual queries
var g_queryList map[string]queries.Query = map[string]queries.Query{
//put your manual query here
}

type Env struct {
	flagsilent *bool
	flagcheck  *bool
	filename   *string
	queries    map[string]config.Query
	client     *elastic.Client
	semaphore  chan struct{}
}

func main() {
	//get flags and config
	env := new(Env)
	env.getFlags()
	env.getConfig()

	//semaphore to queue SendRequest, or else timeout gets wrong
	env.semaphore = make(chan struct{}, 1)

	//init log, if silent, nothing will be printed on stdout
	if *env.flagsilent {
		eslog.InitSilent()
	} else {
		eslog.Init()
	}

	//init rotating log if necessary
	env.initRotatingLog()

	//init slack, mail etc.
	env.initIntegrations()

	//get queries from yaml, if flag -c activated check them and exit
	env.parseQueries()

	//init the stats for every queries and the dispatcher for workers
	initStats()
	worker.StartDispatcher(getNbWorkers())

	//connect to the elasticsearch cluster via env.client
	env.connect()

	for name, check := range g_queryList {
		go launchQuery(check, name, env)
	}

	if isServer() {
		launchServer()
	} else {
		//select to wait forever
		select {}
	}
}

func launchServer() {
	path, port := serverPath(), serverPort()
	var auth *BasicAuth

	eslog.Info("%s : launching server on path %s and port %s", os.Args[0], path, port)
	if IsServerAuthentication() {
		auth = NewBasicAuth(getServerLogin(), getServerPassword())
		auth.setDisplayFunc(collectorDisplay)
		http.HandleFunc(path, auth.BasicAuthHandler)
	} else {
		http.HandleFunc(path, collectorDisplay)
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// function that handles the life of each query
func launchQuery(c queries.Query, name string, env *Env) {
	schedule := new(scheduler)
	retries := getMaxRetries()
	send := new(sender)
	stats := queryStats{true, false, retries, 0, "None"}
	var query elastic.Query

	//query initiation
	eslog.Info("%s : getting query configuration", name)
	c.SetQueryConfig(config.G_Config.ManualConfig.List)
	query, err := c.BuildQuery()
	if err != nil {
		eslog.Error("%s : failed to build query, %s", name, err.Error())
		stats.IsUp = false
		if isServer() {
			go collectorUpdate(stats, name)
		}
		return
	}

	//scheduler and sender initiation
	schedInfo, ok := env.queries[strings.ToLower(name)]
	if ok {
		schedule.initScheduler(&schedInfo)
	} else {
		schedule.initSchedulerDefault()
	}
	err = send.initSender(&schedInfo)
	if err != nil {
		eslog.Error("%s : initSender failed, %s", name, err.Error())
		stats.IsUp = false
		if isServer() {
			go collectorUpdate(stats, name)
		}
		return
	}
	eslog.Info("%s : Starting...", name)

	//loop forever
	for {
		//try to send request. If fails, continue while decreasing attempts, or
		//die if retries reach 0.
		env.semaphore <- struct{}{}
		results, err := send.SendRequest(env.client, query)
		<-env.semaphore

		if err != nil {
			eslog.Error(err.Error())
			stats.Tries -= 1
			if stats.Tries == 0 {
				eslog.Error("%s : max attempts reached, stopping query", name)
				stats.IsUp = false
				if isServer() {
					go collectorUpdate(stats, name)
				}
				return
			} else {
				//retry after schedule
				eslog.Warning("%s : failed to connect, number of attempts left : %d", name, stats.Tries)
				if isServer() {
					go collectorUpdate(stats, name)
				}
				schedule.wait()
				continue
			}
		}

		// interpet the results, if any
		if results != nil && results.Hits != nil && results.Hits.TotalHits > 0 {
			//request succeeded, restart attempts
			stats.Tries = getMaxRetries()
			eslog.Warning("%s : found a total of %d results", name, results.Hits.TotalHits)
			yes := c.CheckCondition(results)
			if yes {
				//condition is verified. Should we enter alert status ?
				if !schedule.isAlertOnlyOnce || (schedule.isAlertOnlyOnce && !schedule.alertState) {
					eslog.Alert("%s : Action triggered", name)
					c.DoAction(results)
					schedule.alertState = true
					stats.AlertStatus = true
					stats.LastAlert = time.Now().Format(TIMELAYOUT)
					stats.NbAlerts++
				}
			} else if !yes {
				//condition not verified. Exiting alert status, triggering onAlertEnd() if necessary
				if schedule.alertState == true && schedule.isAlertEndMsg == true {
					c.OnAlertEnd()
				}
				schedule.alertState = false
				stats.AlertStatus = false
			}
		} else {
			// no results found
			eslog.Info("%s : no result found", name)
			if schedule.alertState == true && schedule.isAlertEndMsg == true {
				c.OnAlertEnd()
			}
			schedule.alertState = false
			stats.AlertStatus = false
		}
		//update the stats and display them, if necessary
		if isServer() {
			go collectorUpdate(stats, name)
		}
		//wait, and do it again
		schedule.wait()
	}
}

func (e *Env) connect() {
	var err error

	config := config.G_Config.Config
	if config == nil || config.Cluster_addr == "" {
		log.Fatal("No config info to start connection ! Check your yml")
	}

	eslog.Info("%s : connection attempt to %s", os.Args[0], config.Cluster_addr)
	if isAuthentication() {
		e.client, err = elastic.NewClient(
			elastic.SetSniff(false),
			elastic.SetURL(config.Cluster_addr),
			elastic.SetBasicAuth(getAuthLogin(), getAuthPassword()))
	} else {
		e.client, err = elastic.NewClient(elastic.SetSniff(false),
			elastic.SetURL(config.Cluster_addr))
	}
	if err != nil {
		log.Fatal(err.Error())
	} else {
		eslog.Info("%s : connection succeeded", os.Args[0])
	}
}

func (e *Env) getFlags() {
	e.flagcheck = flag.Bool("c", false, "Check the queries from the yml and exit")
	e.flagsilent = flag.Bool("q", false, "Silent output")
	e.filename = flag.String("f", "config.yml", "Specify config file name")
	flag.Parse()
}

// read YAML config and set it in config global variables
func (e *Env) getConfig() {
	source, err := ioutil.ReadFile(*e.filename)
	if err != nil {
		log.Fatal(err)
	}
	if err = e.setConfig(source); err != nil {
		log.Fatal(err)
	}
}

func (e *Env) setConfig(source []byte) error {
	var ultim config.Config
	var manual config.ManualConfig

	err := yaml.Unmarshal(source, &ultim)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(source, &manual)
	if err != nil {
		return err
	}

	config.G_Config.ManualConfig = &manual
	config.G_Config.Config = &ultim
	e.queries = ultim.QueryList
	return nil
}

func (e *Env) parseQueries() int {
	var size int

	autoQueryList := getAutoQueryList(e.queries)

	//add autoqueries to queryList
	if size = len(autoQueryList); size > 0 {
		m := initAutoQueries(autoQueryList)
		for k, v := range m {
			g_queryList[k] = v
			if *e.flagcheck {
				eslog.Info("%s : added autoquery %s", os.Args[0], k)
			}
		}
		eslog.Info("%s : %d autoqueries added", os.Args[0], size)
	}
	if *e.flagcheck {
		e.checkAndExit()
	}
	return size
}

// function to check query list, try to init them and build query, and exit with
// displaying relevant errors, if any
func (e *Env) checkAndExit() {
	var err bool
	var err2 error
	var errcount uint32

	if len(g_queryList) == 0 {
		eslog.Warning("%s : No query added, %d", os.Args[0])
		errcount++
	}

	for k, v := range g_queryList {
		eslog.Info("%s : initiating...", k)
		err = v.SetQueryConfig(config.G_Config.ManualConfig.List)
		if err {
			eslog.Error("%s : failed to get config", k)
			errcount++
			continue
		}
		_, err2 = v.BuildQuery()
		if err2 != nil {
			eslog.Error("%s : error while building query, %s", k, err2.Error())
			errcount++
		}

		//test scheduler and sender
		schedule := new(scheduler)
		send := new(sender)
		schedInfo, ok := e.queries[strings.ToLower(k)]
		if ok {
			if errinit := schedule.initScheduler(&schedInfo); errinit != nil {
				eslog.Error("%s : error in schedule info, %s", k, errinit.Error())
				errcount++
			}
		} else {
			eslog.Error("%s : error while getting query info in yaml", k)
			errcount++
		}
		if err2 = send.initSender(&schedInfo); err2 != nil {
			eslog.Error("%s : error while getting query information, %s", k, err2.Error())
			errcount++
		}
	}

	if errcount > 0 {
		eslog.Warning("%s : Detected %d errors", os.Args[0], errcount)
	} else {
		eslog.Info("%s : Everything is ok !", os.Args[0])
	}
	os.Exit(0)
}

func (e *Env) initIntegrations() {
	esmail.Init()
	esslack.Init()
}

func (e *Env) initRotatingLog() {
	if !*e.flagcheck && isRotatingLog() {
		path := config.G_Config.Config.Log_path
		if path == "" {
			path = "./"
		}
		err := eslog.InitRotatingLog(
			path+"/"+config.G_Config.Config.Log_name,
			config.G_Config.Config.Rotate_every,
			config.G_Config.Config.Number_of_files,
		)
		if err != nil {
			eslog.Error("%s : "+err.Error(), os.Args[0])
		}
	} else if !*e.flagcheck && config.G_Config.Config.Log && len(config.G_Config.Config.Log_name) == 0 {
		eslog.Error("%s : no filename specified for log", os.Args[0])
	}
}

//getters from config
func isServer() bool {
	return config.G_Config.Config.Server_mode
}

func serverPath() string {
	return config.G_Config.Config.Server_path
}

func serverPort() string {
	return config.G_Config.Config.Server_port
}

func getMaxRetries() int {
	return config.G_Config.Config.Max_retries
}

func getNbWorkers() int {
	return config.G_Config.Config.Workers
}

func isAuthentication() bool {
	return len(getAuthLogin()) > 0 && len(getAuthPassword()) > 0
}

func getAuthLogin() string {
	return config.G_Config.Config.Auth_login
}

func getAuthPassword() string {
	return config.G_Config.Config.Auth_password
}

func isRotatingLog() bool {
	return config.G_Config.Config.Log && len(config.G_Config.Config.Log_name) > 0 &&
		config.G_Config.Config.Number_of_files > 0 && config.G_Config.Config.Rotate_every > 0
}

func IsServerAuthentication() bool {
	return len(config.G_Config.Config.Server_login) > 0 && len(config.G_Config.Config.Server_password) > 0
}

func getServerLogin() string {
	return config.G_Config.Config.Server_login
}

func getServerPassword() string {
	return config.G_Config.Config.Server_password
}
