package core

import "github.com/robertkrimen/otto"

func (sc *Scanner) InitVM() {
	sc.VM = otto.New()
	sc.LoadEngineScripts()
}

func (sc Scanner) LoadEngineScripts() {

}
