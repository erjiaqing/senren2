package endpoint

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/erjiaqing/senren2/pkg/util"
)

func checkLogin(ctx context.Context, req senrenrpc.HasSession, domain string, state map[string]string) {
	session := util.CheckSession(req.GetSession())

	row := db.DB.QueryRowContext(ctx, "SELECT uid, username, role FROM user WHERE guid = ? AND domain = ?", session, domain)
	row2 := db.DB.QueryRowContext(ctx, "SELECT uid, username, role FROM user WHERE guid = ? AND domain = ?", session, "0000000000000000")

	if row == nil {
		return
	}

	var (
		uid   string
		role  string
		uname string

		guid   string
		grole  string
		guname string
	)

	state["global_login"] = session

	logrus.Debugf("Global Login: %s @ domain %s", session, domain)

	if err := row.Scan(&uid, &uname, &role); err != nil {
		role = "NONE"
		uid = session
	}

	if err := row2.Scan(&guid, &guname, &grole); err != nil {
		grole = "NONE"
	}

	state["role"] = role
	state["uid"] = uid

	state["guid"] = guid
	state["gname"] = guname
	state["grole"] = grole

	logrus.Debugf("Role: %s @ Uid %s , GUid %s", role, uid, guid)
}
