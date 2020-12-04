package serverService

import (
	"fmt"
	"github.com/easysoft/zendata/src/action"
	"github.com/easysoft/zendata/src/server/repo"
	constant "github.com/easysoft/zendata/src/utils/const"
)

type PreviewService struct {
	defRepo   *serverRepo.DefRepo
	fieldRepo *serverRepo.FieldRepo
	referRepo *serverRepo.ReferRepo
}

func (s *PreviewService) PreviewDefData(defId uint) (data string) {
	def, _ := s.defRepo.Get(defId)

	lines := action.Generate("", def.Path, "", constant.FormatData, "")
	data = s.linesToStr(lines)

	return
}
func (s *PreviewService) PreviewFieldData(fieldId uint) (data string) {
	field, _ := s.fieldRepo.Get(fieldId)
	fields := field.Field
	def, _ := s.defRepo.Get(field.DefID)

	lines := action.Generate("", def.Path, fields, constant.FormatData, "")
	data = s.linesToStr(lines)

	return
}
func (s *PreviewService) linesToStr(lines []interface{}) (data string) {
	for index, line := range lines {
		if index > 0 {
			data += "<br/>"
		}
		data += fmt.Sprint(line)
	}

	return
}

func NewPreviewService(defRepo *serverRepo.DefRepo, fieldRepo *serverRepo.FieldRepo) *PreviewService {
	return &PreviewService{defRepo: defRepo, fieldRepo: fieldRepo}
}
