package genHelper

import (
	"github.com/easysoft/zendata/internal/pkg/domain"
	"strings"
)

func IsSelectExcelWithExpr(field domain.DefField) bool {
	if strings.Index(field.Select, "${") > -1 || strings.Index(field.Where, "${") > -1 {
		return true
	} else {
		return false
	}
}
