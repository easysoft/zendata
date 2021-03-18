package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/easysoft/zendata/src/test/proto/defaults"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/emicklei/proto"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

var (
	flagSet *flag.FlagSet
)

func main() {
	flagSet = flag.NewFlagSet("zd", flag.ContinueOnError)
	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")
	flagSet.Parse(os.Args[1:])

	dir, _ := os.Getwd()

	//sample := new(gen.Person)
	//allocate.Zero(&sample)
	//defaults.Set(person)

	sample := new(defaults.Sample)
	defaults.Set(sample)

	bytes, _ := yaml.Marshal(sample)
	fileUtils.WriteFile(path.Join(dir, "src/test/proto/gen/person.yaml"), string(bytes))

	bytes, _ = json.Marshal(sample)
	fileUtils.WriteFile(path.Join(dir, "src/test/proto/gen/person.json"), string(bytes))

	//reader, _ := os.Open(path.Join(dir, "src/test/proto/person.proto"))
	//defer reader.Close()
	//
	//parser := proto.NewParser(reader)
	//definition, _ := parser.Parse()
	//
	//proto.Walk(definition,
	//	proto.WithService(handleService),
	//	proto.WithMessage(handleMessage))
}

func handleService(s *proto.Service) {
	fmt.Println(s.Name)
}

func handleMessage(m *proto.Message) {
	fmt.Println(m.Name)
}
