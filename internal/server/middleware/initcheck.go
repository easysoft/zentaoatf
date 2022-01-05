package middleware

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// InitCheck 初始化检测中间件
func InitCheck() iris.Handler {
	return func(ctx *context.Context) {
		if dao.GetDB() == nil {
			ctx.StopWithJSON(http.StatusOK, domain.Response{Code: domain.NeedInitErr.Code, Data: nil, Msg: domain.NeedInitErr.Msg})
		} else {
			ctx.Next()
		}
		// 处理请求
	}
}
