package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v2"
)

type SqlParseService struct {
}

func (s *SqlParseService) GenYamlFromSql(file string) {
	startTime := time.Now().Unix()

	sql := fileUtils.ReadFile(file)
	statementMap, pkMap, fkMap := s.getCreateStatement(sql)

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

	if vari.GlobalVars.Output != "" {
		vari.GlobalVars.Output = fileUtils.AddSepIfNeeded(vari.GlobalVars.Output)
		outFile := vari.GlobalVars.Output + "keys.yaml"
		fileUtils.WriteFile(outFile, content)
	} else {
		logUtils.PrintTo(content)
	}

	// gen table yaml files
	for tableName, statement := range statementMap {
		createSql := statement

		columns, types := s.getColumnsFromCreateStatement(createSql, nil)

		def := domain.DefSimple{}
		def.Init(tableName, "automated export", "", "1.0")

		for _, col := range columns {
			field := domain.FieldSimple{Field: col}

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
				field.Range = types[col].Rang

				field.Type = types[col].Type
				field.Format = types[col].Format
				field.From = types[col].From
				field.Use = types[col].Use
				field.From = types[col].From
				field.Select = types[col].From
				field.Prefix = types[col].Prefix

				field.Note = types[col].Note

			}

			def.Fields = append(def.Fields, field)
		}

		bytes, _ := yaml.Marshal(&def)
		content := strings.ReplaceAll(string(bytes), "'", "")
		if vari.GlobalVars.Output != "" {
			vari.GlobalVars.Output = fileUtils.AddSepIfNeeded(vari.GlobalVars.Output)
			outFile := vari.GlobalVars.Output + tableName + ".yaml"
			fileUtils.WriteFile(outFile, content)
		} else {
			logUtils.PrintTo(content)
		}
	}

	entTime := time.Now().Unix()
	logUtils.PrintTo(i118Utils.I118Prt.Sprintf("generate_yaml", len(statementMap), vari.GlobalVars.Output, entTime-startTime))
}

func (s *SqlParseService) getCreateStatement(content string) (statementMap map[string]string, pkMap map[string]string, fkMap map[string][2]string) {
	statementMap = map[string]string{}
	pkMap = map[string]string{}
	fkMap = map[string][2]string{}

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

func (s *SqlParseService) getColumnsFromCreateStatement(ddl string, recordsMap map[string][]interface{}) (
	fieldLines []string, fieldInfo map[string]helper.FieldTypeInfo) {
	fieldInfo = map[string]helper.FieldTypeInfo{}

	re := regexp.MustCompile("(?iU)\\s*(\\S+)\\s.*\n")
	results := re.FindAllStringSubmatch(ddl, -1)

	for _, items := range results {
		temp := strings.ToLower(items[0])
		if strings.Contains(temp, " table ") || strings.Contains(temp, " key ") {
			continue
		}

		typ, name, param := s.getFieldData(items)
		fieldInfo[name] = helper.GenerateFieldDefByMetadata(typ, param, name, recordsMap[name])

		fieldLines = append(fieldLines, name)

	}

	return fieldLines, fieldInfo
}

func (s *SqlParseService) getFieldData(item []string) (typ, name string, param string) {
	colName := item[1]
	name = strings.ReplaceAll(colName, "`", "")

	myExp := regexp.MustCompile(colName + `\s([a-zA-Z]+)\((?U:(.*))\)`)
	result := myExp.FindStringSubmatch(item[0])

	if result != nil { // type with length like int(10)
		typ = result[1]
		param = result[2]
	} else {
		typ = strings.Split(strings.Fields(item[0])[1], "(")[0]
	}

	typ = strings.TrimSuffix(typ, ",")

	typ = strings.ToLower(typ)
	name = strings.ToLower(name)

	return
}

func (s *SqlParseService) genKeysYaml(pkMap map[string]string) {
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

	if vari.GlobalVars.Output != "" {
		vari.GlobalVars.Output = fileUtils.AddSepIfNeeded(vari.GlobalVars.Output)
		outFile := vari.GlobalVars.Output + "keys.yaml"
		fileUtils.WriteFile(outFile, content)

	} else {
		logUtils.PrintTo(content)
	}
}

func (s *SqlParseService) genTablesYaml(statementMap map[string]string,
	pkMap map[string]string, fkMap map[string][2]string, recordsMap map[string]map[string][]interface{}) {
	for tableName, statement := range statementMap {
		createStr := statement

		columns, types := s.getColumnsFromCreateStatement(createStr, recordsMap[tableName])

		def := domain.DefSimple{}
		def.Init(tableName, "automated export", "", "1.0")

		for _, col := range columns {
			field := domain.FieldSimple{Field: col}

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
				field.Range = types[col].Rang

				field.Type = types[col].Type
				field.Loop = types[col].Loop
				field.Format = types[col].Format
				field.From = types[col].From
				field.Use = types[col].Use
				field.From = types[col].From
				field.Select = types[col].Select
				field.Prefix = types[col].Prefix

				field.Note = types[col].Note
			}

			def.Fields = append(def.Fields, field)
		}

		bytes, _ := yaml.Marshal(&def)
		content := strings.ReplaceAll(string(bytes), "'", "")
		if vari.GlobalVars.Output != "" {
			vari.GlobalVars.Output = fileUtils.AddSepIfNeeded(vari.GlobalVars.Output)
			outFile := vari.GlobalVars.Output + tableName + ".yaml"
			fileUtils.WriteFile(outFile, content)
		} else {
			logUtils.PrintTo(content)
		}
	}
}

type TableInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}
