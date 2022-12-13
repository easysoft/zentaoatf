package serverDomain

type HostHeartbeatReq struct {
	MacAddress string `json:"macAddress"`
}
type HostHeartbeatResp struct {
	Token  string `json:"token" yaml:"token"`
	Server string `json:"server" yaml:"server"`
}
type HostResponse struct {
	Code string      `json:"code"` // Enums consts.ResultCode
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
