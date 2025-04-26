package parse_table

import (
	"github.com/easysoft/zendata/pkg/utils/vari"
	"testing"
)

func TestGenMsDefFromFieldName(t *testing.T) {
	vari.GlobalVars.Table = TableByFieldName
	parse(t)
	checkByFieldName(t)
}
