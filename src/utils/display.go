package utils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func InitScreenSize() {
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

	Conf.Width = width
	Conf.Height = height
}

func PrintWholeLine(msg string, char string, attr color.Attribute) {
	fixLen := (Conf.Width - len(msg)) / 2
	if fixLen < 0 { // debug mode
		fixLen = 6
	}

	fmt.Println(fixLen)

	fixStr := strings.Repeat(char, fixLen)

	clr := color.New(attr)
	clr.Printf("%s%s%s", fixStr, msg, fixStr)
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
