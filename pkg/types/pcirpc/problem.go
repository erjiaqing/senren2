package pcirpc

import "github.com/erjiaqing/senren2/pkg/types/base"

type CreateProblemRequest struct {
	Session
	Problem *base.PCIProblem `json:"problem"`
}

type CreateProblemResponse struct {
	SuccessError
	Uid int64 `json:"uid"`
}

type CreateProblemEditSessionRequest struct {
	Session
	ProblemAccessKey
	ProblemId
}

type CreateProblemEditSessionResponse struct {
	SuccessError
	Uid string `json:"uid"`
}

type CloseProblemEditSessionRequest struct {
	Session
	ProblemAccessKey
}

type CloseProblemEditSessionResponse SuccessErrorOnly

type GetProblemRequest struct {
	Session
	ProblemAccessKey
	ProblemId
}

type GetProblemResponse struct {
	SuccessError
	Problem *base.PCIProblem `json:"problem"`
}

type GetProblemsRequest struct {
	Session
}

type GetProblemsResponse struct {
	SuccessError
	Problems []*base.PCIProblem `json:"problems"`
}

type GetProblemDescriptionRequest struct {
	Session
	ProblemAccessKey
	ProblemId
}

type GetProblemDescriptionResponse struct {
	SuccessError
	Description string `json:"description"`
}

type GetProblemVersionsRequest GetProblemRequest

type GetProblemVersionsResponse struct {
	SuccessError
	Versions []*base.ProblemVersionState `json:"versions"`
}

type CheckProblemRequest struct {
	Session
	ProblemAccessKey
	ProblemId
	Version string `json:"version"`
}

type CheckProblemResponse SuccessErrorOnly

type GetProblemAccessKeysRequest struct {
	Session
	ProblemAccessKey
	ProblemId
}

type GetProblemAccessKeysResponse struct {
	SuccessError
	Keys []*base.PCIACL `json:"keys"`
}

type CreateProblemAccessKeyRequest struct {
	Session
	ProblemAccessKey
	ProblemId
	Permissions string `json:"perms"`
}

type CreateProblemAccessKeyResponse struct {
	SuccessError
	Key *base.PCIACL `json:"key"`
}

type ProblemAccessKey struct {
	Key string `json:"key"`
}

type HasProblemAccessKey interface {
	GetKey() string
}

func (p *ProblemAccessKey) GetKey() string {
	return p.Key
}

type ProblemId struct {
	UID int64 `json:"uid"`
}

type HasProblemId interface {
	GetId() int64
}

func (p *ProblemId) GetId() int64 {
	return p.UID
}
