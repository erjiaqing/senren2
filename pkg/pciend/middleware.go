package pciend

import (
	"context"
	"fmt"
	"strings"

	"github.com/erjiaqing/senren2/pkg/pcidb"
	"github.com/erjiaqing/senren2/pkg/types/pcirpc"
	"github.com/sirupsen/logrus"
)

func resolveProblemAccessKey(ctx context.Context, k pcirpc.HasProblemAccessKey, state map[string]string) {
	prob := int64(0)
	perm := ""
	pkey := ""

	row := pcidb.PCIDB.QueryRowContext(ctx, "SELECT aclpkey, puid, perm FROM acl WHERE aclkey = ?", k.GetKey())

	if err := row.Scan(&pkey, &prob, &perm); err != nil {
		return
	}

	state["PKEY"] = pkey
	state["PROB"] = fmt.Sprintf("%d", prob)

	perms := strings.Split(perm, "|")
	for _, v := range perms {
		state["PERM_"+v] = "G"
		logrus.Infof("Granted Permission: %s", v)
	}
}
