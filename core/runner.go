package core

import (
	"fmt"
	"github.com/YouChenJun/hscan/libs"
	"github.com/YouChenJun/hscan/utils"
	"github.com/fatih/color"
	"github.com/panjf2000/ants"
	"strings"
	"sync"
)

func (sc *Scanner) PrepareStepParams() {
	sc.Params = sc.target
	for _, routine := range sc.Routines {
		for _, module := range routine.ParsedModules {
			if len(module.Params) > 0 {
				for _, param := range module.Params {
					for k, v := range param {
						_, exit := sc.Params[k]
						if exit {
							continue
						}
						v = ReplaceData(v, sc.Params)
						if strings.HasPrefix(v, "~/") {
							v = utils.NormalizePath(v)
						}
						sc.Params[k] = v
					}
				}
			}
		}
	}
	sc.ReplaceRoutine()
}

func (sc *Scanner) ReplaceRoutine() {
	var routines []libs.Routine
	for _, rawRoutine := range sc.Routines {
		var routine libs.Routine
		for _, module := range rawRoutine.ParsedModules {
			module = ReplaceReports(module, sc.Params)

			sc.Reports = append(sc.Reports, module.Report.Final...)
			module.PreRun = ReplaceSlice(module.PreRun, sc.Params)

			for i, step := range module.Steps {
				module.Steps[i].StepTimeout = ReplaceData(step.StepTimeout, sc.Params)
				module.Steps[i].Threads = ReplaceData(step.Threads, sc.Params)
				module.Steps[i].Label = ReplaceData(step.Label, sc.Params)

				module.Steps[i].Required = ReplaceSlice(step.Required, sc.Params)

				module.Steps[i].Commands = ReplaceSlice(step.Commands, sc.Params)
				module.Steps[i].Scripts = ReplaceSlice(step.Scripts, sc.Params)

			}
			module.PostRun = ReplaceSlice(module.PostRun, sc.Params)
			routine.ParsedModules = append(routine.ParsedModules, module)
		}
		routines = append(routines, routine)
	}
	sc.Routines = routines
	var totalSteps, totalModules int
	parameters := make(map[string]string)
	for _, routine := range sc.Routines {
		for _, module := range routine.ParsedModules {
			for _, param := range module.Params {
				for k, v := range param {
					_, exist := parameters[k]
					if exist {
						continue
					}
					parameters[k] = v
				}
			}
			totalSteps += len(module.Steps)
			totalModules++
		}
	}
	for k, v := range sc.Params {
		parameters[k] = v
	}
	var toggleFlags, skippingFlags []string
	for key, value := range parameters {
		colorKey := color.HiMagentaString(key)

		if value == "true" {
			value = color.GreenString(value)
		} else if value == "false" {
			value = color.RedString(value)
		}

		if strings.HasPrefix(key, "enable") {
			toggleFlags = append(toggleFlags, fmt.Sprintf("%v=%v", colorKey, value))
		}

		if strings.HasPrefix(key, "skip") {
			skippingFlags = append(skippingFlags, fmt.Sprintf("%v=%v", colorKey, value))
		}

	}
	if len(toggleFlags) > 0 || len(skippingFlags) > 0 {
		utils.InforF("Toggleable and skippable parameter usage: %v, %v", strings.Join(toggleFlags, ", "), strings.Join(skippingFlags, ", "))
		utils.InforF("You can skip/enable certain parameters to speed up scanning or get more results. View more usage with %v", color.HiBlueString("hscan workflow view -v -f %v", sc.RoutineName))
	}
}

// RunRoutine 控制并发 执行扫描
func (sc *Scanner) RunRoutine(modules []libs.Module) {
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(sc.Cfg.Concurrency*10, func(m interface{}) {
		module := m.(libs.Module)
		utils.DebugF("[RunRoutine]module:%v", module)

		//sc.RunModule(module)
		wg.Done()
	}, ants.WithPreAlloc(true))
	defer p.Release()
	for _, module := range modules {
		//if funk.ContainsString(sc.Cfg.Exclude, module.Name) {
		//	utils.BadBlockF("%v 模块被排除", color.CyanString(module.Name))
		//	continue
		//}
		//utils.DebugF("[RunRoutine]运行模块:%v", module)
		p.Invoke(module)
		wg.Add(1)
	}
	wg.Wait()
}
