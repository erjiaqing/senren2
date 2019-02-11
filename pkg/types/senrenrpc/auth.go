package senrenrpc

import "github.com/erjiaqing/senren2/pkg/types/base"

type Session struct {
	Sid string `json:"sid"`
}

type HasSession interface {
	GetSession() string
}

func (s *Session) GetSession() string {
	return s.Sid
}

type AuthRequest struct {
	Domain
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Session
	SuccessError
}

type LogoutRequest struct {
	Session
}

type LogoutResponse SuccessErrorOnly

type CheckInviteRequest struct {
	Domain
	InviteCode string `json:"invite_code"`
	InvitePass string `json:"invite_pass"`
}

type ConfirmInviteRequest CheckInviteRequest

type CheckInviteResponse struct {
	SuccessError
	Invite base.DomainInvite `json:"invite"`
}

type ConfirmInviteResponse SuccessErrorOnly

type WhoAmIRequest struct {
	Session
	Domain
}

type WhoAmIResponse struct {
	SuccessError
	User  *base.User `json:"user"`
	GUser *base.User `json:"user_global"`
}
