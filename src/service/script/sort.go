package scriptService

import (
	"fmt"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"regexp"
	"strings"
)

func Sort(cases []string) {
	for _, file := range cases {
		if fileUtils.FileExist(file) {
			script := fileUtils.ReadFile(file)

			regx := regexp.MustCompile(`(?s)\[case\](?U:.*)(\[.*)\[esac\]`)
			arr := regx.FindStringSubmatch(script)

			content := ""
			if len(arr) > 1 {
				checkpoints := arr[1]
				content = commonUtils.RemoveBlankLine(checkpoints)
			}

			level := 0
			groupIdx := 0
			stepIdx := 0

			under := ""

			info := make([]string, 0)
			for _, line := range strings.Split(content, "\n") {
				lineTrim := strings.TrimSpace(line)

				if isGroupIdent(line) {
					info = append(info, lineTrim)

					level = 0
					continue
				}

				isStepsId := isStepsIdent(line)
				isExpectsId := isExpectsIdent(lineTrim)

				isGroupTtl := false
				if !isStepsId && !isExpectsId && isGroupTitle(lineTrim) {
					isGroupTtl = true
				}

				if isStepsId { // [3.1. steps]
					indent := fmt.Sprintf("%d.%d.", groupIdx, stepIdx)
					lineTrim = replaceIndent1(lineTrim, indent, "steps")

					prefix := strings.Repeat(" ", level*2)

					info = append(info, prefix+lineTrim)
					under = "steps"

					continue
				} else if isExpectsId { // [3.1. expects]
					indent := fmt.Sprintf("%d.%d.", groupIdx, stepIdx)
					lineTrim = replaceIndent1(lineTrim, indent, "expects")

					prefix := strings.Repeat(" ", level*2)

					info = append(info, prefix+lineTrim)
					under = "expects"

					if level == 0 {
						groupIdx++
					} else {
						stepIdx++
					}

					continue
				} else if isGroupTtl { // [3. groupe2...]
					info = append(info, lineTrim)
					under = ""
					level = 1

					groupIdx++
					continue
				} else if isStepWithExpect(lineTrim) { // 3.2. step3.2... >> xx
					prefix := strings.Repeat(" ", 2)

					info = append(info, prefix+lineTrim)
					under = ""

					if level == 0 {
						groupIdx++
					} else {
						stepIdx++
					}
					continue
				} else {
					prefix := strings.Repeat(" ", 2)
					if under != "" {
						prefix += strings.Repeat(" ", 2)
					} else {
						if level == 0 {
							groupIdx++
						} else {
							stepIdx++
						}
					}

					info = append(info, prefix+lineTrim)
					continue
				}
			}

			println(strings.Join(info, "\n"))

			// replace info
			// re, _ := regexp.Compile(`(?s)\[case\].*\[esac\]`)
			// content = re.ReplaceAllString(content, "[case]\n"+strings.Join(info, "\n")+"\n\n[esac]")
			//
			// fileUtils.WriteFile(file, content)
			return
		}
	}

	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("success_sort_steps", len(cases)), -1)
}

func isGroupIdent(str string) bool {
	pass, _ := regexp.MatchString(`(?i)\[\s*group\s*\]`, str)
	return pass
}

func isStepsIdent(str string) bool {
	pass, _ := regexp.MatchString(`(?i)\[.*steps\.*\]`, str)
	return pass
}

func isExpectsIdent(str string) bool {
	pass, _ := regexp.MatchString(`(?i)\[.*expects\.*\]`, str)
	return pass
}

func isGroupTitle(str string) bool {
	pass, _ := regexp.MatchString(`(?i)\[.*\]`, str)
	return pass
}

func isStepWithExpect(str string) bool {
	return strings.Contains(str, ">>")
}

func replaceIndent1(str string, indent string, tag string) string {
	regx, _ := regexp.Compile(`\[[\d\.\s]*` + tag + `.*\]`)
	str = regx.ReplaceAllString(str, fmt.Sprintf("[%s %s]", indent, tag))

	return str
}

func replaceIndent2(str string, indent string, tag string) string {
	regx, _ := regexp.Compile(`[\d\.\s]*` + tag + `.*`)
	str = regx.ReplaceAllString(str, fmt.Sprintf("[%s %s]", indent, tag))

	return str
}
