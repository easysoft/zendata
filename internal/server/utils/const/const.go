package serverConst

import consts "github.com/easysoft/zendata/internal/pkg/const"

const (
	HeartBeatInterval    = 60
	CheckUpgradeInterval = 30

	AgentRunTime = 30 * 60
	AgentLogDir  = "log"

	QiNiuURL         = "https://dl.cnezsoft.com/" + consts.AppName + "/"
	AgentVersionURL  = QiNiuURL + "version.txt"
	AgentDownloadURL = QiNiuURL + "%s/%s/" + consts.AppName + ".zip"
)
