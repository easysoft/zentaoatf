package zentaoHelper

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	httpUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/http"
)

func GetProfile(config commDomain.WorkspaceConf) (profile commDomain.ZentaoUserProfile, err error) {
	err = Login(config)
	if err != nil {
		return
	}

	url := GenApiUrl("/user", nil, config.Url)
	bytes, err := httpUtils.Get(url)
	if err != nil {
		return
	}

	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		return
	}
	bytes, err = json.Marshal(jsn.Get("profile"))
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &profile)
	profile.Avatar = config.Url + profile.Avatar

	return
}
