package libs

type Cfg struct {
	TaskId        int64  //task id use for db
	Concurrency   int    // task Concurrency
	Tactics       string //scan tactics
	Threads       int    //
	ScriptTimeout string //run timeout
	NoDB          bool
	Exclude       []string //if you wan't run this module
	Resume        bool     //Resume module
	NoPreRun      bool     //dot run pre part
	NoPostRun     bool
	NoClean       bool // dot clean data in scan
	AgentName     string
	Ip            string
	Env
	Mics
	Update
	DataBase
	Scan
	Flow
}

type Mics struct {
	FullHelp bool //show all help info
	Debug    bool // show more log info
}

type DataBase struct {
	DBType       string
	DBUser       string
	DBPass       string
	DBHost       string
	DBPort       string
	DBName       string
	DBConnection string
}
