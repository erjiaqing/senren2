package senrenrpc

import (
	"github.com/erjiaqing/senren2/pkg/types/base"
)

type GetProblemRequest GetDomainObjectRequest

type GetProblemResponse struct {
	SuccessError
	Problem *base.Problem `json:"problem"`
}

type GetProblemsRequest GetDomainObjectsRequest

type GetProblemsResponse struct {
	SuccessError
	Problems []*base.Problem `json:"problems"`
}

type CreateProblemRequest struct {
	Session
	Problem base.Problem `json:"problem"`
}

type CreateProblemResponse CreateDomainObjectResponse
