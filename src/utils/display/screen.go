package display

import (
	"github.com/easysoft/zendata/src/utils/common"
	"github.com/easysoft/zendata/src/utils/shell"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func GetScreenSize() (int, int) {
	var cmd string
	var width int
	var height int

	if commonUtils.IsWin() {
		cmd = "mode" // tested for win7
		out, _ := shellUtils.ExeShell(cmd)

		//out := `设备状态 CON:
		//		---------
		//			行:　       300
		//			列:　　     80
		//			键盘速度:   31
		//			键盘延迟:　 1
		//			代码页:     936
		//`
		myExp := regexp.MustCompile(`CON:\s+[^\s]+\s*(.*?)(\d+)\s\s*(.*?)(\d+)\s`)
		arr := myExp.FindStringSubmatch(out)
		if len(arr) > 4 {
			height, _ = strconv.Atoi(strings.TrimSpace(arr[2]))
			width, _ = strconv.Atoi(strings.TrimSpace(arr[4]))
		}
	} else {
		width, height = noWindowsSize()
	}

	return width, height
}

func noWindowsSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	output := string(out)

	if err != nil {
		return 0, 0
	}
	width, height, err := parseSize(output)

	return width, height
}
func parseSize(input string) (int, int, error) {
	parts := strings.Split(input, " ")
	h, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	w, err := strconv.Atoi(strings.Replace(parts[1], "\n", "", 1))
	if err != nil {
		return 0, 0, err
	}
	return int(w), int(h), nil
}
