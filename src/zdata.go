package main

import (
	"flag"
	"github.com/easysoft/zendata/src/action"
	configUtils "github.com/easysoft/zendata/src/utils/config"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"syscall"
)

var (
	def    string
	count  int
	fields string

	out   string
	table = "text"
	format = "text"
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

	flagSet.StringVar(&def, "d", "", "")
	flagSet.StringVar(&def, "def", "", "")

	flagSet.IntVar(&count, "c", 10, "")
	flagSet.IntVar(&count, "count", 10, "")

	flagSet.StringVar(&fields, "field", "", "")

	flagSet.StringVar(&out, "o", "", "")
	flagSet.StringVar(&out, "out", "", "")

	flagSet.StringVar(&table, "t", "", "")
	flagSet.StringVar(&table, "table", "", "")

	flagSet.StringVar(&format, "f", "", "")
	flagSet.StringVar(&format, "format", "", "")

	flagSet.BoolVar(&vari.Verbose, "v", false, "")
	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

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
	//case "upgrade":
	//	upgrade(os.Args)
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

func upgrade(args []string) {
	if err := flagSet.Parse(args[2:]); err == nil {
		action.Upgrade()
	}
}

func gen(args []string) {
	if err := flagSet.Parse(args[2:]); err == nil {
		action.Generate(def, count, fields, out, format, table)
	}
}

func init() {
	cleanup()
	logUtils.InitLogger()
	configUtils.InitConfig()
}

func cleanup() {
	color.Unset()
}
