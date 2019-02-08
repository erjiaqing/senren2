package endpoint

import (
	"context"
	"time"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/erjiaqing/senren2/pkg/util"
)

func getProblem(ctx context.Context, req *senrenrpc.GetProblemRequest, state map[string]string, res *senrenrpc.GetProblemResponse) {
	r := &base.Problem{}
	queryString := "SELECT uid, rootuid, domain, alias, title, content, releasetime, problemci, score, language_limit FROM problem WHERE uid = ? AND domain = ?"
	if len(req.UID) != 16 && len(req.UID) <= 8 {
		queryString = "SELECT uid, rootuid, domain, alias, title, content, releasetime, problemci, score, language_limit FROM problem WHERE alias = ? AND domain = ?"
	} else if len(req.UID) != 16 {
		res.Success = false
		res.Error = "illegal problem uid or alias name"
		return
	}
	if req.UID == noUID {
		r.ReleaseTime = time.Now()
		res.Problem = r
		res.Success = true
		return
	}
	row := db.DB.QueryRow(queryString, req.UID, req.Domain)
	if err := row.Scan(&r.Uid, &r.RootUid, &r.Domain, &r.Alias, &r.Title, &r.Description, &r.ReleaseTime, &r.ProblemCI, &r.Score, &r.LanguageLimit); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	res.Problem = r
	res.Success = true
}

func getProblems(ctx context.Context, req *senrenrpc.GetProblemsRequest, state map[string]string, res *senrenrpc.GetProblemsResponse) {
	row, err := db.DB.Query("SELECT uid, rootuid, domain, alias, title, score, releasetime FROM problem WHERE domain = ?", req.Domain)
	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	ret := make([]*base.Problem, 0)

	for row.Next() {
		r := &base.Problem{}
		if err := row.Scan(&r.Uid, &r.RootUid, &r.Domain, &r.Alias, &r.Title, &r.Score, &r.ReleaseTime); err != nil {
			row.Close()
			res.Success = false
			res.Error = err.Error()
			return
		}
		ret = append(ret, r)
	}

	res.Problems = ret
	res.Success = true
}

func createProblem(ctx context.Context, req *senrenrpc.CreateProblemRequest, state map[string]string, res *senrenrpc.CreateProblemResponse) {
	dbExec := "UPDATE problem SET title = ? , content = ? , releasetime = ?, problemci = ?, score = ?, language_limit = ?, alias = ? WHERE uid = ? AND (rootuid = ? OR 1 = 1) AND domain = ?"

	tDomain := senrenrpc.Domain(req.Problem.Domain)
	tDomain.Convert()
	req.Problem.Domain = string(tDomain)

	if req.Problem.Uid == "" || req.Problem.Uid == noUID {
		req.Problem.Uid = util.GenUid()
		if req.Problem.RootUid == "" {
			req.Problem.RootUid = req.Problem.Uid
		}

		dbExec = "INSERT INTO problem (title, content, releasetime, problemci, score, language_limit, alias, uid, rootuid, domain) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	}

	if _, err := db.DB.Exec(dbExec, req.Problem.Title, req.Problem.Description, req.Problem.ReleaseTime, req.Problem.ProblemCI, req.Problem.Score, req.Problem.LanguageLimit, req.Problem.Alias, req.Problem.Uid, req.Problem.RootUid, req.Problem.Domain); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Domain = senrenrpc.Domain(req.Problem.Domain)
	res.UID = req.Problem.Uid
	res.Success = true
}
