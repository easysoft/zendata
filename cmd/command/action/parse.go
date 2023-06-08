package action

import "github.com/easysoft/zendata/internal/command"

func GenYaml(input string) {
	mainCtrl, _ := command.InitCtrl()
	mainCtrl.GenYaml(input)
}
