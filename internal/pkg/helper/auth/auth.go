package authUtils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/easysoft/zentaoatf/pkg/consts"

	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
)

func AddBearTokenIfNeeded(req *http.Request) {
	if strings.Index(req.URL.Path, "api.php") > -1 && serverConfig.CONFIG.AuthToken != "" {
		req.Header.Set(consts.Authorization, GenAuthorization())
	}
}

func GenAuthorization() (ret string) {
	return fmt.Sprintf("%s %s", consts.Bearer, serverConfig.CONFIG.AuthToken)
}

func GetTokenInAuthorization(value string) (token string) {
	return strings.Replace(value, consts.Bearer+" ", "", -1)
}
