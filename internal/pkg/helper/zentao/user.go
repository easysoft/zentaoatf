package zentaoHelper

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	"strings"
)

func GetProfile(config commDomain.WorkspaceConf) (profile commDomain.ZentaoUserProfile, err error) {
	err = Login(config)
	if err != nil {
		return
	}

	url := GenApiUrl("/user", nil, config.Url)
	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(url, err.Error())
		return
	}

	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}
	bytes, err = json.Marshal(jsn.Get("profile"))
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	err = json.Unmarshal(bytes, &profile)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	profile.Avatar = config.Url + strings.TrimLeft(profile.Avatar, "/")

	return
}
