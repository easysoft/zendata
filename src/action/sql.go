package action

import (
	"github.com/easysoft/zendata/src/model"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

func ParseSql(file string, out string) {
	startTime := time.Now().Unix()

	statements := getCreateStatement(file)
	for tableName, statement := range statements {
		columns := getColumnsFromCreateStatement(statement)

		def := model.DefSimple{}
		def.Init(tableName, "automated export", "", "1.0")

		for _, col := range columns {
			field := model.FieldSimple{}
			field.Init(col)
			def.Fields = append(def.Fields, field)
		}

		bytes, _ := yaml.Marshal(&def)
		content := string(bytes)

		if out != "" {
			out = fileUtils.AddSepIfNeeded(out)
			outFile := out + tableName + ".yaml"
			WriteToFile(outFile, content)
		} else {
			logUtils.PrintTo(content)
		}
	}

	entTime := time.Now().Unix()
	logUtils.PrintTo(i118Utils.I118Prt.Sprintf("generate_yaml", len(statements), out, entTime - startTime ))
}

func getCreateStatement(file string) map[string]string {
	statements := map[string]string{}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", file))
		return statements
	}

	re := regexp.MustCompile("(?siU)(CREATE TABLE.*;)")
	arr := re.FindAllString(string(content), -1)
	for _, item := range arr {
		re := regexp.MustCompile("(?i)CREATE TABLE.*\\s+(.+)\\s+\\(") // get table name
		firstLine := strings.Split(item, "\n")[0]
		arr2 := re.FindAllStringSubmatch(firstLine, -1)

		if len(arr2) > 0 && len(arr2[0]) > 1 {
			tableName := arr2[0][1]
			tableName = strings.ReplaceAll(tableName, "`", "")
			statements[tableName] = item
		}
	}

	return statements
}

func getColumnsFromCreateStatement(sent string) []string {
	fieldLines := make([]string, 0)

	re := regexp.MustCompile("(?iU)\\s*(\\S+)\\s.*\n")
	arr := re.FindAllStringSubmatch(string(sent), -1)
	for _, item := range arr {
		line := strings.ToLower(item[0])
		if !strings.Contains(line, " table ") && !strings.Contains(line, " key ") {
			field := item[1]
			field = strings.ReplaceAll(field, "`", "")
			fieldLines = append(fieldLines, field)
		}
	}

	return fieldLines
}