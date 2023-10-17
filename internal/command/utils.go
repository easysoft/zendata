package command

import (
	"fmt"
	"time"

	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

func PrintStartInfo() {
	vari.GlobalVars.StartTime = time.Now()

	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("Start at %s.", vari.GlobalVars.StartTime.Format("2006-01-02 15:04:05")))
	}
}

func PrintEndInfo() {
	vari.GlobalVars.EndTime = time.Now()

	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("End at %s.", vari.GlobalVars.EndTime.Format("2006-01-02 15:04:05")))
		dur := vari.GlobalVars.EndTime.Unix() - vari.GlobalVars.StartTime.Unix()

		logUtils.PrintTo(fmt.Sprintf("Duriation %d sec.", dur))
	}
}
