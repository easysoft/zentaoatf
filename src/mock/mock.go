package mock

var (
	BaseUrl = "http://localhost:8888/"
)

func GetUrl(uri string) string {
	return BaseUrl + uri
}
