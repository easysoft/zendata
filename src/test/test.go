package main

import (
	"encoding/base64"
	"log"
	"net/url"
)

func main() {
	str := "http://baidu.com"
	log.Println(base64.StdEncoding.EncodeToString([]byte(str)))
	log.Println(url.QueryEscape(str))
}
