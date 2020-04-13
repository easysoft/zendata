package main

import (
	"flag"
	"github.com/easysoft/zendata/src/action"
	configUtils "github.com/easysoft/zendata/src/utils/config"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"syscall"
)

var (
	language string

	file  string
	count int
	fields string
	parse bool

	out   string
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

	flagSet.IntVar(&count, "c", 10, "")
	flagSet.IntVar(&count, "count", 10, "")

	flagSet.StringVar(&fields, "field", "", "")

	flagSet.BoolVar(&parse, "p", false, "")
	flagSet.BoolVar(&parse, "parse", false, "")

	flagSet.StringVar(&file, "o", "", "")
	flagSet.StringVar(&file, "out", "", "")

	flagSet.StringVar(&file, "t", "", "")
	flagSet.StringVar(&file, "table", "", "")

	flagSet.BoolVar(&help, "h", false, "")
	flagSet.BoolVar(&help, "help", false, "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "run", ".")
	}

	switch os.Args[1] {
	case "gen":
		gen(os.Args)
	case "set", "-set":
		action.Set()
	case "help", "-h":
		logUtils.PrintUsage()

	default: // gen
		if len(os.Args) > 1 {
			args := []string{os.Args[0], "gen"}
			args = append(args, os.Args[1:]...)

			gen(args)
		} else {
			logUtils.PrintUsage()
		}

	}
}

func gen(args []string) {
	if err := flagSet.Parse(args[2:]); err == nil {
		action.Generate(file, count, fields, out, table)
	}
}

func init() {
	cleanup()
	configUtils.InitConfig()
}

func cleanup() {
	color.Unset()
}
