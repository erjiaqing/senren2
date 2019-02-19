package pcirpc

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
