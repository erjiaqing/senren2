package senrenrpc

import (
	"github.com/erjiaqing/senren2/pkg/types/base"
)

type GetProblemRequest GetDomainObjectRequest

type GetProblemResponse struct {
	SuccessError `json:"result"`
	Problem      *base.Problem `json:"problem"`
}

type GetProblemsRequest GetDomainObjectsRequest

type GetProblemsResponse struct {
	SuccessError `json:"result"`
	Problems     []*base.Problem `json:"problems"`
}

type CreateProblemRequest struct {
	Session `json:"session"`
	Problem base.Problem `json:"problem"`
}

type CreateProblemResponse CreateDomainObjectResponse
