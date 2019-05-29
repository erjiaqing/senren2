package senrenrpc

import (
	"github.com/erjiaqing/senren2/pkg/types/base"
)

type GetProblemRequest GetDomainObjectRequest

type GetProblemResponse struct {
	SuccessError
	Problem *base.Problem `json:"problem"`
}

type GetProblemsRequest GetDomainObjectsPager

type GetProblemsResponse struct {
	PagerSuccessError
	Problems []*base.Problem `json:"problems"`
}

type GetPCIDescriptionRequest GetDomainObjectRequest

type GetPCIDescriptionResponse struct {
	SuccessError
	Description string `json:"description"`
}

type CreateProblemRequest struct {
	Session
	Problem base.Problem `json:"problem"`
}

type CreateProblemResponse CreateDomainObjectResponse
