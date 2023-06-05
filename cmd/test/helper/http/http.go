package httpHelper

import (
	"encoding/json"
	"errors"
	"github.com/bitly/go-simplejson"
	testConsts "github.com/easysoft/zendata/cmd/test/conf"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

func Get(url, token string) (ret []byte, err error) {
	ret, _, err = GetCheckForward(url, token)
	return
}

func GetCheckForward(url, token string) (ret []byte, isForward bool, err error) {
	logUtils.InfofIfVerbose("===DEBUG=== request: %s", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}

	if !strings.Contains(url, "/tokens") {
		req.Header.Add(testConsts.Token, token)
	}

	resp, err := client.Do(req)
	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}
	defer resp.Body.Close()

	isForward = req.URL.Path != resp.Request.URL.Path

	if !IsSuccessCode(resp.StatusCode) {
		logUtils.InfofIfVerbose(color.RedString("read response failed, StatusCode: %d.", resp.StatusCode))
		err = errors.New(resp.Status)
		return
	}

	ret, err = ioutil.ReadAll(resp.Body)
	logUtils.InfofIfVerbose("===DEBUG=== response: %s", logUtils.ConvertUnicode(ret))

	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("read response failed, error ", err.Error()))
		return
	}

	if len(ret) == 0 {
		return
	}

	jsn, err := simplejson.NewJson(ret)
	if err != nil {
		return
	}
	errMsg, _ := jsn.Get("error").String()
	if strings.ToLower(errMsg) == "unauthorized" {
		err = errors.New(consts.UnAuthorizedErr.Message)
		return
	}

	return
}

func Post(url, token string, data interface{}) (ret []byte, err error) {
	return PostOrPut(url, token, "POST", data)
}
func Put(url, token string, data interface{}) (ret []byte, err error) {
	return PostOrPut(url, token, "PUT", data)
}

func PostOrPut(url, token string, method string, data interface{}) (ret []byte, err error) {
	logUtils.InfofIfVerbose("===DEBUG=== request: %s", url)

	client := &http.Client{}

	dataBytes, err := json.Marshal(data)
	logUtils.InfofIfVerbose("===DEBUG=== data: %s", string(dataBytes))

	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("marshal request failed, error: %s.", err.Error()))
		return
	}

	dataStr := string(dataBytes)

	req, err := http.NewRequest(method, url, strings.NewReader(dataStr))
	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	if !strings.Contains(url, "/tokens") {
		req.Header.Add(testConsts.Token, token)
	}

	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	if !IsSuccessCode(resp.StatusCode) {
		logUtils.InfofIfVerbose(color.RedString("post request return '%s'.", resp.Status))
		err = errors.New(resp.Status)
		return
	}

	ret, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	logUtils.InfofIfVerbose("===DEBUG=== response: %s", logUtils.ConvertUnicode(ret))

	if err != nil {
		logUtils.InfofIfVerbose(color.RedString("read response failed, error: %s.", err.Error()))
		return
	}

	if len(ret) == 0 {
		return
	}

	jsn, err := simplejson.NewJson(ret)
	if err != nil {
		return
	}
	errMsg, _ := jsn.Get("error").String()
	if strings.ToLower(errMsg) == "unauthorized" {
		err = errors.New(consts.UnAuthorizedErr.Message)
		return
	}

	return
}

func IsSuccessCode(code int) (success bool) {
	return code >= 200 && code <= 299
}
