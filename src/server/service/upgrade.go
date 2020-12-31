package serverService

import (
	"errors"
	"fmt"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	serverConst "github.com/easysoft/zendata/src/server/utils/const"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	configUtils "github.com/easysoft/zendata/src/utils/config"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	shellUtils "github.com/easysoft/zendata/src/utils/shell"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"github.com/inconshreveable/go-update"
	"github.com/mholt/archiver/v3"
	"os"
	"strconv"
	"strings"
)

type UpgradeService struct {
}

func NewUpgradeService() *UpgradeService {
	return &UpgradeService{}
}

func (s *UpgradeService) CheckUpgrade() {
	pth := vari.AgentLogDir + "version.txt"
	serverUtils.Download(serverConst.AgentVersionURL, pth)

	content := strings.TrimSpace(fileUtils.ReadFile(pth))
	version, _ := strconv.ParseFloat(content, 64)
	if vari.Config.Version < version {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("find_new_ver", content), color.FgCyan)

		versionStr := fmt.Sprintf("%.1f", version)
		pass, err := s.DownloadFile(versionStr)
		if pass && err == nil {
			s.RestartVersion(versionStr)
		}
	} else {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("no_new_ver", content), color.FgCyan)
	}
}

func (s *UpgradeService) DownloadFile(version string) (pass bool, err error) {
	os := commonUtils.GetOs()
	if commonUtils.IsWin() {
		os = fmt.Sprintf("%s%d", os, strconv.IntSize)
	}
	url := fmt.Sprintf(serverConst.AgentDownloadURL, version, os)

	dir := vari.AgentLogDir + version

	pth := dir + ".zip"
	err = serverUtils.Download(url, pth)
	if err != nil {
		return
	}

	md5Url := url + ".md5"
	md5Pth := pth + ".md5"
	err = serverUtils.Download(md5Url, md5Pth)
	if err != nil {
		return
	}

	pass = s.checkMd5(pth, md5Pth)
	if !pass {
		msg := i118Utils.I118Prt.Sprintf("fail_md5_check", pth)
		logUtils.PrintToWithColor(msg, color.FgCyan)
		err = errors.New(msg)
		return
	}

	fileUtils.RmDir(dir)
	fileUtils.MkDirIfNeeded(dir)
	err = archiver.Unarchive(pth, dir)

	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_unzip", pth), color.FgCyan)
		return
	}

	return
}

func (s *UpgradeService) RestartVersion(version string) (err error) {
	currExePath := vari.ExeDir + constant.AppName
	bakExePath := currExePath + "_bak"
	newExePath := vari.AgentLogDir + version + constant.PthSep + constant.AppName + constant.PthSep + constant.AppName
	if commonUtils.IsWin() {
		currExePath += ".exe"
		bakExePath += ".exe"
		newExePath += ".exe"
	}
	logUtils.PrintTo(currExePath)

	rd, _ := os.Open(newExePath)
	err = update.Apply(rd, update.Options{OldSavePath: bakExePath})
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_upgrade",
			vari.Config.Version, version, err.Error()), color.FgRed)
	} else {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("success_upgrade",
			vari.Config.Version, version), color.FgCyan)

		// update config file
		vari.Config.Version, _ = strconv.ParseFloat(version, 64)
		configUtils.SaveConfig(vari.Config)
	}

	return
}

func (s *UpgradeService) checkMd5(filePth, md5Pth string) (pass bool) {
	expectVal := fileUtils.ReadFile(md5Pth)
	actualVal := shellUtils.ExeShellWithOutput("md5sum " + filePth + " | awk '{print $1}'")[0]

	return strings.TrimSpace(actualVal) == strings.TrimSpace(expectVal)
}
