package shellUtils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
)

func Exec(cmdStr string) (ret string, err error) {
	return ExecInDir(cmdStr, "")
}

func ExecInDir(cmdStr string, dir string) (ret string, err error) {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	if dir != "" {
		cmd.Dir = dir
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	ret = out.String()
	return
}

func ExeWithOutput(cmdStr string) []string {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	output := make([]string, 0)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return output
	}

	cmd.Start()

	if err != nil {
		output = append(output, fmt.Sprint(err))
		return output
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(strings.TrimRight(line, "\n"))
		output = append(output, line)
	}

	cmd.Wait()

	return output
}
