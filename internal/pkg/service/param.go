package service

import (
	"errors"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

type ParamService struct {
}

func (s *ParamService) CheckParams() (err error) {
	if len(vari.GlobalVars.DefData.Fields) == 0 {
		err = errors.New("")
	} else if vari.GlobalVars.DefData.Type == constant.DefTypeArticle && vari.GlobalVars.OutputFile == "" { // gen article
		errMsg := i118Utils.I118Prt.Sprintf("gen_article_must_has_out_param")
		logUtils.PrintErrMsg(errMsg)
		err = errors.New(errMsg)
	}

	return
}

func (s *ParamService) FixTotalNum() {
	if vari.GlobalVars.DefData.Type == constant.DefTypeArticle {
		vari.GlobalVars.Total = 1
	}

	if vari.GlobalVars.Total < 0 {
		vari.GlobalVars.Total = constant.DefaultNumber
	}
}
