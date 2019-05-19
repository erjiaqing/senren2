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

	// try domain user first
	qry := db.DB.QueryRow("SELECT `uid`, `guid`, `username`, `passwd` FROM `user` WHERE `domain` = ? AND `username` = ?", req.Domain, req.Username)
	if qry == nil {
		res.Success = false
		res.Error = "Failed when query database"
		return
	}
	if err := qry.Scan(&user.Uid, &user.GUid, &user.Username, &user.Password); err != nil {
		res.Success = false
		res.Error = "Domain User notfound"
	}
	if user.Password != "" && !util.CheckPass(user.Password, req.Password) {
		res.Success = false
		res.Error = "Domain Wrong password"
		// return
	}

	// then try global user
	if !res.Success {
		// Olive Auth
		authres, err := util.AuthOlive(req.Username, req.Password)
		// authres, err := true, error(nil)
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
		qry := db.DB.QueryRow("SELECT `guid`, `username`, `nickname` FROM `user` WHERE `domain` = ? AND `username` = ?", publicDomain, req.Username)
		if qry == nil {
			res.Success = false
			res.Error = "Failed when query database"
			return
		}
		if err := qry.Scan(&user.GUid, &user.Username, &user.Nickname); err != nil {
			user.Uid = util.GenUid()
			logrus.Infof("New user from olive: %s", user.Uid)
			if _, err := db.DB.Exec("INSERT INTO `user` (`uid`, `guid`, `passwd`, `username`, `nickname`, `domain`, `role`) VALUES (?, ?, ?, ?, ?, ?, ?)", user.Uid, user.Uid, "", req.Username, req.Username, publicDomain, "USER"); err != nil {
				res.Success = false
				res.Error = err.Error()
				return
			}
		}
	}

	res.Error = ""
	res.Sid = util.SignSession(user.GUid)
	res.Success = true
}

func logoutUser(ctx context.Context, req *senrenrpc.LogoutRequest, state map[string]string, res *senrenrpc.LogoutResponse) {
	// Just remove session from client XD
}

func whoami(ctx context.Context, req *senrenrpc.WhoAmIRequest, state map[string]string, res *senrenrpc.WhoAmIResponse) {
	retuser := &base.User{}
	qry := db.DB.QueryRow("SELECT `uid`, `guid`, `domain`, `username`, `nickname`, `role` FROM `user` WHERE `uid` = ?", state["uid"])

	if err := qry.Scan(&retuser.Uid, &retuser.GUid, &retuser.Domain, &retuser.Username, &retuser.Nickname, &retuser.Role); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	if retuser.Domain != state["domain"] {
		retuser.Role = "NONE"
	}

	retuser2 := &base.User{}
	qry2 := db.DB.QueryRow("SELECT `uid`, `guid`, `username`, `nickname`, `role` FROM `user` WHERE `uid` = ?", state["guid"])

	if err := qry2.Scan(&retuser2.Uid, &retuser2.GUid, &retuser2.Username, &retuser2.Nickname, &retuser2.Role); err == nil {
		res.GUser = retuser2
	}

	res.Success = true

	res.User = retuser
}
