package core

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/YouChenJun/hscan/libs"
	"github.com/YouChenJun/hscan/utils"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
)

func ParseInput(raw string, cfg libs.Cfg) map[string]string {
	TargetInfo := ParseTarget(raw)
	dir, err := os.Getwd()
	if err == nil {
		TargetInfo["cwd"] = dir
	}
	TargetInfo["threads"] = cast.ToString(cfg.Threads)
	TargetInfo["version"] = libs.VERSION

	//time
	TargetInfo["today"] = time.Now().Format("2006-01-02")
	TargetInfo["date"] = time.Now().Format("2006-01-02T15:05:05")
	TargetInfo["timeStamp"] = utils.GetTS()

	TargetInfo["baseFolder"] = utils.NormalizePath(strings.TrimLeft(cfg.Env.BaseFolder, "/"))
	TargetInfo["binaries"] = cfg.Env.BinariesFolder
	TargetInfo["data"] = cfg.Env.DataFolder
	TargetInfo["workflow"] = cfg.Env.WorkflowFolder
	TargetInfo["Workspaces"] = cfg.Env.WorkspacesFolder

	TargetInfo["Workspace"] = utils.CleanPath(raw)
	TargetInfo["output"] = path.Join(TargetInfo["Workspaces"], TargetInfo["Workspace"])
	TargetInfo["taskId"] = cast.ToString(cfg.TaskId)

	utils.DebugF("raw:%v", raw)
	utils.DebugF("targetInfo:%v", TargetInfo)
	return TargetInfo
}

func ParseTarget(raw string) map[string]string {
	target := make(map[string]string)
	if raw == "" {
		return target
	}
	target["target"] = raw
	return target
}

// ParseFlow parse flow content
func ParseFlow(flowFile string) (libs.Flow, error) {
	utils.DebugF("Parsing workflow:%v", flowFile)
	var flow libs.Flow
	yamlcontent, err := os.ReadFile(flowFile)
	if err != nil {
		utils.ErrorF("YAML parsing %v err:%v", flowFile, err)
		return flow, err
	}
	err = yaml.Unmarshal(yamlcontent, &flow)
	if err != nil {
		utils.ErrorF("unmarshal error:%v", err)
		return flow, err
	}
	utils.DebugF("flow content:%v", flow)
	return flow, nil
}
