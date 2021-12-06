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

				if types[col].(FieldTypeInfo).fieldType == "DATE" || types[col].(FieldTypeInfo).fieldType == "TIME" || types[col].(FieldTypeInfo).fieldType == "YEAR" || types[col].(FieldTypeInfo).fieldType == "DATETIME" || types[col].(FieldTypeInfo).fieldType == "TIMESTAMP" {
					field.Range = strconv.Quote("20210821 000000:60")
					field.Format = strconv.Quote("YY/MM/DD hh:mm:ss")
					field.Type = "timestamp"
				} else if types[col].(FieldTypeInfo).fieldType == "CHAR" || types[col].(FieldTypeInfo).fieldType == "VARCHAR" || types[col].(FieldTypeInfo).fieldType == "TINYTEXT" || types[col].(FieldTypeInfo).fieldType == "TEXT" || types[col].(FieldTypeInfo).fieldType == "MEDIUMTEXT" || types[col].(FieldTypeInfo).fieldType == "LONGTEXT" {
					field.Range = types[col].(FieldTypeInfo).rang
					field.Loop = "'3'" // default value of loop
					field.Loopfix = "_"
				} else {
					field.Range = types[col].(FieldTypeInfo).rang
				}
				field.Note = types[col].(FieldTypeInfo).note

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

func getColumnsFromCreateStatement2(sent string) (fieldLines []string, fieldInfo map[string]interface{}) {
	fieldLines = make([]string, 0)
	re := regexp.MustCompile("(?iU)\\s*(\\S+)\\s.*\n")
	arr := re.FindAllStringSubmatch(string(sent), -1)
	fieldInfo = make(map[string]interface{})

	for _, item := range arr {
		line := strings.ToLower(item[0])
		if !strings.Contains(line, " table ") && !strings.Contains(line, " key ") {
			fieldTmp := item[1]
			field := strings.ReplaceAll(fieldTmp, "`", "")
			fieldInfo[field] = parseFieldInfo(strings.ToUpper(fieldTmp), item[0])
			fieldLines = append(fieldLines, field)
		}
	}

	return fieldLines, fieldInfo
}

func parseFieldInfo(fieldTmp, fieldTypeInfo string) (FieldTypeInfoIns FieldTypeInfo) {

	var ret []string
	isUnsigned := false
	fieldTypeInfo = strings.ToUpper(fieldTypeInfo)
	if strings.Contains(fieldTypeInfo, "UNSIGNED") { // judge unsigned
		isUnsigned = true
	}

	myExp := regexp.MustCompile(fieldTmp + `\s([A-Z]+)\(([^,]*),?([^,]*)\)`)
	ret = myExp.FindStringSubmatch(fieldTypeInfo)
	if ret != nil {
		FieldTypeInfoIns = judgeFieldType(ret[1], ret[2], isUnsigned)
	} else {
		fieldType := strings.Split(strings.Fields(fieldTypeInfo)[1], "(")[0]
		FieldTypeInfoIns = judgeFieldType(fieldType, "", isUnsigned)
	}

	return FieldTypeInfoIns
}

type FieldTypeInfo struct {
	fieldType string
	note      string
	rang      string
}

func judgeFieldType(fieldType, num string, isUnsigned bool) (FieldTypeInfoIns FieldTypeInfo) {
	ran := ""
	note := ""
	switch fieldType {
	// integer
	case "BIT":
		ran = "0,1"
	case "TINYINT":
		if num == "1" {
			ran = "0,1"
		} else {
			ran = "0-255"
		}
	case "SMALLINT":
		ran = "0-65535"
	case "MEDIUMINT":
		ran = "0-65535"
		note = `"MEDIUMINT [0,2^24-1]"`
	case "INT", "INTEGER":
		ran = "0-100000"
		note = `"INI [0,2^32-1]"`
	case "BIGINT":
		ran = "0-100000"
		note = `"BIGINT [0,2^64-1]"`
	// floating-point
	case "FLOAT":
		ran = "1.01-99.99:0.01"
		note = `"FLOAT"`
	case "DOUBLE":
		ran = "1.01-99.99:0.01"
		note = `"DOUBLE"`
	// fixed-point
	case "DECIMAL":
		ran = "123.45"
		note = `"DECIMAL"`
	// character string
	case "CHAR":
		ran = `"a-z"`
	case "VARCHAR":
		ran = `"VARCHAR"`
	case "TINYTEXT":
		ran = `"TINYTEXT"`
	case "TEXT":
		ran = `"TEXT"`
	case "MEDIUMTEXT":
		ran = `"MEDIUMTEXT"`
	case "LONGTEXT":
		ran = `"LONGTEXT"`
	// binary data
	case "TINYBLOB":
		ran = `"TINYBLOB"`
	case "BLOB":
		ran = `"BLOB"`
	case "MEDIUMBLOB":
		ran = `"MEDIUMBLOB"`
	case "LONGBLOB":
		ran = `"LONGBLOB"`
	case "BINARY":
		ran = "BINARY"
	case "VARBINARY":
		ran = "VARBINARY"
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
	// other type
	case "ENUM":
		ran = `"ENUM"`
	case "SET":
		ran = `"SET"`
	case "GEOMETRY":
		ran = `"GEOMETRY"`
	case "POINT":
		ran = `"POINT"`
	case "LINESTRING":
		ran = `"LINESTRING"`
	case "POLYGON":
		ran = `"POLYGON"`
	case "MULTIPOINT":
		ran = `"MULTIPOINT"`
	case "MULTILINESTRING":
		ran = `"MULTILINESTRING"`
	case "MULTIPOLYGON":
		ran = `"MULTIPOLYGON"`
	case "GEOMETRYCOLLECTION":
		ran = `"GEOMETRYCOLLECTION"`
	default:
	}

	FieldTypeInfoIns.rang = ran
	FieldTypeInfoIns.note = note
	FieldTypeInfoIns.fieldType = fieldType

	return FieldTypeInfoIns
}
