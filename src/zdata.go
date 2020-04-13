package main

import (
	"flag"
	configUtils "github.com/easysoft/zendata/src/utils/config"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"syscall"
)

var (
	language string

	file  string
	count string
	field string
	parse string

	out   bool
	table string
	help  bool

	flagSet *flag.FlagSet
)

func main() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		cleanup()
		os.Exit(0)
	}()

	flagSet = flag.NewFlagSet("zdata", flag.ContinueOnError)

	flagSet.StringVar(&file, "f", "", "")
	flagSet.StringVar(&file, "file", "", "")

	flagSet.StringVar(&file, "c", "", "")
	flagSet.StringVar(&file, "count", "", "")

	flagSet.StringVar(&file, "field", "", "")

	flagSet.StringVar(&file, "p", "", "")
	flagSet.StringVar(&file, "parse", "", "")

	flagSet.StringVar(&file, "o", "", "")
	flagSet.StringVar(&file, "out", "", "")

	flagSet.StringVar(&file, "t", "", "")
	flagSet.StringVar(&file, "table", "", "")

	flagSet.StringVar(&file, "h", "", "")
	flagSet.StringVar(&file, "help", "", "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "run", ".")
	}

	switch os.Args[1] {
	case "checkout", "co":
		if err := flagSet.Parse(os.Args[2:]); err == nil {

		}
	case "help", "-h":
		logUtils.PrintUsage()

	default: // run
		if len(os.Args) > 1 {
			args := []string{os.Args[0], "run"}
			args = append(args, os.Args[1:]...)

			//run(args)
		} else {
			logUtils.PrintUsage()
		}

	}
}

func init() {
	cleanup()

	vari.RunFromCui = false
	configUtils.InitConfig()
}

func cleanup() {
	color.Unset()
}
