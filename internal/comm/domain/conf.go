package commDomain

type ProjectConf struct {
	Language string `json:"language"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`

	Javascript string `json:"javascript"`
	Lua        string `json:"lua"`
	Perl       string `json:"perl"`
	Php        string `json:"php"`
	Python     string `json:"python"`
	Ruby       string `json:"ruby"`
	Tcl        string `json:"tcl"`
	Autoit     string `json:"autoit"`

	Version string `json:"version"`
	IsWin   bool   `json:"isWin"`
}
