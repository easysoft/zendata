package action

import (
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
)

func WriteToFile(path string, content string) {
	fileUtils.WriteFile(path, content)
	logUtils.Screen(i118Utils.I118Prt.Sprintf("save_testdata", path))
}
