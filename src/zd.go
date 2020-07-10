package main

import (
	"flag"
	"github.com/easysoft/zendata/src/action"
	configUtils "github.com/easysoft/zendata/src/utils/config"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

var (
	deflt  string
	yaml   string
	count  int
	fields string

	input  string
	output string
	table  string
	format = constant.FormatText

	viewRes string
	viewDetail string

	example bool
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

	flagSet.StringVar(&yaml, "c", "", "")
	flagSet.StringVar(&yaml, "config", "", "")

	flagSet.StringVar(&input, "i", "", "")
	flagSet.StringVar(&input, "input", "", "")

	flagSet.IntVar(&count, "n", 10, "")
	flagSet.IntVar(&count, "lines", 10, "")

	flagSet.StringVar(&fields, "F", "", "")
	flagSet.StringVar(&fields, "field", "", "")

	flagSet.StringVar(&output, "o", "", "")
	flagSet.StringVar(&output, "output", "", "")

	flagSet.StringVar(&table, "t", "table_name", "")
	flagSet.StringVar(&table, "table", "table_name", "")

	flagSet.StringVar(&viewRes, "v", "", "")
	flagSet.StringVar(&viewDetail, "vv", "", "")

	flagSet.StringVar(&vari.HeadSep, "H", "\t", "")
	flagSet.StringVar(&vari.HeadSep, "human", "\t", "")

	flagSet.IntVar(&vari.Length, "length", 0, "")
	flagSet.StringVar(&vari.LeftPad, "leftPad", "", "")
	flagSet.StringVar(&vari.RightPad, "rightPad", "", "")

	flagSet.BoolVar(&vari.HttpService, "s", false, "")

	flagSet.BoolVar(&example, "e", false, "")
	flagSet.BoolVar(&example, "example", false, "")

	flagSet.BoolVar(&help, "h", false, "")
	flagSet.BoolVar(&help, "help", false, "")

	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-help")
	}

	switch os.Args[1] {
	case "-e", "-example":
		logUtils.PrintExample()
	case "-h", "-help":
		logUtils.PrintUsage()
	default:
		if os.Args[1][0:1] == "-" {
			args := []string{os.Args[0], "gen"}
			args = append(args, os.Args[1:]...)
			os.Args = args
		}

		gen(os.Args)
	}
}

func gen(args []string) {
	flagSet.SetOutput(ioutil.Discard)
	if err := flagSet.Parse(args[2:]); err == nil {
		if vari.HeadSep != "" {
			vari.WithHead = true
		}

		if output != "" {
			ext := strings.ToLower(path.Ext(output))
			if len(ext) > 1 {
				ext = strings.TrimLeft(ext,".")
			}
			if stringUtils.InArray(ext, constant.Formats) {
				format = ext
			}
		}

		if input != "" {
			action.ParseSql(input, output)
		} else {
			action.Generate(deflt, yaml, count, fields, output, format, table)
		}
	} else {
		logUtils.PrintUsage()
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
