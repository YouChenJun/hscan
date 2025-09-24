package cmd

import (
	"os"
	"strings"
	"sync"

	"github.com/YouChenJun/hscan/core"
	"github.com/YouChenJun/hscan/libs"
	"github.com/YouChenJun/hscan/utils"
	"github.com/fatih/color"
	"github.com/panjf2000/ants"
	"github.com/spf13/cobra"
)

func init() {
	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "CLI scan",
		Long:  core.Banner(),
		RunE:  runScan,
	}
	scanCmd.SetHelpFunc(ScanHelp)
	RootCmd.AddCommand(scanCmd)
	scanCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if cfg.FullHelp {
			cmd.Help()
			os.Exit(0)
		}
	}
}

func runScan(_ *cobra.Command, _ []string) error {
	if cfg.Scan.InputFile != "" {
		//	filename
		if utils.FileExists(cfg.Scan.InputFile) {
			cfg.Scan.Inputs = append(cfg.Scan.Inputs, utils.ReadingFileUnique(cfg.Scan.InputFile)...)
		}
	}

	////if from stdin
	//stat, _ := os.Stdin.Stat()
	//if stat.Mode()&os.ModeCharDevice == 0 {
	//	sc := bufio.NewScanner(os.Stdin)
	//	for sc.Scan() {
	//		target := strings.TrimSpace(sc.Text())
	//		if err := sc.Err(); err == nil && target != "" {
	//			cfg.Scan.Inputs = append(cfg.Scan.Inputs, target)
	//		}
	//	}
	//}

	utils.InforF("Using the %v Engine %v by %v", color.GreenString(libs.BINARY), libs.AUTHOR, libs.VERSION)
	var wg sync.WaitGroup

	p, _ := ants.NewPoolWithFunc(cfg.Concurrency, func(i interface{}) {
		CreateScanner(i)
		wg.Done()
	}, ants.WithPreAlloc(true))
	defer p.Release()

	for _, target := range cfg.Scan.Inputs {
		wg.Add(1)
		_ = p.Invoke(strings.TrimSpace(target))
	}
	wg.Wait()
	return nil
}
func CreateScanner(j interface{}) {
	target := j.(string)
	sc, err := core.InitCLIScanner(target, cfg)
	if err != nil {
		utils.ErrorF("Init scanner error:%v", err)
		return
	}
	//to scan
	sc.Scan()
}
