package serverService

import (
	"fmt"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	serverConst "github.com/easysoft/zendata/src/server/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
)

type CronService struct {
	upgradeService *UpgradeService
}

func NewCronService(upgradeService *UpgradeService) *CronService {
	return &CronService{upgradeService: upgradeService}
}

func (s *CronService) Init() {
	serverUtils.AddTaskFuc(
		i118Utils.I118Prt.Sprintf("check_update"),
		fmt.Sprintf("@every %ds", serverConst.CheckUpgradeInterval),
		func() {
			s.upgradeService.CheckUpgrade()
		},
	)
}
