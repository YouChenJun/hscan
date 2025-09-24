package core

import (
	"os"

	"github.com/YouChenJun/hscan/libs"
	"github.com/YouChenJun/hscan/utils"
	"github.com/fatih/color"
	"github.com/robertkrimen/otto"
)

type Scanner struct {
	Input           string
	Cfg             libs.Cfg
	target          map[string]string
	VM              *otto.Otto
	DoneFile        string
	RuntimeFile     string
	WorkspaceFolder string
	RuntimeInfo     libs.RuntimeInfo

	Workspace   string
	RoutineName string
	RoutineType string
	RoutinePath string
}

func (sc *Scanner) Scan() {
	utils.InforF("Scan Threads:%v", color.HiMagentaString("%v", sc.Cfg.Threads))
	sc.Cfg.Scan.TargetInfo = sc.target
	utils.DebugF("scan target:%v", sc.target)
	utils.MakeDir(sc.target["output"])
	sc.DoneFile = sc.target["output"] + "/done"
	sc.RuntimeFile = sc.target["output"] + "/runtime"
	sc.WorkspaceFolder = sc.target["output"]

	os.Remove(sc.DoneFile)
	utils.InforF("More info at %v", color.CyanString(sc.RuntimeFile))
	utils.InforF("Scan workflow:%v", color.CyanString(sc.Cfg.Scan.FlowName))

	sc.NewRuntime()

}

func (sc *Scanner) PreWorker() {
	sc.target = ParseInput(sc.Input, sc.Cfg)
	if sc.Cfg.Scan.FlowName != "" {
		sc.RoutineName = sc.Cfg.FlowName
	}
	sc.PrepareWorkflow()
}
