package zd

import (
	"embed"
	"github.com/easysoft/zendata/pkg/utils/common"
	"os"
)

//go:embed res
var resFileSys embed.FS

func ReadResData(path string) (ret []byte, err error) {
	if commonUtils.IsRelease() {
		ret, err = resFileSys.ReadFile(path)
	} else {
		ret, err = os.ReadFile(path)
	}

	return
}
