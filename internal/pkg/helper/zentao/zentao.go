package zentaoHelper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/bitly/go-simplejson"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
)

func GetConfig(baseUrl string) (err error) {
	url := baseUrl + "?mode=getconfig"
	bytes, err := httpUtils.Get(url)
	if err != nil {
		return
	}

	json, _ := simplejson.NewJson(bytes)
	commConsts.ZenTaoVersion, _ = json.Get("version").String()
	commConsts.SessionId, _ = json.Get("sessionID").String()
	commConsts.SessionVar, _ = json.Get("sessionVar").String()
	requestType, _ := json.Get("requestType").String()
	commConsts.RequestType = requestType
	commConsts.RequestFix, _ = json.Get("requestFix").String()

	return
}

func Login(config commDomain.WorkspaceConf) (err error) {
	url := GenApiUrl("tokens", nil, config.Url)

	params := map[string]string{
		"account":  config.Username,
		"password": config.Password,
	}
	bodyBytes, err := httpUtils.Post(url, params)
	if err != nil {
		err = ZentaoLoginErr(err.Error())
		return
	}

	jsn, err := simplejson.NewJson(bodyBytes)
	if err != nil {
		err = ZentaoLoginErr(err.Error())

		return
	}

	if jsn == nil {
		return
	}
	mp, err := jsn.Map()
	if err != nil {
		err = ZentaoLoginErr(err.Error())

		return
	}

	val, ok := mp["token"]
	if ok {
		commConsts.SessionId = val.(string)
		if commConsts.Verbose {
			logUtils.Info(i118Utils.Sprintf("success_to_login"))
		}

	} else {
		err = ZentaoLoginErr(fmt.Sprintf("%#v", mp["error"]))

		return
	}

	return
}

func ListLang() (langs []serverDomain.ZentaoLang, err error) {
	for key, mp := range commConsts.LangMap {
		langs = append(langs, serverDomain.ZentaoLang{Code: key, Name: mp["name"]})
	}

	return
}

func FixUrl(url string) (ret string) {
	regx := regexp.MustCompile(`(http|https):\/\/.+`)
	result := regx.FindStringSubmatch(url)
	if result == nil {
		return
	}

	regx = regexp.MustCompile(`[^:\/]\/`)
	result = regx.FindStringSubmatch(url)
	if result == nil { // without /
		ret = url
	} else {
		index := strings.LastIndex(url, "/")
		if url[index+1:] != "zentao" {
			ret = url[:index+1]
		} else {
			ret = url
		}
	}

	return
}

func FixSiteUrl(orginUrl string) (ret string) {
	u, _ := url.Parse(orginUrl)
	ret = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	ret += "/"
	if len(u.Path) >= 7 && u.Path[:7] == "/zentao" {
		ret = ret + "zentao/"
	}
	return
}
