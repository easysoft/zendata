package gen

import (
	"fmt"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	shellUtils "github.com/easysoft/zendata/pkg/utils/shell"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	outputDir = "out"
	bufFile   = "data.bin"
)

func GenerateFromProtobuf(protoFile string) (content, pth string) {
	outputDir := generateCls(protoFile)

	convertFile := generateConverter(outputDir)
	content, pth = generateBinData(convertFile)

	return
}

func generateBinData(convertFile string) (content, pth string) {
	dir := filepath.Dir(convertFile)

	phpExeFile := "php"
	if commonUtils.IsWin() { // use build-in php runtime
		phpExeFile = filepath.Join(vari.ZdPath, "runtime", "php", "php7", "php.exe")
	}
	cmdStr := phpExeFile + " convert.php"
	out, _ := shellUtils.ExecInDir(cmdStr, dir)
	if vari.Verbose {
		logUtils.PrintTo(out)
	}

	pth = filepath.Join(dir, bufFile)

	return
}

func generateConverter(dir string) (pth string) {
	srcFile := filepath.Join(vari.ZdPath, "runtime", "protobuf", "convert.php")
	pth = filepath.Join(dir, "convert.php")

	content := fileUtils.ReadFile(srcFile)
	content = strings.ReplaceAll(content, "${cls_name}", vari.ProtoCls)

	fileUtils.WriteFile(pth, content)

	return
}

func generateCls(protoFile string) (ret string) {
	outputDir := filepath.Join(fileUtils.GetAbsoluteDir(protoFile), outputDir)
	fileUtils.RmFile(outputDir)
	fileUtils.MkDirIfNeeded(outputDir)

	platform := commonUtils.GetOs()
	execFile := "protoc"
	if commonUtils.IsWin() {
		platform += fmt.Sprintf("%d", strconv.IntSize)
		execFile += ".exe"
	}

	execFile = filepath.Join(vari.ZdPath, "runtime", "protobuf", "bin", platform, execFile)

	cmdStr := fmt.Sprintf("%s --php_out=%s %s", execFile, outputDir, protoFile)
	shellUtils.Exec(cmdStr)

	ret = outputDir

	return
}
