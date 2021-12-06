package action

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseSql(file string, out string) {
	startTime := time.Now().Unix()

	statementMap, pkMap, fkMap := getCreateStatement(file)

	// gen key yaml files
	inst := model.ZdInstances{Title: "keys", Desc: "automated export"}

	for tableName, keyCol := range pkMap {
		item := model.ZdInstancesItem{}

		item.Instance = fmt.Sprintf("%s_%s", tableName, keyCol)
		item.Range = "1-100000"

		inst.Instances = append(inst.Instances, item)
	}

	bytes, _ := yaml.Marshal(&inst)
	content := strings.ReplaceAll(string(bytes), "'-'", "\"\"")

	if out != "" {
		out = fileUtils.AddSepIfNeeded(out)
		outFile := out + "keys.yaml"
		WriteToFile(outFile, content)
	} else {
		logUtils.PrintTo(content)
	}

	// gen table yaml files
	for tableName, statement := range statementMap {
		createStr := statement

		columns, types := getColumnsFromCreateStatement2(createStr)

		def := model.DefSimple{}
		def.Init(tableName, "automated export", "", "1.0")

		for _, col := range columns {
			field := model.FieldSimple{Field: col}

			// pk
			isPk := col == pkMap[tableName]
			fkInfo, isFk := fkMap[col]

			if isPk {
				field.From = "keys.yaml"
				field.Use = fmt.Sprintf("%s_%s", tableName, col)
			} else if isFk {
				field.From = "keys.yaml"
				field.Use = fmt.Sprintf("%s_%s{:1}", fkInfo[0], fkInfo[1])
			} else {
				if types[col] == `"DATE"` || types[col] == `"TIME"` || types[col] == `"YEAR"` || types[col] == `"DATETIME"` || types[col] == `"TIMESTAMP"` {
					field.Range = strconv.Quote("20210821 000000:60")
					field.Type = "timestamp"
					field.Format = strconv.Quote("YY/MM/DD hh:mm:ss")
				} else {
					field.Range = types[col]
				}
			}

			def.Fields = append(def.Fields, field)
		}

		bytes, _ := yaml.Marshal(&def)
		content := strings.ReplaceAll(string(bytes), "'", "")
		if out != "" {
			out = fileUtils.AddSepIfNeeded(out)
			outFile := out + tableName + ".yaml"
			WriteToFile(outFile, content)
		} else {
			logUtils.PrintTo(content)
		}
	}

	entTime := time.Now().Unix()
	logUtils.PrintTo(i118Utils.I118Prt.Sprintf("generate_yaml", len(statementMap), out, entTime-startTime))
}

func getCreateStatement(file string) (statementMap map[string]string, pkMap map[string]string, fkMap map[string][2]string) {
	statementMap = map[string]string{}
	pkMap = map[string]string{}
	fkMap = map[string][2]string{}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", file))
		return
	}

	re := regexp.MustCompile("(?siU)(CREATE TABLE.*;)")
	arr := re.FindAllString(string(content), -1)
	for _, item := range arr {
		re := regexp.MustCompile("(?i)CREATE TABLE.*\\s+(\\S+)\\s+\\(") // get table name
		firstLine := strings.Split(item, "\n")[0]
		arr2 := re.FindAllStringSubmatch(firstLine, -1)

		if len(arr2) > 0 && len(arr2[0]) > 1 {
			tableName := arr2[0][1]
			tableName = strings.ReplaceAll(tableName, "`", "")
			statementMap[tableName] = item

			re3 := regexp.MustCompile("(?i)PRIMARY KEY\\s+\\((\\S+)\\)")
			arr3 := re3.FindAllStringSubmatch(item, -1)
			if len(arr3) > 0 {
				for _, childArr := range arr3 {
					pkMap[tableName] = strings.ReplaceAll(childArr[1], "`", "")
				}
			}

			re4 := regexp.MustCompile("(?i)FOREIGN KEY\\s+\\((\\S+)\\) REFERENCES (\\S+) \\((\\S+)\\)")
			arr4 := re4.FindAllStringSubmatch(item, -1)
			if len(arr4) > 0 {
				for _, childArr := range arr4 {
					col := strings.ReplaceAll(childArr[1], "`", "")
					toTable := strings.ReplaceAll(childArr[2], "`", "")
					toCol := strings.ReplaceAll(childArr[3], "`", "")

					fkMap[col] = [2]string{toTable, toCol}
				}
			}

		}
	}

	return
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

func getColumnsFromCreateStatement2(sent string) (fieldLines []string, fieldInfo map[string]string) {
	fieldLines = make([]string, 0)
	re := regexp.MustCompile("(?iU)\\s*(\\S+)\\s.*\n")
	arr := re.FindAllStringSubmatch(string(sent), -1)
	fieldInfo = map[string]string{}
	for _, item := range arr {
		line := strings.ToLower(item[0])
		if !strings.Contains(line, " table ") && !strings.Contains(line, " key ") {
			field := item[1]
			field = strings.ReplaceAll(field, "`", "")
			fieldInfo[field] = parseFieldInfo(item[0])
			fieldLines = append(fieldLines, field)
		}
	}

	return fieldLines, fieldInfo
}

func parseFieldInfo(fieldType string) (ran string) {
	var ret []string

	fieldType = strings.ToUpper(fieldType)
	if strings.Contains(fieldType, "UNSIGNED") {
		fmt.Println("start 0")
	}

	fieldType = strings.ReplaceAll(strings.Fields(fieldType)[1], ",", "")
	myExp := regexp.MustCompile("([A-Z]*)\\((.*?)\\)")
	ret = myExp.FindStringSubmatch(fieldType)
	if ret != nil {
		ran = judgeFieldType(ret[1], ret[2])
	} else {
		ran = judgeFieldType(fieldType, "")
	}

	return ran
}

func judgeFieldType(fieldType, num string) (ran string) {

	switch fieldType {
	case "TINYINT":
		ran = "-128-127"
		if num == "1" {
			ran = "0-1"
		}
	case "SMALLINT":
		ran = "-32768-32767"
	case "MEDIUMINT":
		ran = "-8388608~8388607"
	case "INT", "INTEGER":
		ran = "-2147483648~2147483647"
	case "FLOAT":
		ran = "123.457"
	case "BIGINT":
		ran = "BIGINT"
	case "DOUBLE":
		ran = "DOUBLE"
	case "DECIMAL":
		ran = "DECIMAL"

	// Date and time type
	case "DATE":
		ran = `"DATE"`
	case "TIME":
		ran = `"TIME"`
	case "YEAR":
		ran = `"YEAR"`
	case "DATETIME":
		ran = `"DATETIME"`
	case "TIMESTAMP":
		ran = `"TIMESTAMP"`
	// String type
	case "CHAR":
		ran = `"CHAR"`
	case "VARCHAR":
		ran = `"VARCHAR"`
	case "TINYBLOB":
		ran = `"TINYBLOB"`
	case "TINYTEXT":
		ran = `"TINYTEXT"`
	case "BLOB":
		ran = `"BLOB"`
	case "TEXT":
		ran = `"TEXT"`
	case "MEDIUMBLOB":
		ran = `"MEDIUMBLOB"`
	case "MEDIUMTEXT":
		ran = `"MEDIUMTEXT"`
	case "LONGBLOB":
		ran = `"LONGBLOB"`
	case "LONGTEXT":
		ran = `"LONGTEXT"`
	default:
	}

	return ran
}
