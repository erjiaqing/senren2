package senrenrpc

import (
	"github.com/erjiaqing/senren2/pkg/types/base"
)

type GetDomainRequest GetDomainObjectRequest

type GetDomainResponse struct {
	SuccessError
	Domain *base.DomainInfo `json:"domain"`
}

type GetDomainsRequest GetDomainObjectsRequest

type GetDomainsResponse struct {
	SuccessError
	Domain []*base.DomainInfo `json:"domains"`
}

type CreateDomainRequest struct {
	Session
	DomainUser *base.DomainInfo `json:"domain"`
}

type CreateDomainResponse CreateDomainObjectResponse

type GetDomainInviteReuqest GetDomainObjectRequest

type GetDomainInviteResponse struct {
	SuccessError
	DomainInvite *base.DomainInvite `json:"domain_invite"`
}

type GetDomainInvitesReuqest GetDomainObjectsRequest

type GetDomainInvitesResponse struct {
	SuccessError
	DomainInvites []*base.DomainInvite `json:"domain_invites"`
}

type CreateDomainInviteRequest struct {
	Session
	DomainInvite *base.DomainInvite `json:"domain_invite"`
}

type CreateDomainInviteResponse CreateDomainObjectResponse

type JoinDomainRequest struct {
	Session
	Domain
	InviteCode string `json:"invite_code"`
}

type JoinDomainResponse SuccessErrorOnly

type UpdateDomainUserRequest struct {
	Session
	Domain
	UserUID string `json:"user_uid"`
	Status  string `json:"status"`
	Role    string `json:"role"`
}
