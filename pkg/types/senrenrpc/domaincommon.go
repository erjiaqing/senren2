package senrenrpc

type GetDomainObjectRequest struct {
	Session
	Domain `json:"domain"`
	UID    string `json:"uid"`
	Filter string `json:"filter"`
}

type GetDomainObjectsRequest struct {
	Session
	Domain `json:"domain"`
	From   int64  `json:"from"`
	To     int64  `json:"to"`
	Filter string `json:"filter"`
}

type GetDomainObjectsPager struct {
	Session
	Domain `json:"domain"`
	Page   int    `json:"page"`
	Count  int    `json:"count"`
	Filter string `json:"filter"`
}

type CreateDomainObjectResponse struct {
	SuccessError
	Domain `json:"domain"`
	UID    string `json:"uid"`
}
