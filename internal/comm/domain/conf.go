package commDomain

type ProjectConf struct {
	Language string `json:"language"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`

	Javascript string `ini:"omitthis,omitempty" json:"javascript"`
	Lua        string `ini:"omitthis,omitempty" json:"lua"`
	Perl       string `ini:"omitthis,omitempty" json:"perl"`
	Php        string `ini:"omitthis,omitempty" json:"php"`
	Python     string `ini:"omitthis,omitempty" json:"python"`
	Ruby       string `ini:"omitthis,omitempty" json:"ruby"`
	Tcl        string `ini:"omitthis,omitempty" json:"tcl"`
	Autoit     string `ini:"omitthis,omitempty" json:"autoit"`

	Version string `json:"version"`
	IsWin   bool   `ini:"-" json:"isWin"`
}
