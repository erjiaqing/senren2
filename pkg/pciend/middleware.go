package pciend

import (
	"context"
	"fmt"
	"strings"

	"github.com/erjiaqing/senren2/pkg/pcidb"
	"github.com/erjiaqing/senren2/pkg/types/pcirpc"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/erjiaqing/senren2/pkg/util"
	"github.com/sirupsen/logrus"
)

func checkLogin(ctx context.Context, req senrenrpc.HasSession, state map[string]string) {
	session := util.CheckSession(req.GetSession())
	state["USER"] = session
}

func resolveProblemAccessKey(ctx context.Context, k pcirpc.HasProblemAccessKey, uid int64, state map[string]string) {
	prob := int64(0)
	perm := ""
	pkey := ""
	powner := ""

	row := pcidb.PCIDB.QueryRowContext(ctx, "SELECT aclpkey, puid, perm FROM acl WHERE aclkey = ?", k.GetKey())

	if err := row.Scan(&pkey, &prob, &perm); err != nil {
		row2 := pcidb.PCIDB.QueryRowContext(ctx, "SELECT owner FROM problem WHERE uid = ?", uid)
		row2.Scan(&powner)
		if powner == "" || powner != state["USER"] {
			return
		}
	}

	state["PKEY"] = pkey
	state["PROB"] = fmt.Sprintf("%d", prob)

	perms := strings.Split(perm, "|")
	for _, v := range perms {
		state["PERM_"+v] = "G"
		logrus.Infof("Granted Permission: %s", v)
	}
}
