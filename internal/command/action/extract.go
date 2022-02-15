package action

import (
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
)

func Extract(files []string) error {
	return scriptUtils.Extract(files)
}
