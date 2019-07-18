package utils

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type ScreenSizeStruct struct {
	Width  int
	Height int
}

var ScreenSize ScreenSizeStruct

func GetTermSize() ScreenSizeStruct {
	if ScreenSize.Width != 0 {
		return ScreenSize
	}

	var cmd string
	var width int
	var height int

	if IsWin() {
		cmd = "mode" // tested for win7
		out, _ := ExeShell(cmd)

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

	ScreenSize.Width = width
	ScreenSize.Height = height

	return ScreenSize
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

	if width == 0 {
		fmt.Printf("  IDE debug mode will return:")
	}
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
