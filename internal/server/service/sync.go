package serverService

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
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

func (s *SyncService) SyncData() {
	files := s.ResService.LoadRes("")

	fileMap := map[string][]domain.ResFile{}
	for _, fi := range files {
		if fileMap[fi.ResType] == nil {
			fileMap[fi.ResType] = make([]domain.ResFile, 0)
		}

		fileMap[fi.ResType] = append(fileMap[fi.ResType], fi)
	}

	s.DefService.Sync(fileMap[consts.ResTypeYaml])
	s.RangesService.Sync(fileMap[consts.ResTypeRanges])
	s.InstancesService.Sync(fileMap[consts.ResTypeInstances])
	s.ConfigService.Sync(fileMap[consts.ResTypeConfig])

	s.ExcelService.Sync(fileMap[consts.ResTypeExcel])
	s.TextService.Sync(fileMap[consts.ResTypeText])
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
