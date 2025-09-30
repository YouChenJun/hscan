package libs

type Flow struct {
	Name     string
	Desc     string
	Usage    string
	Routines []Routine
}

type Routine struct {
	RoutineName   string
	FlowFolder    string `yaml:"flow"`
	Timeout       string `yaml:"timeout"`
	Modules       []string
	ParsedModules []Module
}
type Module struct {
	Name       string
	Desc       string
	ModulePath string

	Report struct {
		Final []string
	}
	Params  []map[string]string
	Steps   []Step
	PreRun  []string `yaml:"pre_run"`  // pre run commands
	PostRun []string `yaml:"post_run"` // final run commands
}

type Step struct {
	StepTimeout string `yaml:"stepTimeout"`
	Label       string
	Required    []string
	Commands    []string
	Scripts     []string
	Threads     string
}
