package config

import (
	"fmt"
	"path"

	"github.com/YouChenJun/hscan/libs"
	"github.com/YouChenJun/hscan/utils"
	"github.com/spf13/viper"
)

func SetEnvConfig(cfg *libs.Cfg) {
	e := viper.New()
	e.SetConfigName("hscan-var")
	e.SetConfigType("yaml")
	e.AddConfigPath(path.Dir(cfg.Env.EnvConfigFile))
	err := e.ReadInConfig()
	agentName := fmt.Sprintf("xscan_%s", utils.GenHash(utils.GetTS())[:8])

	if err != nil {
		e.SetDefault("Agent", map[string]string{
			"agent_name": agentName,
			"ip":         utils.GetPublicIP(),
		})
		if err = e.WriteConfigAs(cfg.Env.EnvConfigFile); err != nil {
			utils.ErrorF("Write env config error: %s", err)
			return
		}
	}
	configIp := e.GetString("Agent.ip")
	if utils.GetPublicIP() != configIp {
		e.Set("Agent.ip", utils.GetPublicIP())
		e.Set("Agent.agent_name", agentName)
		if err = e.WriteConfig(); err != nil {
			utils.ErrorF("Write env config error: %s", err)
			return
		}
		utils.InforF("hscan agent name :%v,current ip :%v", agentName, utils.GetPublicIP())
	}
	envs := e.GetStringMapString("Agent")
	cfg.AgentName = envs["agent_name"]
	cfg.Ip = envs["ip"]
}
