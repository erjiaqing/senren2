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
	Domains []*base.DomainInfo `json:"domains"`
}

type CreateDomainRequest struct {
	Session
	Domain *base.DomainInfo `json:"domain"`
}

type CreateDomainResponse CreateDomainObjectResponse

type GetDomainInviteRequest GetDomainObjectRequest

type GetDomainInviteResponse struct {
	SuccessError
	DomainInvite *base.DomainInvite `json:"domain_invite"`
}

type GetDomainInvitesRequest GetDomainObjectsRequest

type GetDomainInvitesResponse struct {
	SuccessError
	DomainInvites []*base.DomainInvite `json:"domain_invites"`
}

type CreateDomainInviteRequest struct {
	Session
	Domain       `json:"domain"`
	DomainInvite *base.DomainInvite `json:"domain_invite"`
}

type CreateDomainInviteResponse CreateDomainObjectResponse

type JoinDomainRequest struct {
	Session
	Domain         `json:"domain"`
	NickName       string `json:"nickname"`
	InviteCode     string `json:"invite_code"`
	InvitePassword string `json:"invite_password"`
}

type JoinDomainResponse SuccessErrorOnly

type UpdateDomainUserRequest struct {
	Session
	Domain  `json:"domain"`
	UserUID string `json:"user_uid"`
	Status  string `json:"status"`
	Role    string `json:"role"`
}

type UpdateDomainUserResponse SuccessErrorOnly

type GetDomainUserRequest GetDomainObjectRequest

type GetDomainUserResponse struct {
	SuccessError
	User *base.User `json:"user"`
}

type GetDomainUsersRequest GetDomainObjectsRequest

type GetDomainUsersResponse struct {
	SuccessError
	Users []*base.User `json:"users"`
}

type GetPCISidRequest GetDomainObjectRequest

type GetPCISidResponse AuthResponse
