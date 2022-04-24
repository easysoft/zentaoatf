package zentaoHelper

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	httpUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/http"
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
		err = ZentaoRequestErr(err.Error())
		return
	}

	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}
	bytes, err = json.Marshal(jsn.Get("profile"))
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	err = json.Unmarshal(bytes, &profile)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	profile.Avatar = config.Url + strings.TrimLeft(profile.Avatar, "/")

	return
}
