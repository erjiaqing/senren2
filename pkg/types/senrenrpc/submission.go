package senrenrpc

import (
	"github.com/erjiaqing/senren2/pkg/types/base"
)

type GetSubmissionRequest GetDomainObjectRequest

type GetSubmissionResponse struct {
	SuccessError
	Submission *base.Submission `json:"submission"`
}

type GetSubmissionsRequest GetDomainObjectsRequest

type GetSubmissionsResponse struct {
	SuccessError
	Submissions []*base.Submission `json:"submissions"`
}

type CreateSubmissionsRequest struct {
	Session
	Submission base.Submission `json:"submission"`
}

type CreateSubmissionsResponse CreateDomainObjectResponse
