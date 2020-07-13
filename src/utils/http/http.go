package httpUtils

import (
	"fmt"
	"github.com/easysoft/zendata/src/utils/vari"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func PostForm(urlStr string, data url.Values) (interface{}, bool) {
	if vari.Verbose {
		log.Print(urlStr)
	}

	resp, err := http.PostForm(urlStr, data)

	if err != nil {
		log.Print(err.Error())
		return nil, false
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err.Error())
		return nil, false
	}

	log.Print(string(body))
	return body, true
}

func GenUrl(server string, path string) string {
	server = UpdateUrl(server)
	url := fmt.Sprintf("%s%s", server, path)
	return url
}

func UpdateUrl(url string) string {
	if strings.LastIndex(url, "/") < len(url)-1 {
		url += "/"
	}
	return url
}