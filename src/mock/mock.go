package mock

var (
	BaseUrl       = "http://localhost:8888/"
	ImportProject = "importProject"
	ReportBugs    = "reportBugs"
)

func GetUrl(uri string) string {
	return BaseUrl + uri
}
