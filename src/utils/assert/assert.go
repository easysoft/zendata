package assertUtils

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	langUtils "github.com/easysoft/zentaoatf/src/utils/lang"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func GetCaseByDirAndFile(files []string) []string {
	cases := make([]string, 0)

	for _, file := range files {
		GetAllScriptsInDir(file, &cases)
	}

	return cases
}

func GetAllScriptsInDir(filePthParam string, files *[]string) error {
	sep := string(os.PathSeparator)

	if !fileUtils.IsDir(filePthParam) { // first call, param is file
		regx := langUtils.GetSupportLanguageExtRegx()

		pass, _ := regexp.MatchString(`.*\.`+regx+`$`, filePthParam)

		if pass {
			pass := zentaoUtils.CheckFileIsScript(filePthParam)
			if pass {
				*files = append(*files, filePthParam)
			}
		}

		return nil
	}

	filePthParam = fileUtils.AbosutePath(filePthParam)

	dir, err := ioutil.ReadDir(filePthParam)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if commonUtils.IngoreFile(name) {
			continue
		}

		if fi.IsDir() { // 目录, 递归遍历
			GetAllScriptsInDir(filePthParam+name+sep, files)
		} else {
			path := filePthParam + name
			regx := langUtils.GetSupportLanguageExtRegx()
			pass, _ := regexp.MatchString("^*.\\."+regx+"$", path)

			if pass {
				pass = zentaoUtils.CheckFileIsScript(path)
				if pass {
					*files = append(*files, path)
				}
			}
		}
	}

	return nil
}

func GetScriptByIdsInDir(dirPth string, idMap map[int]string, files *[]string) error {
	dirPth = fileUtils.AbosutePath(dirPth)

	sep := string(os.PathSeparator)

	if commonUtils.IngoreFile(dirPth) {
		return nil
	}

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if fi.IsDir() { // 目录, 递归遍历
			GetScriptByIdsInDir(dirPth+name+sep, idMap, files)
		} else {
			regx := langUtils.GetSupportLanguageExtRegx()
			pass, _ := regexp.MatchString("^*.\\."+regx+"$", name)

			if !pass {
				continue
			}

			path := dirPth + name

			pass, id, _, _ := zentaoUtils.GetCaseInfo(path)
			if pass {
				_, ok := idMap[id]

				if ok {
					*files = append(*files, path)
				}
			}
		}
	}

	return nil
}

func GetCaseIdsInSuiteFile(name string, fileIdMap *map[int]string) {
	content := fileUtils.ReadFile(name)

	for _, line := range strings.Split(content, "\n") {
		idStr := strings.TrimSpace(line)
		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err == nil {
			(*fileIdMap)[id] = ""
		}
	}
}

func GetFailedCasesDirectlyFromTestResult(resultFile string) []string {
	cases := make([]string, 0)

	extName := path.Ext(resultFile)

	if extName == "."+constant.ExtNameResult {
		resultFile = strings.Replace(resultFile, extName, "."+constant.ExtNameJson, -1)
	}

	content := fileUtils.ReadFile(resultFile)

	var report model.TestReport
	json.Unmarshal([]byte(content), &report)

	for _, cs := range report.FuncResult {
		if cs.Status != "pass" {
			cases = append(cases, cs.Path)
		}
	}

	return cases
}

func GetScriptType(scripts []string) []string {
	vari.ScriptExtToNameMap = langUtils.GetExtToNameMap()

	exts := make([]string, 0)
	for _, script := range scripts {
		ext := path.Ext(script)
		if ext != "" {
			ext = ext[1:]
			name := vari.ScriptExtToNameMap[ext]

			if !stringUtils.FindInArr(name, exts) {
				exts = append(exts, name)
			}
		}
	}

	return exts
}
