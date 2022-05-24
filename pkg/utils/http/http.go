package httpUtils

import (
	"fmt"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func GenUrl(server string, port int, path string) string {
	url := fmt.Sprintf("http://%s:%d/%s", server, port, path)
	return url
}
