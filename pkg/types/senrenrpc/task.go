package senrenrpc

import "github.com/erjiaqing/senren2/pkg/types/base"

type GetTaskRequest GetDomainObjectRequest

type GetTaskResponse struct {
	SuccessError
	Task *base.ServerTask `json:"task"`
}
