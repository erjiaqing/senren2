package pciend

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/erjiaqing/senren2/pkg/httpreq"
	"github.com/erjiaqing/senren2/pkg/util"

	"github.com/erjiaqing/senren2/pkg/pcidb"
	"github.com/erjiaqing/senren2/pkg/repo"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/pcirpc"
)

var createEditSessionMutex sync.Mutex
var editorServer string

func init() {
	editorServer = os.Getenv("PCI_EDITOR_SERV")
	if editorServer == "" {
		editorServer = "http://localhost:8078"
	}
}

func getProblem(ctx context.Context, req *pcirpc.GetProblemRequest, state map[string]string, res *pcirpc.GetProblemResponse) {

}

func createProblem(ctx context.Context, req *pcirpc.CreateProblemRequest, state map[string]string, res *pcirpc.CreateProblemResponse) {
	tx, err := pcidb.PCIDB.BeginTx(ctx, nil)
	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	if req.Problem.Uid == -1 {
		// do create
		qry, err := tx.ExecContext(ctx, "INSERT INTO problem (title, remoteURL, currentVersion, editSession, owner, state) VALUES (?, ?, ?, ?, ?, ?)", req.Problem.Title, ".", "", "-", state["login"], "NEW")
		if err != nil {
			res.Success = false
			res.Error = err.Error()
		}

		req.Problem.Owner = state["login"]
		req.Problem.CurrentVersion = ""
		req.Problem.Uid, _ = qry.LastInsertId()
		req.Problem, err = repo.CreateProblemRepo(req.Problem)

		if err != nil {
			res.Success = false
			res.Error = err.Error()
			return
		}

		if _, err := tx.ExecContext(ctx, "UPDATE problem SET remoteURL = ? WHERE uid = ?", req.Problem.RemoteURL, req.Problem.Uid); err != nil {
			res.Success = false
			res.Error = err.Error()
			return
		}
	}

	// TODO: check currentVersion is valid (exists and built)
	if _, err := tx.ExecContext(ctx, "UPDATE problem SET title = ?, currentVersion = ? WHERE uid = ? AND owner = ?", req.Problem.Title, req.Problem.CurrentVersion, req.Problem.Uid, state["login"]); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	if err := tx.Commit(); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Success = true
	res.Uid = req.Problem.Uid
}

func createProblemEditSession(ctx context.Context, req *pcirpc.CreateProblemEditSessionRequest, state map[string]string, res *pcirpc.CreateProblemEditSessionResponse) {
	// call problem editor to clone and create problem edit session
	createEditSessionMutex.Lock()
	defer createEditSessionMutex.Unlock()

	tprobRow := pcidb.PCIDB.QueryRowContext(ctx, "SELECT remoteURL FROM problem WHERE uid = ?", state["PROB"])
	probURL := ""

	if err := tprobRow.Scan(&probURL); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	probURL = repo.GetCloneAddress(probURL)

	if probURL == "" {
		res.Success = false
		res.Error = "Failed to get problem repo url"
		return
	}

	dat, code, err := httpreq.POSTJson(editorServer+"/api/repo/pull", map[string]string{
		"repo": probURL,
	})

	if code >= 300 {
		res.Success = false
		res.Error = "500 - " + string(dat)
		return
	}

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	uid := string(dat)

	if uid == "" {
		res.Success = false
		res.Error = "failed to clone"
		return
	}

	pcidb.PCIDB.ExecContext(ctx, "UPDATE problem SET editSession = ? WHERE uid = ?", uid, state["PROB"])
	res.Success = true
	res.Uid = uid
}

func closeProblemEditSession(ctx context.Context, req *pcirpc.CloseProblemEditSessionRequest, state map[string]string, res *pcirpc.CloseProblemEditSessionResponse) {
	// call problem editor to ignore uncommited edit and delete problem edit session
}

func createProblemAccessKey(ctx context.Context, req *pcirpc.CreateProblemAccessKeyRequest, state map[string]string, res *pcirpc.CreateProblemAccessKeyResponse) {
	pubKey, err := util.GenerateRandomString(24)
	pubKey = fmt.Sprintf("%08x", time.Now().Unix()) + pubKey
	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	privKey, err := util.GenerateRandomString(32)
	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	_, err = pcidb.PCIDB.ExecContext(ctx, "INSERT INTO acl (aclkey, aclpkey, puid, perm, create_at) VALUES (?, ?, ?, ?, ?)", pubKey, privKey, req.Problem, req.Permissions, time.Now())
	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Success = true
	res.Key = &base.PCIACL{
		Key:           pubKey,
		PrivateKey:    privKey,
		ProblemUID:    req.Problem,
		AccessControl: req.Permissions,
	}
}
