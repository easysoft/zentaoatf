package middleware

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/multi"
)

/**
 * 验证 jwt
 * @method JwtHandler
 */
func JwtHandler() iris.Handler {
	verifier := multi.NewVerifier()
	verifier.Extractors = []multi.TokenExtractor{multi.FromHeader} // extract token only from Authorization: Bearer $token
	verifier.ErrorHandler = func(ctx *context.Context, err error) {
		ctx.JSON(domain.Response{Code: domain.AuthErr.Code, Data: nil, Msg: domain.AuthErr.Msg})
		//ctx.StopWithError(http.StatusUnauthorized, err)
	} // extract token only from Authorization: Bearer $token
	return verifier.Verify()
}
