package repo

import "time"

type user struct {
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	ID        int64  `json:"id"`
	UserName  string `json:"login"`
}

type repo struct {
	CloneURL      string    `json:"clone_url"`
	CreatedAt     time.Time `json:"created_at"`
	DefaultBranch string    `json:"default_branch"`
	Description   string    `json:"description"`
	Empty         bool      `json:"empty"`
	Fork          bool      `json:"fork"`
	Forks         int64     `json:"forks_count"`
	FullName      string    `json:"full_name"`
	HTMLURL       string    `json:"html_url"`
	ID            int64     `json:"id"`
	Mirror        bool      `json:"mirror"`
	Name          string    `json:"name"`
	Owner         *user     `json:"owner"`
}

type payloadUser struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	UserName string `json:"username"`
}

type payloadCommit struct {
	Author       *payloadUser         `json:"author"`
	Committer    *payloadUser         `json:"committer"`
	ID           string               `json:"id"`
	Message      string               `json:"message"`
	TimeStamp    time.Time            `json:"timestamp"`
	URL          string               `json:"url"`
	Verification *payloadVerification `json:"verification"`
}

type payloadVerification struct {
	Payload   string `json:"payload"`
	Reason    string `json:"reason"`
	Signature string `json:"signature"`
	Verified  bool   `json:"verified"`
}

type branch struct {
	Commit *payloadCommit `json:"commit"`
	Name   string         `json:"name"`
}
