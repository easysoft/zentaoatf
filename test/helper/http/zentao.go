package httpHelper

import (
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
)

func Login() (err error) {
	config := commDomain.WorkspaceConf{
		Url:      constTestHelper.ZtfUrl,
		Username: constTestHelper.ZentaoUsername,
		Password: constTestHelper.ZentaoPassword,
	}

	zentaoHelper.Login(config)

	return
}
