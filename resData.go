package zd

import (
	"embed"
	"os"
	"strings"

	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
)

//go:embed res
var resFileSys embed.FS

func ReadResData(path string) (ret []byte, err error) {
	if commonUtils.IsRelease() {
		path = strings.ReplaceAll(path, "\\", "/")
		ret, err = resFileSys.ReadFile(path)
	} else {
		ret, err = os.ReadFile(path)
	}

	return
}
