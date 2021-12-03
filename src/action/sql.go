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

		columns := getColumnsFromCreateStatement2(createStr)
		// todo
		// handle fieldTypes
		//rans := getFieldType(infos)

		def := model.DefSimple{}
		def.Init(tableName, "automated export", "", "1.0")

		for m, col := range columns {
			field := model.FieldSimple{Field: m}

			// pk
			isPk := m == pkMap[tableName]
			fkInfo, isFk := fkMap[m]

			if isPk {
				field.From = "keys.yaml"
				field.Use = fmt.Sprintf("%s_%s", tableName, m)
			} else if isFk {
				field.From = "keys.yaml"
				field.Use = fmt.Sprintf("%s_%s{:1}", fkInfo[0], fkInfo[1])
			} else {
				field.Range = col
			}

			def.Fields = append(def.Fields, field)
		}

		bytes, _ := yaml.Marshal(&def)
		content := strings.ReplaceAll(string(bytes), "'-'", "\"\"")

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

func getColumnsFromCreateStatement2(sent string) (fieldInfo map[string]string) {
	//fieldLines := make([]string, 0)
	re := regexp.MustCompile("(?iU)\\s*(\\S+)\\s.*\n")
	arr := re.FindAllStringSubmatch(string(sent), -1)
	//fieldTypes := make([]string, 0)
	fieldInfo = map[string]string{}
	for _, item := range arr {
		line := strings.ToLower(item[0])
		if !strings.Contains(line, " table ") && !strings.Contains(line, " key ") {
			field := item[1]
			fieldType := strings.Split(strings.TrimSpace(strings.Split(strings.Split(item[0], field)[1], ",")[0]), " ")[0]
			//fieldTypes = append(fieldTypes, fieldType)
			field = strings.ReplaceAll(field, "`", "")
			fieldInfo[field] = judge(fieldType)
			//fieldLines = append(fieldLines, field)
		}
	}

	return fieldInfo
}

func judge(fieldType string) (ran string) {
	ran = ""
	typ := strings.ToUpper(strings.Split(fieldType, "(")[0])
	switch typ {
	// value type
	case "TINYINT":
		ran = "-128-127"
		fmt.Println("TINYINT")
	case "SMALLINT":
		ran = "-32768-32767"
		fmt.Println("SMALLINT")
	case "MEDIUMINT":
		ran = "-8388608~8388607"
		fmt.Println("MEDIUMINT")
	case "INT", "INTEGER":
		ran = "-2147483648~2147483647"
		fmt.Println("INT/INTEGER")
	case "FLOAT":
		ran = "123.457"
		fmt.Println("FLOAT")
	case "BIGINT":
		ran = "BIGINT"
		fmt.Println("BIGINT")
	case "DOUBLE":
		ran = "DOUBLE"
		fmt.Println("DOUBLE")
	case "DECIMAL":
		ran = "DECIMAL"
		fmt.Println("DECIMAL")

	// Date and time type
	case "DATE":
		ran = "2008-12-2"
		fmt.Println("DATE")
	case "TIME":
		ran = "12:25:36"
		fmt.Println("TIME")
	case "YEAR":
		ran = "1901/2155"
		fmt.Println("YEAR")
	case "DATETIME":
		ran = "20210101 000000-20210101 230000:60"
		fmt.Println("DATETIME")
	case "TIMESTAMP":
		ran = "1970-01-01 00:00:00/2038"
		fmt.Println("TIMESTAMP")
	// String type
	case "CHAR":
		ran = "CHAR"
		fmt.Println("CHAR")
	case "VARCHAR":
		ran = "VARCHAR"
		fmt.Println("VARCHAR")
	case "TINYBLOB":
		ran = "TINYBLOB"
		fmt.Println("TINYBLOB")
	case "TINYTEXT":
		ran = "TINYTEXT"
		fmt.Println("TINYTEXT")
	case "BLOB":
		ran = "BLOB"
		fmt.Println("BLOB")
	case "TEXT":
		ran = "TEXT"
		fmt.Println("TEXT")
	case "MEDIUMBLOB":
		ran = "MEDIUMBLOB"
		fmt.Println("MEDIUMBLOB")
	case "MEDIUMTEXT":
		ran = "MEDIUMTEXT"
		fmt.Println("MEDIUMTEXT")
	case "LONGBLOB":
		ran = "LONGBLOB"
		fmt.Println("LONGBLOB")
	case "LONGTEXT":
		ran = "LONGTEXT"
		fmt.Println("LONGTEXT")
	default:
		ran = "0-255"
		fmt.Println("other types")
	}
	return ran
}

func getFieldType(fieldTypes []string) (rans []string) {
	ran := "0-255"
	for _, v := range fieldTypes {
		typ := strings.ToUpper(strings.Split(v, "(")[0])
		switch typ {
		// value type
		case "TINYINT":
			ran = "-128-127"
			fmt.Println("TINYINT")
		case "SMALLINT":
			ran = "-32768-32767"
			fmt.Println("SMALLINT")
		case "MEDIUMINT":
			ran = "-8388608~8388607"
			fmt.Println("MEDIUMINT")
		case "INT", "INTEGER":
			ran = "-2147483648~2147483647"
			fmt.Println("INT/INTEGER")
		case "FLOAT":
			ran = "123.457"
			fmt.Println("FLOAT")
		case "BIGINT":
			ran = "BIGINT"
			fmt.Println("BIGINT")
		case "DOUBLE":
			ran = "DOUBLE"
			fmt.Println("DOUBLE")
		case "DECIMAL":
			ran = "DECIMAL"
			fmt.Println("DECIMAL")

		// Date and time type
		case "DATE":
			ran = "2008-12-2"
			fmt.Println("DATE")
		case "TIME":
			ran = "12:25:36"
			fmt.Println("TIME")
		case "YEAR":
			ran = "1901/2155"
			fmt.Println("YEAR")
		case "DATETIME":
			ran = "20210101 000000-20210101 230000:60"
			fmt.Println("DATETIME")
		case "TIMESTAMP":
			ran = "1970-01-01 00:00:00/2038"
			fmt.Println("TIMESTAMP")
		// String type
		case "CHAR":
			ran = "CHAR"
			fmt.Println("CHAR")
		case "VARCHAR":
			ran = "VARCHAR"
			fmt.Println("VARCHAR")
		case "TINYBLOB":
			ran = "TINYBLOB"
			fmt.Println("TINYBLOB")
		case "TINYTEXT":
			ran = "TINYTEXT"
			fmt.Println("TINYTEXT")
		case "BLOB":
			ran = "BLOB"
			fmt.Println("BLOB")
		case "TEXT":
			ran = "TEXT"
			fmt.Println("TEXT")
		case "MEDIUMBLOB":
			ran = "MEDIUMBLOB"
			fmt.Println("MEDIUMBLOB")
		case "MEDIUMTEXT":
			ran = "MEDIUMTEXT"
			fmt.Println("MEDIUMTEXT")
		case "LONGBLOB":
			ran = "LONGBLOB"
			fmt.Println("LONGBLOB")
		case "LONGTEXT":
			ran = "LONGTEXT"
			fmt.Println("LONGTEXT")
		default:
			ran = "0-255"
			fmt.Println("other types")
		}
		rans = append(rans, ran)
	}
	return rans
}
