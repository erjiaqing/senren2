package endpoint

import (
	"context"
	"strings"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/erjiaqing/senren2/pkg/util"
)

func getContest(ctx context.Context, req *senrenrpc.GetContestRequest, state map[string]string, res *senrenrpc.GetContestResponse) {
	res.Contest = &base.Contest{}
	row := db.DB.QueryRowContext(ctx, "SELECT uid, domain, title, content, type, problem_list, start_time, end_time, open_time, close_time, freeze_time, release_time FROM contest WHERE `uid` = ? AND `domain` = ?", req.UID, req.Domain)
	if err := row.Scan(&res.Contest.Uid, &res.Contest.Domain, &res.Contest.Title, &res.Contest.Type, &res.Contest.ProblemList, &res.Contest.StartTime, &res.Contest.EndTime, &res.Contest.FreezeTime, &res.Contest.ReleaseTime); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	res.Success = true
}

func getContests(ctx context.Context, req *senrenrpc.GetContestsRequest, state map[string]string, res *senrenrpc.GetContestsResponse) {
	// TODO: filter
	filters := strings.Split(req.Filter, ",")
	res.Contests = make([]*base.Contest, 0)
	rows, err := db.DB.QueryContext(ctx, "SELECT uid, domain, title, content, type, problem_list, start_time, end_time, open_time, close_time, freeze_time, release_time FROM contest WHERE `uid` = ? AND `domain` = ?", filters[0], req.Domain)
	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Success = true

	for rows.Next() {
		tContest := &base.Contest{}
		if err := rows.Scan(&tContest.Uid, &tContest.Domain, &tContest.Title, &tContest.Type, &tContest.StartTime, &tContest.EndTime, &tContest.FreezeTime, &tContest.ReleaseTime); err != nil {
			res.Success = false
			res.Error = err.Error()
			rows.Close()
			break
		}
		res.Contests = append(res.Contests, tContest)
	}
}

// also updates contest
func createContest(ctx context.Context, req *senrenrpc.CreateContestRequest, state map[string]string, res *senrenrpc.CreateContestResponse) {
	dbExec := "UPDATE contest SET title = ?, content = ?, type = ?, problem_list = ?, start_time = ?, end_time = ?, open_time = ?, close_time = ?, freeze_time = ?, release_time = ? WHERE `uid` = ? AND `domain` = ?"

	if req.Contest.Uid == "" || req.Contest.Uid == noUID {
		// No uid, create a contest
		req.Contest.Uid = util.GenUid()
		dbExec = "INSERT INTO contest (title, content, type, problem_list, start_time, end_time, open_time, close_time, freeze_time, release_time, uid, domain) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	}

	_, err := db.DB.ExecContext(ctx, dbExec,
		req.Contest.Title, req.Contest.Description, req.Contest.Type, req.Contest.ProblemList, req.Contest.StartTime, req.Contest.EndTime, req.Contest.OpenTime, req.Contest.CloseTime, req.Contest.FreezeTime, req.Contest.ReleaseTime, req.Contest.Uid, req.Contest.Domain)
	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	res.Success = true
	res.Domain = req.Contest.Domain
	res.UID = req.Contest.Uid
}
