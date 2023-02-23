package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/easysoft/zendata/cmd/test/proto/defaults"
	"github.com/easysoft/zendata/cmd/test/proto/dist"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	httpUtils "github.com/easysoft/zendata/pkg/utils/http"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v3"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
)

var (
	flagSet   *flag.FlagSet
	protoFile string
)

func main() {
	flagSet = flag.NewFlagSet("zd", flag.ContinueOnError)
	flagSet.StringVar(&protoFile, "f", "", "")
	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")
	flagSet.Parse(os.Args[1:])

	dir := fileUtils.GetAbsDir(protoFile)
	//yamlPath := getYamlFile(protoFile)
	//request(yamlPath, 8)

	fld := model.DefField{}
	person := new(dist.Person)

	defaults.Set(person, &fld)
	def := convertFieldToDef(fld)

	bytes, _ := yaml.Marshal(def)
	str := stringUtils.ConvertYamlStringToMapFormat(bytes)
	fileUtils.WriteFile(path.Join(dir, "dist/person.yaml"), str)

	bytes, _ = json.Marshal(person)
	fileUtils.WriteFile(path.Join(dir, "dist/person.json"), string(bytes))
}

func convertFieldToDef(fld model.DefField) (def model.DefData) {
	def.Version = "1.0"
	def.Title = fld.Field
	def.Author = "ProtoBuf"

	def.Fields = fld.Fields

	return
}

func getYamlFile(protoFile string) string {
	dir := fileUtils.GetAbsDir(protoFile)
	base := path.Base(protoFile)
	yamlPath := path.Join(dir, strings.Replace(base, ".proto", ".yaml", 1))

	return yamlPath
}

func request(config string, count int) {
	requestWithParams("127.0.0.1", 8848, "", count, "", config)
}
func requestWithParams(host string, port int, fields string, count int, defaultt, config string) {
	urlStr := httpUtils.GenUrl(host, port, fmt.Sprintf("?F=%s&lines=%d", fields, count))
	data := url.Values{}

	defaultContent := fileUtils.ReadFile(defaultt)
	configContent := fileUtils.ReadFile(config)

	data.Add("default", defaultContent)
	data.Add("config", configContent)
	data.Add("lines", "8")

	resp, _ := httpUtils.PostForm(urlStr, data)
	log.Println(resp.([]byte))
}
