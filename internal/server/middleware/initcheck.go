package middleware

import (
	"net/http"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/internal/server/core/dao"
	"github.com/easysoft/zentaoatf/pkg/domain"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func InitCheck() iris.Handler {
	return func(ctx *context.Context) {
		if dao.GetDB() == nil {
			ctx.StopWithJSON(http.StatusOK,
				domain.Response{Code: commConsts.NeedInitErr.Code, Msg: i118Utils.Sprintf(commConsts.NeedInitErr.Message)})
			return
		}

		lang := ctx.URLParam("lang")
		if lang == "" {
			lang = commConsts.Language
		}
		if lang != commConsts.Language {
			commConsts.Language = lang
			i118Utils.Init(commConsts.Language, commConsts.AppServer)
		}

		ctx.Next()
	}
}
