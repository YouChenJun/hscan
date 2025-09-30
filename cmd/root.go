package cmd

import (
	"fmt"
	"os"

	"github.com/YouChenJun/hscan/config"
	"github.com/YouChenJun/hscan/core"
	"github.com/YouChenJun/hscan/libs"
	"github.com/YouChenJun/hscan/utils"
	"github.com/spf13/cobra"
)

var cfg = libs.Cfg{}

var RootCmd = &cobra.Command{
	Use:   libs.BINARY,
	Short: fmt.Sprintf("%s - %s", libs.BINARY, libs.DESC),
	Long:  core.Banner(),
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfg.Env.RootFolder, "rootFolder", fmt.Sprintf("~/.%s/", libs.BINARY), "hscan root folder")
	RootCmd.PersistentFlags().StringVar(&cfg.Env.BaseFolder, "baseFolder", fmt.Sprintf("~/%s-base/", libs.BINARY), "hscan base folder")
	RootCmd.PersistentFlags().StringVar(&cfg.Env.ConfigFile, "configFile", fmt.Sprintf("~/%s/config.yaml", libs.BINARY), "configFile path")
	RootCmd.PersistentFlags().StringVar(&cfg.Env.WorkspacesFolder, "wfFolder", fmt.Sprintf("~/%s-base/workspace", libs.BINARY), "workspace folder")

	//misc
	RootCmd.PersistentFlags().BoolVarP(&cfg.Mics.Debug, "debug", "d", false, "debug mode")
	//scan
	RootCmd.PersistentFlags().StringSliceVarP(&cfg.Scan.Inputs, "target", "t", []string{}, "scan target or target list")
	RootCmd.PersistentFlags().StringVarP(&cfg.Scan.FlowName, "flow", "f", "general", "scan workflow")
	RootCmd.PersistentFlags().IntVarP(&cfg.Concurrency, "concurrency", "c", 1, "scan Concurrency")
	RootCmd.PersistentFlags().StringVar(&cfg.Tactics, "tactic", "default", "Please select your scanning strategy [default aggressive gently]")

	RootCmd.SetHelpFunc(RootHelp)
	cobra.OnInitialize(initConfig)
	RootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if cfg.Mics.FullHelp {
			cmd.UsageString()
			os.Exit(0)
		}
	}
}

func initConfig() {
	utils.InitLog(&cfg)
	if err := config.InitConfig(&cfg); err != nil {
		utils.ErrorF("failed to initialize config:%v ", err)
		utils.BadBlockF("fatal", "Please use root run hscan or install hscan")
	}
	config.ParsingConfig(&cfg)
}
