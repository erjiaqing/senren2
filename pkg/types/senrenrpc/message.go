package senrenrpc

import "github.com/erjiaqing/senren2/pkg/types/base"

type GetMessageRequest GetDomainObjectRequest

type GetMessageResponse struct {
	SuccessError `json:"result"`
	Message      *base.Message `json:"message"`
}

type GetMessagesRequest GetDomainObjectsRequest

type GetMessagesResponse struct {
	SuccessError `json:"result"`
	Messages     []*base.Message `json:"messages"`
}

type CreateMessagesRequest struct {
	Session `json:"session"`
	Message base.Message `json:"message"`
}

type CreateMessagesResponse CreateDomainObjectResponse

type MarkMessageReadRequest struct {
	Session `json:"session"`
	Message base.Message `json:"message"`
}

type MarkMessageReadResponse SuccessErrorOnly
