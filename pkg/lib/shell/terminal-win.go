//go:build windows
// +build windows

package shellUtils

import (
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
		if win.Row > 0 && win.Col > 0 {
			return win
		}

		line = strings.TrimSpace(line)

		if strings.Contains(line, "行") || strings.Contains(line, "Lines") {
			re := regexp.MustCompile(`\d+`)
			rs := re.FindAllString(line, -1)
			if len(rs) == 0 {
				continue
			}

			row, err := strconv.ParseUint(rs[0], 10, 16)
			if err != nil {
				continue
			}

			win.Row = uint16(row)
		}

		if strings.Contains(line, "列") || strings.Contains(line, "Columns") {
			re := regexp.MustCompile(`\d+`)
			rs := re.FindAllString(line, -1)
			if len(rs) == 0 {
				continue
			}

			col, err := strconv.ParseUint(rs[0], 10, 16)
			if err != nil {
				continue
			}

			win.Col = uint16(col)
		}
	}

	return win
}
