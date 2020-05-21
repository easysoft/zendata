package action

import (
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"os"
	"path/filepath"
	"time"
)

func ParseSql(file string, out string) {
	startTime := time.Now().Unix()
	vari.InputDir = filepath.Dir(file) + string(os.PathSeparator)

	files := make([]string, 0)



	entTime := time.Now().Unix()
	logUtils.Screen(i118Utils.I118Prt.Sprintf("generate_yaml", len(files), out, entTime - startTime ))
}