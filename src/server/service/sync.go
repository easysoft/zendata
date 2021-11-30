package serverService

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
)

type SyncService struct {
	DefService       *DefService       `inject:""`
	FieldService     *FieldService     `inject:""`
	RangesService    *RangesService    `inject:""`
	InstancesService *InstancesService `inject:""`
	ConfigService    *ConfigService    `inject:""`
	ExcelService     *ExcelService     `inject:""`
	TextService      *TextService      `inject:""`
	ReferService     *ReferService     `inject:""`
	ResService       *ResService       `inject:""`
}

func (s *SyncService) SyncData(mode string) { // TODO: overwrite or not
	files := s.ResService.LoadRes("")

	fileMap := map[string][]model.ResFile{}
	for _, fi := range files {
		if fileMap[fi.ResType] == nil {
			fileMap[fi.ResType] = make([]model.ResFile, 0)
		}
		fileMap[fi.ResType] = append(fileMap[fi.ResType], fi)
	}

	s.DefService.Sync(fileMap[constant.ResTypeYaml])
	s.RangesService.Sync(fileMap[constant.ResTypeRanges])
	s.InstancesService.Sync(fileMap[constant.ResTypeInstances])
	s.ConfigService.Sync(fileMap[constant.ResTypeConfig])

	s.ExcelService.Sync(fileMap[constant.ResTypeExcel])
	s.TextService.Sync(fileMap[constant.ResTypeText])
}

func NewSyncService(
	defService *DefService,
	fieldService *FieldService,
	rangesService *RangesService,
	instancesService *InstancesService,
	configService *ConfigService,
	excelService *ExcelService,
	textService *TextService,
	referService *ReferService,
	resService *ResService) *SyncService {
	return &SyncService{
		DefService:       defService,
		FieldService:     fieldService,
		RangesService:    rangesService,
		InstancesService: instancesService,
		ConfigService:    configService,
		ExcelService:     excelService,
		TextService:      textService,
		ReferService:     referService,
		ResService:       resService,
	}
}
