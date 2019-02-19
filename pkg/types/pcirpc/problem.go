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
}

type GetProblemResponse struct {
	SuccessError
	Problem *base.PCIProblem `json:"problem"`
}

type GetProblemAccessKeyRequest struct {
	Session
	Problem int64 `json:"problem"`
}

type GetProblemAccessKeyResponse struct {
	SuccessError
	Keys []*base.PCIACL `json:"keys"`
}

type CreateProblemAccessKeyRequest struct {
	Session
	Problem     int64  `json:"problem"`
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
