package httpUtils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/easysoft/zendata/pkg/utils/vari"
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
	body, err := io.ReadAll(resp.Body)
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
