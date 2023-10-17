package service

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	shellUtils "github.com/easysoft/zendata/pkg/utils/shell"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

const (
	outputDir = "out"
	bufFile   = "data.bin"
)

type ProtobufService struct {
	ResService     *ResService     `inject:""`
	FieldService   *FieldService   `inject:""`
	CombineService *CombineService `inject:""`
	OutputService  *OutputService  `inject:""`
}

func (s *ProtobufService) GenerateProtobuf(protoFile string) (content, pth string) {
	outputDir := s.generateCls(protoFile)

	convertFile := s.generateConverter(outputDir)

	content, pth = s.generateBinData(convertFile)

	return
}

func (s *ProtobufService) generateBinData(convertFile string) (content, pth string) {
	dir := filepath.Dir(convertFile)

	phpExeFile := "php"
	if commonUtils.IsWin() { // use build-in php runtime
		phpExeFile = filepath.Join(vari.WorkDir, "runtime", "php", "php7", "php.exe")
	}
	cmdStr := phpExeFile + " convert.php"
	out, _ := shellUtils.ExecInDir(cmdStr, dir)
	if vari.Verbose {
		logUtils.PrintTo(out)
	}

	pth = filepath.Join(dir, bufFile)

	return
}

func (s *ProtobufService) generateConverter(dir string) (pth string) {
	srcFile := filepath.Join(vari.WorkDir, "runtime", "protobuf", "convert.php")
	pth = filepath.Join(dir, "convert.php")

	content := fileUtils.ReadFile(srcFile)
	content = strings.ReplaceAll(content, "${cls_name}", vari.ProtoCls)

	fileUtils.WriteFile(pth, content)

	return
}

func (s *ProtobufService) generateCls(protoFile string) (ret string) {
	outputDir := filepath.Join(fileUtils.GetAbsoluteDir(protoFile), outputDir)
	fileUtils.RmFile(outputDir)
	fileUtils.MkDirIfNeeded(outputDir)

	platform := commonUtils.GetOs()
	execFile := "protoc"
	if commonUtils.IsWin() {
		platform += fmt.Sprintf("%d", strconv.IntSize)
		execFile += ".exe"
	}

	execFile = filepath.Join(vari.WorkDir, "runtime", "protobuf", "bin", platform, execFile)

	cmdStr := fmt.Sprintf("%s --php_out=%s %s", execFile, outputDir, protoFile)
	shellUtils.Exec(cmdStr)

	ret = outputDir

	return
}
