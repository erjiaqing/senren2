package pciend

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/erjiaqing/senren2/pkg/pcidb"
	"github.com/erjiaqing/senren2/pkg/types/pcirpc"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/erjiaqing/senren2/pkg/util"
	"github.com/sirupsen/logrus"
)

func checkLogin(ctx context.Context, req senrenrpc.HasSession, state map[string]string) {
	session := util.CheckSessionTimeDomain(req.GetSession(), 365*24*time.Hour, "0000000000000001")
	state["USER"] = session
}

func resolveProblemAccessKey(ctx context.Context, k pcirpc.HasProblemAccessKey, uid int64, state map[string]string) {
	prob := int64(0)
	perm := ""
	pkey := ""
	powner := ""

	row2 := pcidb.PCIDB.QueryRowContext(ctx, "SELECT owner FROM problem WHERE uid = ?", uid)
	row2.Scan(&powner)
	logrus.Debugf("Problem: %d Problem Owner: %s Login User: %s", uid, powner, state["USER"])
	if powner != "" && powner == state["USER"] {
		logrus.Debug("Granted Permission: __owner__")
		state["PERM_."] = "G"
		state["PKEY"] = pkey
		state["PROB"] = fmt.Sprintf("%d", uid)
	}

	row := pcidb.PCIDB.QueryRowContext(ctx, "SELECT aclpkey, puid, perm FROM acl WHERE aclkey = ?", k.GetKey())

	if err := row.Scan(&pkey, &prob, &perm); err != nil {
		return
	}

	if prob == -1 && uid > 0 {
		state["PROB"] = fmt.Sprintf("%d", uid)
	} else {
		state["PROB"] = fmt.Sprintf("%d", prob)
	}

	logrus.Debugf("Problem: %s (%d, %d) Problem Owner: %s Login User: %s", state["PROB"], prob, uid, powner, state["USER"])
	perms := strings.Split(perm, "|")
	for _, v := range perms {
		state["PERM_"+v] = "G"
		logrus.Debugf("Granted Permission: %s", v)
	}
}
