package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/YouChenJun/hscan/libs"
	"github.com/YouChenJun/hscan/utils"
	"github.com/fatih/color"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var v *viper.Viper

func InitConfig(cfg *libs.Cfg) error {
	RootFolder := filepath.Dir(utils.NormalizePath(cfg.Env.ConfigFile))
	utils.DebugF(color.YellowString("[*]")+"hscan config file:%v", RootFolder)
	if !utils.FolderExists(RootFolder) {
		if err := os.MkdirAll(RootFolder, 0750); err != nil {
			return err
		}
	}
	//配置各种文件路径信息 basefolder 在cmd的时候已经存在默认值
	BaseFolder := utils.NormalizePath(cfg.Env.BaseFolder)
	utils.DebugF(color.YellowString("[*]")+"hscan BaseFolder:%v", BaseFolder)
	if !utils.FolderExists(BaseFolder) {
		utils.ErrorF("Base folder not found or create at path : %v", BaseFolder)
		utils.Error("Please try reinstalling the hscan !")
		os.Exit(1)
	}

	cfg.Env.EnvConfigFile = path.Join(BaseFolder, "token/hscan-var.yaml")
	utils.DebugF(color.YellowString("[*]")+"hscan EnvConfigFile:%v", cfg.Env.EnvConfigFile)
	if !utils.FolderExists(path.Dir(cfg.Env.EnvConfigFile)) {
		utils.MakeDir(path.Dir(cfg.Env.EnvConfigFile))
	}

	cfg.Env.ExternalFile = path.Join(BaseFolder, "data/external-configs/")
	utils.DebugF(color.YellowString("[*]")+"hscan ExternalFile:%v", cfg.Env.ExternalFile)
	if !utils.FolderExists(path.Dir(cfg.Env.ExternalFile)) {
		utils.MakeDir(path.Dir(cfg.Env.ExternalFile))
	}

	cfg.Env.RuleFile = path.Join(BaseFolder, "data/rules")
	utils.DebugF(color.YellowString("[*]")+"hscan RuleFile:%v", cfg.Env.RuleFile)
	if !utils.FolderExists(path.Dir(cfg.Env.RuleFile)) {
		utils.MakeDir(path.Dir(cfg.Env.RuleFile))
	}

	cfg.Update.LocalMetaData = path.Join(BaseFolder, "dist/public.json")
	utils.DebugF(color.YellowString("[*]")+"hscan LocalMetaData:%v", cfg.Update.LocalMetaData)
	cfg.Env.ConfigFile = utils.NormalizePath(cfg.Env.ConfigFile)

	v = viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path.Dir(cfg.Env.ConfigFile))
	err := v.ReadInConfig()
	if err != nil {
		//读取错误
		v.SetDefault("Database", map[string]string{
			"db_host": utils.GetOSEnv("DB_HOST", "127.0.0.1"),
			"db_port": utils.GetOSEnv("DB_PORT", "3306"),
			"db_user": utils.GetOSEnv("DB_USER", "root"),
			"db_pass": utils.GetOSEnv("DB_PASS", "root"),
			"db_name": utils.GetOSEnv("DB_NAME", "hscan"),
			"db_type": utils.GetOSEnv("DB_TYPE", "NODB"),
		})

		v.SetDefault("update", map[string]string{
			"metadata_url": "https://example.com/public.json",
			"git_repo_url": "https://github.com/xx/xx.git",
		})

		v.SetDefault("Environments", map[string]string{
			"workspaces": cfg.Env.WorkspacesFolder,
			"workflows":  path.Join(BaseFolder, "workflow"),
			"binaries":   path.Join(BaseFolder, "binaries"),
			"data":       path.Join(BaseFolder, "data"),
		})
		v.SetDefault("Tatic", map[string]any{
			"default":    runtime.NumCPU() * 4,
			"aggressive": runtime.NumCPU() * 10,
			"gently":     runtime.NumCPU() * 2,
		})
		utils.DebugF("config file path:%v", cfg.Env.ConfigFile)
		if err = v.WriteConfigAs(cfg.Env.ConfigFile); err != nil {
			utils.ErrorF("WriteConfigAs Error:%v", err)
			return err
		}
		utils.InforF("WriteConfigAs Success at %v", cfg.Env.ConfigFile)
	}
	return nil
}
func LoadConfig(cfg *libs.Cfg) *viper.Viper {
	v = viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path.Dir(cfg.Env.ConfigFile))
	if err := v.ReadInConfig(); err != nil {
		utils.ErrorF("Read config file has error :%v", err)
		os.Exit(1)
	}
	return v
}
func ParsingConfig(cfg *libs.Cfg) {
	v = LoadConfig(cfg)
	SetEnvConfig(cfg)
	SetCfg(cfg)
}

func SetCfg(cfg *libs.Cfg) {
	db := v.GetStringMapString("Database")

	cfg.DataBase.DBType = db["db_type"]
	if cfg.DataBase.DBType == "mysql" {
		cfg.DataBase.DBHost = db["db_host"]
		cfg.DataBase.DBUser = db["db_user"]
		cfg.DataBase.DBPass = db["db_pass"]
		cred := fmt.Sprintf("%v:%v", cfg.DataBase.DBUser, cfg.DataBase.DBPass)
		dest := fmt.Sprintf("%v:%v", cfg.DataBase.DBHost, cfg.DataBase.DBPort)
		cfg.DataBase.DBConnection = fmt.Sprintf("%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local", cred, dest, cfg.DataBase.DBName)
		utils.DebugF("jdbc info:%v", cfg.DataBase.DBConnection)
	}

	envs := v.GetStringMapString("Environments")
	cfg.Env.BinariesFolder = utils.NormalizePath(envs["binaries"])
	cfg.Env.WorkflowFolder = utils.NormalizePath(envs["workflows"])
	cfg.Env.DataFolder = utils.NormalizePath(envs["data"])
	utils.MakeDir(cfg.Env.DataFolder)

	update := v.GetStringMapString("update")
	cfg.Update.MetadataUlr = update["metadata_url"]
	cfg.Update.GitRepoUrl = update["git_repo_url"]
	switch cfg.Tactics {
	case "default":
		cfg.Threads = cast.ToInt(v.Get("Tatic.default"))
	case "aggressive:":
		cfg.Threads = cast.ToInt(v.Get("Tatic.aggressive"))
	case "gently:":
		cfg.Threads = cast.ToInt(v.Get("Tatic.gently"))
	default:
		cfg.Threads = cast.ToInt(v.Get("Tatic.default"))
	}
}
