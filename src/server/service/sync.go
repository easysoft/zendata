package serverService

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
)

type SyncService struct {
	defService *DefService
	fieldService *FieldService
	rangesService *RangesService
	instancesService *InstancesService
	configService *ConfigService
	excelService *ExcelService
	textService *TextService
	referService *ReferService
	resService *ResService
}

func (s *SyncService) SyncData(mode string) { // TODO: overwrite or not
	files := s.resService.LoadRes("")

	fileMap := map[string][]model.ResFile{}
	for _, fi := range files {
		if fileMap[fi.ResType] == nil {
			fileMap[fi.ResType] = make([]model.ResFile, 0)
		}
		fileMap[fi.ResType] = append(fileMap[fi.ResType], fi)
	}

	s.defService.Sync(fileMap[constant.ResTypeYaml])
	s.rangesService.Sync(fileMap[constant.ResTypeRanges])
	s.instancesService.Sync(fileMap[constant.ResTypeInstances])
	s.configService.Sync(fileMap[constant.ResTypeConfig])

	s.excelService.Sync(fileMap[constant.ResTypeExcel])
	s.textService.Sync(fileMap[constant.ResTypeText])
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
		defService: defService,
		fieldService: fieldService,
		rangesService: rangesService,
		instancesService: instancesService,
		configService: configService,
		excelService: excelService,
		textService: textService,
		referService: referService,
		resService: resService,
	}
}
