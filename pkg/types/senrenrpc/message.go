package senrenrpc

import "github.com/erjiaqing/senren2/pkg/types/base"

type GetMessageRequest GetDomainObjectRequest

type GetMessageResponse struct {
	SuccessError
	Message *base.Message `json:"message"`
}

type GetMessagesRequest GetDomainObjectsRequest

type GetMessagesResponse struct {
	SuccessError
	Messages []*base.Message `json:"messages"`
}

type CreateMessagesRequest struct {
	Session
	Message base.Message `json:"message"`
}

type CreateMessagesResponse CreateDomainObjectResponse

type MarkMessageReadRequest struct {
	Session
	Message base.Message `json:"message"`
}

type MarkMessageReadResponse SuccessErrorOnly
