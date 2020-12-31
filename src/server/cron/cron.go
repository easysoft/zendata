package cron

import (
	"fmt"
	"github.com/easysoft/zendata/src/server/service"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	serverConst "github.com/easysoft/zendata/src/server/utils/const"
)

type CronService struct {
	upgradeService *serverService.UpgradeService
}

func NewCronService(upgradeService *serverService.UpgradeService) *CronService {
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
