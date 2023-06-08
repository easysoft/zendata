package gen

import (
	"github.com/easysoft/zendata/cmd/command/action"
	"github.com/easysoft/zendata/cmd/test/consts"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"strings"
)

type Generator struct {
	Total        int
	Configs      []string
	ExportFields []string
	Output       string
	Trim         bool
	Human        bool
	Recursive    bool

	DBType   string
	Table    string
	ProtoCls string
}

func New() (g *Generator) {
	g = &Generator{}

	return
}

func (s *Generator) Gen() (out string) {
	if s.Total == 0 {
		s.Total = 10
	}

	vari.GlobalVars.Total = s.Total
	vari.GlobalVars.ExportFields = s.ExportFields
	vari.GlobalVars.Output = s.Output
	vari.GlobalVars.Trim = s.Trim
	vari.GlobalVars.Human = s.Human
	vari.GlobalVars.Recursive = s.Recursive

	vari.GlobalVars.DBType = s.DBType
	vari.GlobalVars.Table = s.Table
	vari.ProtoCls = s.ProtoCls

	action.GenData(s.Configs)

	if len(s.Configs) > 0 && strings.HasSuffix(s.Configs[0], "proto") {
		out = fileUtils.ReadFile(consts.CommandTestFileProtoOut)

	} else if vari.GlobalVars.Output != "" && !strings.HasSuffix(vari.GlobalVars.Output, "xlsx") {
		out = fileUtils.ReadFile(vari.GlobalVars.Output)

	} else {
		out = consts.Buf.String()
	}

	vari.GlobalVars.Trim = false

	return
}

func (s *Generator) SetConfigs(configs []string) (r *Generator) {
	s.Configs = configs

	r = s
	return r
}

func (s *Generator) SetTotal(total int) (r *Generator) {
	s.Total = total

	r = s
	return r
}

func (s *Generator) SetFields(fields string) (r *Generator) {
	arr := strings.Split(fields, ",")
	s.ExportFields = arr

	r = s
	return r
}

func (s *Generator) SetOutput(pth string) (r *Generator) {
	s.Output = pth

	r = s
	return r
}

func (s *Generator) SetDBTable(tp string) (r *Generator) {
	s.DBType = tp

	r = s
	return r
}
func (s *Generator) SetTable(tbl string) (r *Generator) {
	s.Table = tbl

	r = s
	return r
}

func (s *Generator) SetProtoCls(cls string) (r *Generator) {
	s.ProtoCls = cls

	r = s
	return r
}

func (s *Generator) SetTrim(val bool) (r *Generator) {
	s.Trim = val

	r = s
	return r
}

func (s *Generator) SetHuman(val bool) (r *Generator) {
	s.Human = val

	r = s
	return r
}

func (s *Generator) SetRecursive(val bool) (r *Generator) {
	s.Recursive = val

	r = s
	return r
}
