package libs

type Scan struct {
	InputFile  string //target list this is a file name
	Inputs     []string
	TargetInfo map[string]string
	FlowName   string
}
