package senrenrpc

import "github.com/erjiaqing/senren2/pkg/types/base"

type GetDiscussionRequest GetDomainObjectRequest

type GetDiscussionResponse struct {
	SuccessError `json:"result"`
	Discussion   *base.Discussion `json:"discussion"`
}

type GetDiscussionsRequest GetDomainObjectsRequest

type GetDiscussionsResponse struct {
	SuccessError `json:"result"`
	Discussions  []*base.Discussion `json:"discussions"`
}

type CreateDiscussionsRequest struct {
	Session    `json:"session"`
	Discussion base.Discussion `json:"discussion"`
}

type CreateDiscussionsResponse CreateDomainObjectResponse
