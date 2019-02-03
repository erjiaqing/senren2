package senrenrpc

import (
	"github.com/erjiaqing/senren2/pkg/types/base"
)

type GetContestRequest GetDomainObjectRequest

type GetContestResponse struct {
	SuccessError `json:"result"`
	Contest      *base.Contest `json:"contest"`
}

type GetContestsRequest GetDomainObjectsRequest

type GetContestsResponse struct {
	SuccessError `json:"result"`
	Contests     []*base.Contest `json:"contests"`
}

type CreateContestRequest struct {
	Session `json:"session"`
	Contest base.Contest `json:"problem"`
}

type CreateContestResponse CreateDomainObjectResponse
