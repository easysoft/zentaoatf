//go:build windows
// +build windows

package shellUtils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type window struct {
	Row uint16
	Col uint16
}

func WindowSize() window {
	win := window{0, 0}

	cmd := exec.Command("cmd")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CmdLine: "/c mode con"}

	out, _ := cmd.Output()
	if stringUtils.IsGBK(out) {
		out, _ = simplifiedchinese.GBK.NewDecoder().Bytes(out)
	}
	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		getWindowSizeFromLine(line, &win)
		if win.Row > 0 && win.Col > 0 {
			return win
		}
	}

	return win
}

func getWindowSizeFromLine(line string, win *window) {
	fmt.Println(line)
	if strings.Contains(line, "行") || strings.Contains(line, "Lines") {
		re := regexp.MustCompile(`\d+`)
		rs := re.FindAllString(line, -1)
		if len(rs) == 0 {
			return
		}

		row, err := strconv.ParseUint(rs[0], 10, 16)
		if err != nil {
			return
		}

		win.Row = uint16(row)
	}

	if strings.Contains(line, "列") || strings.Contains(line, "Columns") {
		re := regexp.MustCompile(`\d+`)
		rs := re.FindAllString(line, -1)
		if len(rs) == 0 {
			return
		}

		col, err := strconv.ParseUint(rs[0], 10, 16)
		if err != nil {
			return
		}

		win.Col = uint16(col)
	}

	return
}

func GenFullScreenDivider() string {
	divider := "--------------------------------"
	window := WindowSize()
	if window.Col != 0 {
		divider = strings.Repeat("-", int(window.Col))
	}

	return divider
}
