package serverService

import (
	"fmt"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	serverConst "github.com/easysoft/zendata/src/server/utils/const"
)

type CronService struct {
	upgradeService *UpgradeService
}

func NewCronService(upgradeService *UpgradeService) *CronService {
	return &CronService{upgradeService: upgradeService}
}

func (s *CronService) Init() {
	serverUtils.AddTaskFuc(
		"CheckUpdate",
		fmt.Sprintf("@every %ds", serverConst.CheckUpgradeInterval),
		func() {
			s.upgradeService.CheckUpgrade()
		},
	)
}
