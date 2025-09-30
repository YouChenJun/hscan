package core

import "github.com/YouChenJun/hscan/libs"

func ReplaceReports(module libs.Module, params map[string]string) libs.Module {
	var final []string
	for _, report := range module.Report.Final {
		final = append(final, ReplaceData(report, params))
	}
	module.Report.Final = final

	return module
}
