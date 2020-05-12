package action

import (
	"fmt"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	logUtils "github.com/easysoft/zendata/src/utils/log"
)

func WriteToFile(path string, content string) {
	fileUtils.WriteFile(path, content)
	logUtils.Screen(fmt.Sprintf("Test data saved to %s", path))
}
