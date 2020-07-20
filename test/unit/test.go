package main

import (
	"fmt"
	"strings"
)

func main() {
	//urlStr := "http://127.0.0.1:8848/?F=field3&lines=12"
	//file := "test/code/test.yaml"
	//
	//yamlContent, _ := ioutil.ReadFile(file)
	//yamlContent = gen.ReplaceSpecialChars(yamlContent)
	//
	//data := url.Values{"config": {string(yamlContent)}}
	//httpUtils.PostForm(urlStr, data)

	a := "我的"
	b := "dsfdsfdf是我的吗"
	index := strings.Index(a, b)
	fmt.Println(index)
}
