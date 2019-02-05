package endpoint

import (
	"context"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/erjiaqing/senren2/pkg/util"
	"github.com/sirupsen/logrus"
)

func authUser(ctx context.Context, req *senrenrpc.AuthRequest, state map[string]string, res *senrenrpc.AuthResponse) {
	user := &base.User{}
	if req.Domain == "" {
		// Olive Auth
		authres, err := util.AuthOlive(req.Username, req.Password)
		if err != nil {
			res.Success = false
			res.Error = err.Error()
			return
		}
		if !authres {
			res.Success = false
			res.Error = "Wrong Password"
			return
		}
		qry := db.DB.QueryRow("SELECT `uid`, `username`, `nickname` FROM `user` WHERE `domain` = ? AND `username` = ?", publicDomain, req.Username)
		if qry == nil {
			res.Success = false
			res.Error = "Failed when query database"
			return
		}
		if err := qry.Scan(&user.Uid, &user.Username, &user.Nickname); err != nil {
			user.Uid = util.GenUid()
			logrus.Infof("New user from olive: %s", user.Uid)
			if _, err := db.DB.Exec("INSERT INTO `user` (`uid`, `username`, `nickname`, `domain`) VALUES (?, ?, ?, ?)", user.Uid, req.Username, req.Username, publicDomain); err != nil {
				res.Success = false
				res.Error = err.Error()
				return
			}
		}
	} else {
		qry := db.DB.QueryRow("SELECT `uid`, `username`, `passwd` FROM `user` WHERE `domain` = ? AND `username` = ?", req.Domain, req.Username)
		if qry == nil {
			res.Success = false
			res.Error = "Failed when query database"
			return
		}
		if err := qry.Scan(&user.Uid, &user.Username, &user.Password); err != nil {
			res.Success = false
			res.Error = "User notfound"
			return
		}
		if !util.CheckPass(user.Password, req.Password) {
			res.Success = false
			res.Error = "Wrong password"
			return
		}
	}
	res.Sid = util.SignSession(user.Uid)
	res.Success = true
}

func logoutUser(ctx context.Context, req *senrenrpc.LogoutRequest, state map[string]string, res *senrenrpc.LogoutResponse) {
	// Just remove session from client XD
}
