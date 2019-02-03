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
	Domain   string `json:"domain"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Session      `json:"session"`
	SuccessError `json:"result"`
}

type LogoutRequest struct {
	Session `json:"session"`
}

type LogoutResponse SuccessErrorOnly

type CheckInviteRequest struct {
	Domain     string `json:"domain"`
	InviteCode string `json:"invite_code"`
	InvitePass string `json:"invite_pass"`
}

type ConfirmInviteRequest CheckInviteRequest

type CheckInviteResponse struct {
	SuccessError `json:"result"`
	Invite       base.DomainInvite `json:"invite"`
}

type ConfirmInviteResponse SuccessErrorOnly
