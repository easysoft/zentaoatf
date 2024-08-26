package core

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"

	"github.com/easysoft/zentaoatf/pkg/consts"
	"github.com/easysoft/zentaoatf/pkg/domain"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	authUtils "github.com/easysoft/zentaoatf/internal/pkg/helper/auth"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
)

func Auth() iris.Handler {
	return func(ctx *context.Context) {
		token := authUtils.GetTokenInAuthorization(ctx.GetHeader(consts.Authorization))

		success := false
		if token == serverConfig.CONFIG.AuthToken {
			success = true
		}

		if success {
			ctx.Next()
		} else {
			ctx.StopWithJSON(http.StatusUnauthorized, domain.Response{Code: commConsts.UnAuthorizedErr.Code, Msg: "wrong token"})
		}
	}
}
