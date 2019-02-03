package senrenrpc

type GetDomainObjectRequest struct {
	Session `json:"session"`
	Domain  string `json:"domain"`
	UID     string `json:"uid"`
	Filter  string `json:"filter"`
}

type GetDomainObjectsRequest struct {
	Session `json:"session"`
	Domain  string `json:"domain"`
	From    int64  `json:"from"`
	To      int64  `json:"to"`
	Filter  string `json:"filter"`
}

type CreateDomainObjectResponse struct {
	SuccessError
	Domain string `json:"domain"`
	UID    string `json:"uid"`
}
