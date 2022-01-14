package _rpcUtils

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
)

func Post(url, method, api, mtd string, args interface{}) (resp domain.RpcResp) {
	//cc := &codec.MsgpackCodec{}
	//
	//data, _ := cc.Encode(args)
	//
	//req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewReader(data))
	//if err != nil {
	//	logUtils.Errorf("failed to create request: ", err)
	//	return
	//}
	//
	//h := req.Header
	//h.Set(gateway.XMessageID, "10000")
	//h.Set(gateway.XMessageType, "0")
	//h.Set(gateway.XSerializeType, "3")
	//h.Set(gateway.XServicePath, api)
	//h.Set(gateway.XServiceMethod, mtd)
	//
	//res, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	log.Errorf("failed to call: ", err)
	//}
	//defer res.Body.Close()
	//
	//replyData, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	log.Errorf("failed to read response: ", err)
	//}
	//
	//err = cc.Decode(replyData, &resp)
	//if err != nil {
	//	log.Errorf("failed to decode reply: ", err)
	//}
	//log.Infof("%v -> %v", args, resp)

	return
}
