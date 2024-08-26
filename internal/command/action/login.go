package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
)

func Login() error {
	config := configHelper.LoadByWorkspacePath(commConsts.ZtfDir)
	logUtils.Info(i118Utils.Sprintf("only_test_auth_login", config.Username, config.Url))
	return zentaoHelper.Login(config)
}
