package zd

import (
	"embed"
	"github.com/easysoft/zendata/pkg/utils/common"
	"io/fs"
	"os"
)

//go:embed ui/dist
var uiFileSys embed.FS

func GetUiFileSys() (ret fs.FS, err error) {
	if commonUtils.IsRelease() {
		ret, err = fs.Sub(uiFileSys, "ui/dist")
	} else {
		ret = os.DirFS("ui/dist")
	}

	return
}
