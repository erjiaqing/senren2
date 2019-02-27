package pciend

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

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

func problemUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pid, _ := strconv.Atoi(params["problem"])
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	versions := repo.ParseCommitPush(int64(pid), reqBody)

	stat, err := pcidb.PCIDB.PrepareContext(r.Context(), "INSERT INTO problemVersion (p_uid, version, state, logtime, message) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		w.WriteHeader(500)
		return
	}

	commitedVersions := 0

	for _, v := range versions {
		if _, err := stat.ExecContext(r.Context(), pid, v.Version, v.State, v.LogTime, v.Message); err == nil {
			commitedVersions++
		} else {
			logrus.Errorf("Failed to do SQL: %v", err)
		}
	}
	w.Write([]byte(fmt.Sprintf("success, %d versions logged", commitedVersions)))
}

func getProblem(ctx context.Context, req *pcirpc.GetProblemRequest, state map[string]string, res *pcirpc.GetProblemResponse) {
	if state["PROB"] == "-1" {
		state["PROB"] = fmt.Sprintf("%d", req.ProblemId)
	}
	row := pcidb.PCIDB.QueryRowContext(ctx, "SELECT uid, title, remoteURL, currentVersion, owner, state FROM problem WHERE uid = ?", state["PROB"])
	ret := &base.PCIProblem{}
	if err := row.Scan(&ret.Uid, &ret.Title, &ret.RemoteURL, &ret.CurrentVersion, &ret.Owner, &ret.State); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	res.Success = true
	res.Problem = ret
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

	tprobRow := pcidb.PCIDB.QueryRowContext(ctx, "SELECT remoteURL, editSession FROM problem WHERE uid = ?", state["PROB"])
	probURL, currSession := "", ""

	if err := tprobRow.Scan(&probURL, &currSession); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	if currSession != "" && currSession != "-" {
		res.Success = true
		res.Uid = currSession
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

	uidDict := make(map[string]interface{})
	if err := json.Unmarshal(dat, &uidDict); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	uid, ok := uidDict["text"].(string)

	if !ok || uid == "" {
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

func getProblemDescription(ctx context.Context, req *pcirpc.GetProblemDescriptionRequest, state map[string]string, res *pcirpc.GetProblemDescriptionResponse) {
	row := pcidb.PCIDB.QueryRowContext(ctx, "SELECT currentVersion, remoteURL FROM problem WHERE uid = ?", state["PROB"])
	var rev, remote string
	if err := row.Scan(&rev, &remote); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	ret := repo.FetchProblemDescription(remote, rev)
	res.Success = (ret != nil)
	res.Description = string(ret)
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

func getProblemVersions(ctx context.Context, req *pcirpc.GetProblemVersionsRequest, state map[string]string, res *pcirpc.GetProblemVersionsResponse) {
	row, err := pcidb.PCIDB.QueryContext(ctx, "SELECT version, state, logtime, message FROM problemVersion WHERE p_uid = ? ORDER BY logtime DESC", state["PROB"])
	if err != nil {
		res.Success = false
		res.Error = err.Error()
	}

	res.Versions = make([]*base.ProblemVersionState, 0)
	for row.Next() {
		t := &base.ProblemVersionState{}
		if err := row.Scan(&t.Version, &t.State, &t.LogTime, &t.Message); err != nil {
			res.Success = false
			return
		}
		res.Versions = append(res.Versions, t)
	}
	res.Success = true
}
