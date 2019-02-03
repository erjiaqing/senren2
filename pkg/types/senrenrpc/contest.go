package senrenrpc

import (
	"github.com/erjiaqing/senren2/pkg/types/base"
)

type GetContestRequest GetDomainObjectRequest

type GetContestResponse struct {
	SuccessError
	Contest *base.Contest `json:"contest"`
}

type GetContestsRequest GetDomainObjectsRequest

type GetContestsResponse struct {
	SuccessError
	Contests []*base.Contest `json:"contests"`
}

type CreateContestRequest struct {
	Session
	Contest base.Contest `json:"problem"`
}

type CreateContestResponse CreateDomainObjectResponse
