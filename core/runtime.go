package core

import (
	"time"

	"github.com/YouChenJun/hscan/libs"
	"github.com/YouChenJun/hscan/utils"
	jsoniter "github.com/json-iterator/go"
)

func (sc *Scanner) NewRuntime() {
	sc.RuntimeInfo = libs.RuntimeInfo{
		Target:    sc.Input,
		Workspace: sc.Workspace,
	}
	sc.CreateRuntime()
}

func (sc *Scanner) CreateRuntime() {
	sc.RuntimeInfo.CreatedAt = time.Now()
	sc.RuntimeInfo.UpdatedAt = time.Now()
	if runtimeData, err := jsoniter.MarshalToString(sc.RuntimeInfo); err == nil {
		utils.WriteToFile(sc.RuntimeFile, runtimeData)
	}
}
