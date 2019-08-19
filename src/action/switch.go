package action

import configUtils "github.com/easysoft/zentaoatf/src/utils/config"

func SwitchWorkDir(dir string) error {
	configUtils.SetWorkDir(dir, false)

	return nil
}
