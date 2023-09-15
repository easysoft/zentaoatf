package zentaoHelper

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"

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
	bodyBytes, httpError := httpUtils.Post(url, params)
	if httpError != nil && string(bodyBytes) == "" {
		err = ZentaoLoginErr(httpError.Error())
		return
	}

	jsn, err := simplejson.NewJson(bodyBytes)
	if err != nil {
		if httpError != nil {
			err = ZentaoLoginErr(err.Error())
			return
		}
		err = ZentaoLoginErr(err.Error())

		return
	}

	if jsn == nil {
		return
	}
	mp, err := jsn.Map()
	if err != nil {
		if httpError != nil {
			err = ZentaoLoginErr(err.Error())
			return
		}
		err = ZentaoLoginErr(err.Error())

		return
	}

	if httpError != nil {
		_, ok := mp["error"]
		if ok {
			err = ZentaoLoginErr(strings.TrimRight(mp["error"].(string), ".。"))
			return
		}
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

func LoginSilently(config commDomain.WorkspaceConf) (err error) {
	url := GenApiUrl("tokens", nil, config.Url)

	params := map[string]string{
		"account":  config.Username,
		"password": config.Password,
	}
	bodyBytes, err := httpUtils.Post(url, params)
	if err != nil {
		return
	}

	jsn, err := simplejson.NewJson(bodyBytes)
	if err != nil {
		return
	}

	if jsn == nil {
		return
	}
	mp, err := jsn.Map()
	if err != nil {
		return
	}

	val, ok := mp["token"]
	if ok {
		commConsts.SessionId = val.(string)
		if commConsts.Verbose {
			logUtils.Info(i118Utils.Sprintf("success_to_login"))
		}
	} else {
		if commConsts.Verbose {
			err = ZentaoLoginErr(fmt.Sprintf("%#v", mp["error"]))
		}

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

func GetSiteUrl(srcUrl string) (baseUrl, version string, err error) {
	url1, url2 := FixSiteUrl(srcUrl)

	version = getZentaoVersion(url1)
	if version != "" {
		baseUrl = url1
	} else {
		version = getZentaoVersion(url2)
		if version != "" {
			baseUrl = url2
		}
	}

	if baseUrl == "" {
		err = errors.New(i118Utils.Sprintf("wrong_zentao_url"))
		return
	}

	return
}

func FixSiteUrl(originalUrl string) (url1, url2 string) {
	u, _ := url.Parse(originalUrl)
	url1 = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	url1 += "/"

	pth := strings.Replace(originalUrl, url1, "", -1)
	pth = strings.TrimLeft(pth, "/")
	arr := strings.Split(pth, "/")
	if len(arr) >= 1 {
		url2 = url1 + arr[0]
	}

	url1 = fileUtils.AddUrlPathSepIfNeeded(url1)
	url2 = fileUtils.AddUrlPathSepIfNeeded(url2)

	return
}

//func CheckRestApi(baseUrl string) (err error) {
//	url := baseUrl + "api.php/v1/products"
//	bytes, err := httpUtils.Get(url)
//	if err != nil {
//		return
//	}
//
//	json, err := simplejson.NewJson(bytes)
//	if err != nil {
//		return
//	}
//
//	total, err := json.Get("total").Int()
//	if err != nil {
//		return
//	}
//
//	if total < 1 {
//		err = errors.New(i118Utils.Sprintf("wrong_zentao_url"))
//	}
//
//	return
//}

func getZentaoVersion(baseUrl string) (ret string) {
	url := baseUrl + "?mode=getconfig"
	bytes, isForward, err := httpUtils.GetCheckForward(url)
	if err != nil || isForward {
		return
	}

	json, err := simplejson.NewJson(bytes)
	if err != nil {
		return
	}

	ret, _ = json.Get("version").String()

	return
}
