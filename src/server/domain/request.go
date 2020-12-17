package domain

type ReqData struct {
	Action string `json:"action"`
	Id     int    `json:"id"`
	Mode   string `json:"mode"`
	Data   string `json:"data"`
}

type RespData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`

	WorkDir string `json:"workDir"`
}
