package action

import (
	"github.com/easysoft/zendata/internal/command"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
)

func GenData(files []string) {
	command.PrintStartInfo()

	err := command.SetOutFormat()
	defer logUtils.OutputFileWriter.Close()
	if err != nil {
		return
	}

	mainCtrl, _ := command.InitCtrl()
	mainCtrl.Generate(files)

	command.PrintEndInfo()
}
