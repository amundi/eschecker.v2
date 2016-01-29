package config

// full config struct
type Config struct {
	Cluster_addr    string
	Auth_login      string
	Auth_password   string
	Server_mode     bool
	Server_path     string
	Server_port     string
	Server_login    string
	Server_password string
	Log             bool
	Log_path        string
	Log_name        string
	Rotate_every    int
	Number_of_files int
	Workers         int
	Max_retries     int
	mailinfo
	slackinfo
	QueryList map[string]Query `yaml:"querylist"`
}

//information for each query
type Query struct {
	Schedule       string
	Alert_onlyonce bool
	TimeOut        string
	Alert_endmsg   bool
	Query          QueryInfo
	Actions        Actions
}

//information about query structs
type QueryInfo struct {
	Index     string
	SortBy    string
	SortOrder string
	NbDocs    int
	Limit     int
	Type      string
	Clauses   map[string]interface{}
}

type Actions struct {
	List  []string //list of present actions, for example ["email", "slack"]
	Email Email
	Slack Slack
}

//individual info for each query
type Email struct {
	To    []string
	Title string
	Text  string
}

//individual info for each query
type Slack struct {
	Channel string
	Text    string
	User    string
}

type mailinfo struct {
	Server   string
	Port     int
	Username string
	Password string
}

type slackinfo struct {
	Token string
}

type ManualConfig struct {
	List ManualQueryList `yaml:"querylist"`
}

/*
type ExampleQuery struct {
	ErrorStatus int
  MailList []string
}
*/

//add your manual query config here. Example :
type ManualQueryList struct {
	//ExampleQuery ExampleQuery
}

var G_Config = struct {
	ManualConfig *ManualConfig
	Config       *Config
}{}
