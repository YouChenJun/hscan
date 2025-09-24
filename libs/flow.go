package libs

type Flow struct {
	Name     string
	Desc     string
	Usage    string
	Routines []Routine
}

type Routine struct {
	RoutineName string
	FlowFolder  string `yaml:"flow"`
	Timeout     string `yaml:"timeout"`
	Modules     []string
}
