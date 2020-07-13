package main

import (
	"github.com/easysoft/zendata/src/gen"
	httpUtils "github.com/easysoft/zendata/src/utils/http"
	"io/ioutil"
	"net/url"
)

func main() {
	urlStr := "http://127.0.0.1:8848/?F=field3&lines=12"
	file := "test/code/test.yaml"

	yamlContent, _ := ioutil.ReadFile(file)
	yamlContent = gen.ReplaceSpecialChars(yamlContent)

	data := url.Values{"config": {string(yamlContent)}}
	httpUtils.PostForm(urlStr, data)
}
