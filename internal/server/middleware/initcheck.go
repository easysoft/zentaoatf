package middleware

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	"net/http"

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

		//lang := ctx.URLParam("lang")
		//if lang != commConsts.Language {
		//	commConsts.Language = lang
		//	i118Utils.Init(commConsts.Language, commConsts.AppServer)
		//}

		ctx.Next()
	}
}
