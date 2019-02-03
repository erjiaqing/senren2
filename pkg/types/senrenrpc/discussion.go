package senrenrpc

import "github.com/erjiaqing/senren2/pkg/types/base"

type GetDiscussionRequest GetDomainObjectRequest

type GetDiscussionResponse struct {
	SuccessError
	Discussion *base.Discussion `json:"discussion"`
}

type GetDiscussionsRequest GetDomainObjectsRequest

type GetDiscussionsResponse struct {
	SuccessError
	Discussions []*base.Discussion `json:"discussions"`
}

type CreateDiscussionsRequest struct {
	Session
	Discussion base.Discussion `json:"discussion"`
}

type CreateDiscussionsResponse CreateDomainObjectResponse
