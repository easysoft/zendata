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
	deflt string
	yml string
	count  int
	fields string

	input  string
	output string
	table  = "text"
	format = "text"

	viewRes string
	viewDetail string

	help   bool

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

	flagSet = flag.NewFlagSet("zd", flag.ContinueOnError)

	flagSet.StringVar(&deflt, "d", "", "")
	flagSet.StringVar(&deflt, "default", "", "")

	flagSet.StringVar(&yml, "y", "", "")
	flagSet.StringVar(&yml, "yml", "", "")

	flagSet.StringVar(&input, "i", "", "")
	flagSet.StringVar(&input, "input", "", "")

	flagSet.IntVar(&count, "c", 10, "")
	flagSet.IntVar(&count, "count", 10, "")

	flagSet.StringVar(&fields, "field", "", "")

	flagSet.StringVar(&output, "o", "", "")
	flagSet.StringVar(&output, "output", "", "")

	flagSet.StringVar(&table, "t", "", "")
	flagSet.StringVar(&table, "table", "", "")

	flagSet.StringVar(&format, "f", "text", "")
	flagSet.StringVar(&format, "format", "text", "")

	flagSet.StringVar(&viewRes, "v", "", "")
	flagSet.StringVar(&viewDetail, "vv", "", "")

	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-help")
	} else if os.Args[1][0:1] == "-" {
		args := []string{os.Args[0], "gen"}
		args = append(args, os.Args[1:]...)
		os.Args = args
	}

	switch os.Args[1] {
	case "-s", "-set":
		set()
	case "-h", "-help":
		logUtils.PrintUsage()
	default:
		gen(os.Args)
	}
}

func set() {
	action.Set()
}

func gen(args []string) {
	if err := flagSet.Parse(args[2:]); err == nil {
		if input != "" {
			action.ParseSql(input, output)
		} else {
			action.Generate(deflt, yml, count, fields, output, format, table)
		}
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
