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

type GetContestProblemRequest GetDomainObjectRequest

type GetContestProblemResponse GetProblemResponse

type CreateContestSubmissionRequest = CreateSubmissionsRequest

type CreateContestSubmissionResponse = CreateSubmissionsResponse

type GetContestSubmissionsRequest = GetSubmissionsRequest

type GetContestSubmissionsResponse = GetSubmissionsResponse

type GetContestSubmissionRequest = GetSubmissionRequest

type GetContestSubmissionResponse = GetSubmissionResponse

type GetContestsResponse struct {
	SuccessError
	Contests []*base.Contest `json:"contests"`
}

type CreateContestRequest struct {
	Session
	Contest base.Contest `json:"contest"`
}

type CreateContestResponse CreateDomainObjectResponse
