package zentaoUtils

import (
	"fmt"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/easysoft/zendata/src/utils/vari"
)

func GenApiUri(module string, methd string, param string) string {
	var uri string

	if vari.RequestType == constant.RequestTypePathInfo {
		uri = fmt.Sprintf("%s-%s-%s.json", module, methd, param)
	} else {
		uri = fmt.Sprintf("index.php?m=%s&f=%s&%s&t=json", module, methd, param)
	}

	return uri
}