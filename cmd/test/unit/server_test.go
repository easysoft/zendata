package main

import (
	"net/url"
	"testing"

	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	httpUtils "github.com/easysoft/zendata/pkg/utils/http"
)

func TestServer(t *testing.T) {
	urlStr := httpUtils.GenUrl("10.8.0.134", 8848, "?F=field_common&lines=10")
	data := url.Values{}

	defaultContent := fileUtils.ReadFile("demo/default.yaml")
	configContent := fileUtils.ReadFile("demo/test.yaml")

	data.Add("default", defaultContent)
	data.Add("config", configContent)
	//data.Add("lines", "3")

	httpUtils.PostForm(urlStr, data)
}
