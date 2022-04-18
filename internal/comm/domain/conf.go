package commDomain

type WorkspaceConf struct {
	Language string `json:"language"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`

	Javascript string `ini:"omitempty" json:"javascript"`
	Lua        string `ini:"omitempty" json:"lua"`
	Perl       string `ini:"omitempty" json:"perl"`
	Php        string `ini:"omitempty" json:"php"`
	Python     string `ini:"omitempty" json:"python"`
	Ruby       string `ini:"omitempty" json:"ruby"`
	Tcl        string `ini:"omitempty" json:"tcl"`
	Autoit     string `ini:"omitempty" json:"autoit"`

	Version string `json:"version"`
	IsWin   bool   `ini:"-" json:"isWin"`
}
