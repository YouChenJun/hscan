package core

import (
	"github.com/YouChenJun/hscan/libs"
	"github.com/fatih/color"
)

func Banner() string {
	version := color.HiWhiteString(libs.VERSION)
	author := color.MagentaString(libs.AUTHOR)
	b := color.GreenString(`

██   ██ ███████  ██████  █████  ███    ██ 
 ██ ██  ██      ██      ██   ██ ████   ██ 
  ███   ███████ ██      ███████ ██ ██  ██ 
 ██ ██       ██ ██      ██   ██ ██  ██ ██ 
██   ██ ███████  ██████ ██   ██ ██   ████ 

	`)
	b += "\n\n\t" + color.GreenString(`                  Xscan Next Generation %v`, version) + color.GreenString(` by %v`, author)
	b += "\n\n" + color.HiCyanString(`	                    %s`, libs.DESC) + "\n"
	color.Unset()
	return b
}
