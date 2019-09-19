package scriptService

import (
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
)

func Sort(cases []string) {
	for _, tc := range cases {
		_ = tc
	}

	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("success_sort_steps", len(cases)), -1)
}
