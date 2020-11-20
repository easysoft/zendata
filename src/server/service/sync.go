package serverService

import "github.com/easysoft/zendata/src/model"

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

	defs := fileMap["yaml"]

	s.defService.Sync(defs)
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
