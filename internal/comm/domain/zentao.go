package commDomain

type ZentaoUserProfile struct {
	Id       int    `json:"id"`
	Account  string `json:"account"`
	Realname string `json:"realname"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
