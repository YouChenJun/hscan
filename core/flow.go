package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/YouChenJun/hscan/libs"
	"github.com/YouChenJun/hscan/utils"
	"github.com/thoas/go-funk"
)

func (sc *Scanner) PrepareWorkflow() {
	//list all flow or modules
	if sc.Cfg.Scan.FlowName != "" {
		//scn mode is flow
		sc.RoutineType = "flow"
	}
	sc.Cfg.Scan.FlowName = CheckFlow(sc.Cfg)
	allFlowModules := GetAllFlowModules(sc.Cfg)
	utils.DebugF("this flow all modules:%v", allFlowModules)
	var err error
	sc.RoutinePath = path.Join(sc.Cfg.Env.WorkflowFolder, fmt.Sprintf("%v.yaml", sc.Cfg.Scan.FlowName))
	sc.Cfg.Flow, err = ParseFlow(sc.RoutinePath)
	if err != nil {
		utils.ErrorF("Parse flow error :%v", err)
		os.Exit(1)
	}
	sc.target["flowPath"] = sc.RoutinePath

	for _, routine := range sc.Cfg.Flow.Routines {
		modules := SelectModules(routine.Modules, sc.Cfg)
		routine.RoutineName = fmt.Sprintf("flow-%s", sc.Cfg.Flow.Name)
		for _, module := range modules {
			parseModuleContent, err := ParseModule(module)
			if err != nil && parseModuleContent.Name == "" {
				continue
			}
			sc.TotalSteps = len(parseModuleContent.Steps)
			routine.ParsedModules = append(routine.ParsedModules, parseModuleContent)
		}
		sc.Routines = append(sc.Routines, routine)
	}
	if len(sc.Routines) == 0 {
		utils.ErrorF("flow %s not found, please check your flow name", sc.Cfg.Flow)
	}
}
func GetAllFlowModules(cfg libs.Cfg) []string {
	modePath := path.Join(cfg.Env.WorkspacesFolder, "general/*.yaml")
	if cfg.Scan.FlowName != "" {
		modePath = path.Join(cfg.Env.WorkflowFolder, fmt.Sprintf("%v/*.yaml", cfg.Scan.FlowName))
	}
	modules, err := filepath.Glob(modePath)
	if err != nil {
		utils.DebugF("Get all flow error :%v Please reinstall", err)
		return modules
	}
	return modules
}

// CheckFlow get all workflows make sure flow name is right
func CheckFlow(cfg libs.Cfg) string {
	flowPath := path.Join(cfg.Env.WorkflowFolder, "*.yaml")
	flowYaml, err := filepath.Glob(flowPath)
	if err != nil {
		utils.ErrorF("Get all workflow error :%v Please reinstall", err)
		os.Exit(1)
	}
	for _, flow := range flowYaml {
		fileName := filepath.Base(flow)
		flowNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]

		// 比较文件名(无扩展名)与传入的flowName
		if flowNameWithoutExt == cfg.Scan.FlowName {
			return cfg.Scan.FlowName
		}
	}
	utils.ErrorF("flow %s is not found,will use general flow to scan", cfg.Scan.FlowName)
	return "general"
}

func SelectModules(moduleNames []string, cfg libs.Cfg) []string {
	modePath := path.Join(cfg.Env.WorkflowFolder, "general/*.yaml")
	if cfg.Flow.Name != "" {
		modePath = path.Join(cfg.Env.WorkflowFolder, fmt.Sprintf("%v/*.yaml", cfg.Flow.Name))
	}
	modules, err := filepath.Glob(modePath)
	if err != nil {
		return modules
	}
	var selectedModules []string
	for _, item := range moduleNames {
		selectedModules = append(selectedModules, checkSelectModules(item, modules)...)
	}
	selectedModules = funk.UniqString(selectedModules)
	utils.DebugF("selected modules:%v", selectedModules)
	return selectedModules
}
func checkSelectModules(moduleName string, modules []string) []string {
	var selectedModules []string
	for _, module := range modules {
		baseModuleName := strings.Trim(strings.TrimRight(filepath.Base(module), "yaml"), ".")
		if strings.ToLower(baseModuleName) == strings.ToLower(moduleName) {
			selectedModules = append(selectedModules, module)
		}
	}
	return selectedModules
}
