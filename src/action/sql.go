package action

import (
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func ParseSql(file string, out string) {
	startTime := time.Now().Unix()
	vari.InputDir = filepath.Dir(file) + string(os.PathSeparator)

	sents := getCreateSent(file)
	for _, sent := range sents {
		fields := getFieldsFromCreateSent(sent)
		fields = fields
	}

	entTime := time.Now().Unix()
	logUtils.Screen(i118Utils.I118Prt.Sprintf("generate_yaml", len(sents), out, entTime - startTime ))
}

func getCreateSent(file string) []string {
	sents := make([]string, 0)

	content, err := ioutil.ReadFile(file)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", file))
		return sents
	}

	re := regexp.MustCompile(`(?siU)(CREATE TABLE.*;)`)
	arr := re.FindAllString(string(content), -1)
	for _, item := range arr {
		sents = append(sents,item)
	}

	return sents
}

func getFieldsFromCreateSent(sent string) []string {
	fieldLines := make([]string, 0)

	re := regexp.MustCompile("(?iU)`(.+)`\\s.*,")
	arr := re.FindAllStringSubmatch(string(sent), -1)
	for _, item := range arr {
		fieldLines = append(fieldLines, item[1])
	}

	return fieldLines
}