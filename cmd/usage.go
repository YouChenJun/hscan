package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func RootHelp(cmd *cobra.Command, _ []string) {
	fmt.Println(cmd.UsageString())
	if cfg.FullHelp {
		fmt.Println(cmd.UsageString())
	}
	RootUsage()
}
func RootUsage() {
	var h string
	h += ScanUsage()
	fmt.Println(h)
}
func ScanUsage() string {
	h := color.HiBlueString("CIL scan\n")
	h += "	hscan scan -T file.txt				target is file\n"
	h += "	hscan scan -t target.com			scan domain target.com"
	h += "\n"
	return h
}

func ScanHelp(cmd *cobra.Command, _ []string) {
	if cfg.FullHelp {
		fmt.Println(cmd.UsageString())
	}
	h := ScanUsage()
	fmt.Println(h)
}
