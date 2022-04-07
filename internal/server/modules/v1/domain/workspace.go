package serverDomain

type WorkspaceReqPaginate struct {
	ReqPaginate

	SiteId    int `json:"siteId"`
	ProductId int `json:"productId"`
}
