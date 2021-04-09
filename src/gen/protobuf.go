package gen

import (
	"fmt"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	shellUtils "github.com/easysoft/zendata/src/utils/shell"
	"github.com/easysoft/zendata/src/utils/vari"
	"path"
	"strconv"
	"strings"
)

const (
	outputDir = "output"
)

func GenerateProtobuf(protoFile string) (ret string) {
	outputDir := generateCls(protoFile)

	converterFile := generateConverter(outputDir)

	return converterFile
}

func generateConverter(dir string) (ret string) {
	srcFile := path.Join(vari.ZdPath, "runtime", "protobuf", "convert.php")
	ret = path.Join(dir, "convert.php")

	content := fileUtils.ReadFile(srcFile)
	content = strings.Replace(content, "${cls_path}", vari.ProtoCls, 1)
	content = strings.Replace(content, "${cls_name}", vari.ProtoCls, 1)

	fileUtils.WriteFile(ret, content)

	return
}

func generateCls(protoFile string) (ret string) {
	outputDir := path.Join(fileUtils.GetAbsoluteDir(protoFile), outputDir)
	fileUtils.RmFile(outputDir)
	fileUtils.MkDirIfNeeded(outputDir)

	platform := commonUtils.GetOs()
	execFile := "protoc"
	if commonUtils.IsWin() {
		platform += fmt.Sprintf("%d", strconv.IntSize)
		execFile += ".exe"
	}

	execFile = path.Join(vari.ZdPath, "runtime", "protobuf", "bin", platform, execFile)

	cmdStr := fmt.Sprintf("%s --php_out=%s %s", execFile, outputDir, protoFile)
	shellUtils.Exec(cmdStr)

	ret = outputDir

	return
}
