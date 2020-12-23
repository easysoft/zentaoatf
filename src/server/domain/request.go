package domain

type ReqData struct {
	Action string      `json:"action,omitempty"`
	Id     int         `json:"id,omitempty"`
	Mode   string      `json:"mode,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type RespData struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`

	WorkDir string `json:"workDir,omitempty"`
}
