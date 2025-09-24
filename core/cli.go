package core

import "github.com/YouChenJun/hscan/libs"

func InitCLIScanner(input string, cfg libs.Cfg) (Scanner, error) {
	var sc Scanner
	sc.Input = input
	sc.Cfg = cfg
	sc.PreWorker()
	sc.InitVM()
	return sc, nil
}
