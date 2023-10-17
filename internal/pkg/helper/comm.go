package helper

import (
	"path/filepath"
	"strings"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
)

func IsFromProtobuf(file string) bool {
	return strings.ToLower(filepath.Ext(file)) == "."+consts.FormatProto
}

func IsSelectExcelWithExpr(field domain.DefField) bool {
	if strings.Index(field.Select, "${") > -1 || strings.Index(field.Where, "${") > -1 {
		return true
	} else {
		return false
	}
}
