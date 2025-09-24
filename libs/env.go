package libs

//env file

type Env struct {
	RootFolder       string //default: ~./hscan
	BaseFolder       string //default:~/hscan-base
	DataFolder       string //this is scan result default:~/hscan-base/data/
	WorkflowFolder   string //default:~/hscan-base/workflow/
	BinariesFolder   string // ~/hscan-base/binaries use tools
	WorkspacesFolder string //default:~/hscan-base/workspace

	ConfigFile    string //config file path
	LogFile       string
	EnvConfigFile string //hscan env info path
	ExternalFile  string //bin tools config Folder
	RuleFile      string // bin tools rule folder
}
