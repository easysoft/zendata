package action

import (
	fileUtils "github.com/easysoft/zendata/src/utils/file"
)

func WriteToFile(path string, content string) {
	fileUtils.WriteFile(path, content)
}
