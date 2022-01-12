package commConsts

type ZentaoRequestType string

const (
	PathInfo ZentaoRequestType = "PATH_INFO"
	Get      ZentaoRequestType = "GET"
	Empty    ZentaoRequestType = ""
)

func (e ZentaoRequestType) String() string {
	return string(e)
}
