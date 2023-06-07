package gen

import (
	"github.com/easysoft/zendata/cmd/command/action"
	"github.com/easysoft/zendata/cmd/test/consts"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"strings"
)

type Generator struct {
	Total        int
	Configs      []string
	ExportFields []string
	Trim         bool
	Human        bool
	Recursive    bool
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
	vari.GlobalVars.Trim = s.Trim
	vari.GlobalVars.Human = s.Human
	vari.GlobalVars.Recursive = s.Recursive

	action.GenData(s.Configs)

	out = consts.Buf.String()

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
