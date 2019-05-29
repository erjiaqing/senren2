package endpoint

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/erjiaqing/senren2/pkg/util"
)

func getContest(ctx context.Context, req *senrenrpc.GetContestRequest, state map[string]string, res *senrenrpc.GetContestResponse) {
	res.Contest = &base.Contest{}
	res.Success = true

	if req.UID == "" || req.UID == noUID {
		current := time.Now()
		next15Minute := (current.Unix()/(15*60) + 1) * (15 * 60)
		res.Contest.OpenTime = current
		res.Contest.CloseTime = current.AddDate(0, 0, 7)
		res.Contest.StartTime = time.Unix(next15Minute, 0)
		res.Contest.EndTime = time.Unix(next15Minute+300*60, 0)
		res.Contest.FreezeTime = time.Unix(next15Minute+240*60, 0)
		res.Contest.ReleaseTime = res.Contest.EndTime
		return
	}

	row := db.DB.QueryRowContext(ctx, "SELECT uid, domain, title, content, type, problem_list, start_time, end_time, open_time, close_time, freeze_time, release_time FROM contest WHERE `uid` = ? AND `domain` = ?", req.UID, req.Domain)
	if err := row.Scan(&res.Contest.Uid, &res.Contest.Domain, &res.Contest.Title, &res.Contest.Description, &res.Contest.Type, &res.Contest.ProblemList, &res.Contest.StartTime, &res.Contest.EndTime, &res.Contest.OpenTime, &res.Contest.CloseTime, &res.Contest.FreezeTime, &res.Contest.ReleaseTime); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	current := time.Now()
	if current.Before(res.Contest.StartTime) && !(state["role"] == "ADMIN" || state["role"] == "ROOT") {
		res.Contest.ProblemList = "[]" // hide problem list
	}
}

func getContestProblem(ctx context.Context, req *senrenrpc.GetContestProblemRequest, state map[string]string, res *senrenrpc.GetContestProblemResponse) {
	contest := &base.Contest{}

	row := db.DB.QueryRowContext(ctx, "SELECT problem_list, start_time FROM contest WHERE `uid` = ? AND `domain` = ?", req.UID, req.Domain)
	if err := row.Scan(&contest.ProblemList, &contest.StartTime); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	current := time.Now()
	if current.Before(contest.StartTime) && !(state["role"] == "ADMIN" || state["role"] == "ROOT") {
		contest.ProblemList = "[]" // hide problem list
	}

	if req.Filter == "" {
		req.Filter = "A"
	}

	prob := int(req.Filter[0] - 'A')

	probs := make([]*base.Problem, 0)
	json.Unmarshal([]byte(contest.ProblemList), &probs)

	if prob < 0 || prob >= len(probs) {
		res.Success = false
		res.Error = "Problem not found"
	}

	probUid := probs[prob].Uid

	// temp lift role
	// trole := state["role"]
	state["role"] = "ADMIN"

	tres := &senrenrpc.GetProblemResponse{}

	getProblem(ctx, &senrenrpc.GetProblemRequest{
		Domain: req.Domain,
		UID:    probUid,
	}, state, tres)

	res.Problem = tres.Problem
	res.Success = tres.Success
	res.Error = tres.Error

	if res.Success {
		res.Problem.RootUid = ""
		res.Problem.Alias = req.Filter[0:1]
		res.Problem.ProblemCI = ""
		res.Problem.Title = probs[prob].Title
	}
}

func getContestSubmissions(ctx context.Context, req *senrenrpc.GetContestSubmissionsRequest, state map[string]string, res *senrenrpc.GetContestSubmissionsResponse) {
	contest := &base.Contest{}

	tFilter := strings.Split(req.Filter, "|")
	tFilter = append(tFilter, []string{"", "", ""}...)

	row := db.DB.QueryRowContext(ctx, "SELECT uid, freeze_time, release_time FROM contest WHERE `uid` = ? AND `domain` = ?", tFilter[0], req.Domain)
	if err := row.Scan(&contest.Uid, &contest.FreezeTime, &contest.ReleaseTime); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	current := time.Now()
	hideState := true
	if contest.ReleaseTime.After(current) && !(state["role"] == "ADMIN" || state["role"] == "ROOT") {
		state["enable_contest"] = "U"
	} else {
		state["enable_contest"] = "A"
		hideState = false
	}

	state["contest_uid"] = contest.Uid
	req.Filter = tFilter[1]

	state["extra_config"] = tFilter[2]

	getSubmissions(ctx, req, state, res)

	if hideState {
		for _, t := range res.Submissions {
			if t.UserUid == state["uid"] || !hideState {
				continue
			}
			t.ExecuteMemory = -1
			t.ExecuteTime = -1
			t.Uid = ""
			if hideState && t.SubmitTime.After(contest.FreezeTime) {
				t.Verdict = "PENDING"
				t.Status = "PENDING"
				continue
			} else if t.Status != "PENDING" && t.Verdict != "AC" && t.Verdict != "CE" {
				t.Verdict = "WA"
			}
		}
	}
}

func getContestSubmission(ctx context.Context, req *senrenrpc.GetContestSubmissionRequest, state map[string]string, res *senrenrpc.GetContestSubmissionResponse) {
	contest := &base.Contest{}

	tFilter := strings.Split(req.Filter, "|")
	tFilter = append(tFilter, []string{"", "", ""}...)

	row := db.DB.QueryRowContext(ctx, "SELECT uid, freeze_time, release_time FROM contest WHERE `uid` = ? AND `domain` = ?", tFilter[0], req.Domain)
	if err := row.Scan(&contest.Uid, &contest.FreezeTime, &contest.ReleaseTime); err != nil {
		res.Success = false
		res.Error = "failed to get contest: " + err.Error()
		return
	}

	current := time.Now()
	hideState := true
	if contest.ReleaseTime.After(current) && !(state["role"] == "ADMIN" || state["role"] == "ROOT") {
		state["enable_contest"] = "U"
	} else {
		state["enable_contest"] = "A"
		hideState = false
	}

	state["contest_uid"] = contest.Uid
	req.Filter = tFilter[1]

	state["extra_config"] = tFilter[2]

	getSubmission(ctx, req, state, res)

	if hideState && res.Submission != nil {
		if res.Submission.UserUid != state["uid"] {
			res.Submission = nil
			res.Success = false
			res.Error = "not found"
		} else {
			res.Submission.ExecuteMemory = -1
			res.Submission.ExecuteTime = -1
			res.Submission.JudgerResponse = "{}"
		}
	}
}

func createContestSubmission(ctx context.Context, req *senrenrpc.CreateContestSubmissionRequest, state map[string]string, res *senrenrpc.CreateContestSubmissionResponse) {
	contest := &base.Contest{}

	row := db.DB.QueryRowContext(ctx, "SELECT problem_list, start_time FROM contest WHERE `uid` = ? AND `domain` = ?", req.Submission.ContestUid, req.Submission.Domain)
	if err := row.Scan(&contest.ProblemList, &contest.StartTime); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	current := time.Now()
	if current.Before(contest.StartTime) && !(state["role"] == "ADMIN" || state["role"] == "ROOT") {
		contest.ProblemList = "[]" // hide problem list
	}

	if req.Submission.ProblemUid == "" {
		req.Submission.ProblemUid = "A"
	}

	prob := int(req.Submission.ProblemUid[0] - 'A')

	probs := make([]*base.Problem, 0)
	json.Unmarshal([]byte(contest.ProblemList), &probs)

	if prob < 0 || prob >= len(probs) {
		res.Success = false
		res.Error = "Problem not found"
	}

	req.Submission.ProblemUid = probs[prob].Uid

	state["contest"] = req.Submission.ContestUid

	createSubmission(ctx, req, state, res)
}

func getContests(ctx context.Context, req *senrenrpc.GetContestsRequest, state map[string]string, res *senrenrpc.GetContestsResponse) {
	// TODO: filter
	// filters := strings.Split(req.Filter, ",")
	res.Contests = make([]*base.Contest, 0)
	rows, err := db.DB.QueryContext(ctx, "SELECT uid, domain, title, content, type, problem_list, start_time, end_time, open_time, close_time, freeze_time, release_time FROM contest WHERE `domain` = ?", req.Domain)

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Success = true

	for rows.Next() {
		tContest := &base.Contest{}
		if err := rows.Scan(&tContest.Uid, &tContest.Domain, &tContest.Title, &tContest.Description, &tContest.Type, &tContest.ProblemList, &tContest.StartTime, &tContest.EndTime, &tContest.OpenTime, &tContest.CloseTime, &tContest.FreezeTime, &tContest.ReleaseTime); err != nil {
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
	if state["role"] != "ADMIN" && state["role"] != "ROOT" {
		res.Success = false
		res.Error = "Forbidden"
		return
	}

	dbExec := "UPDATE contest SET title = ?, content = ?, type = ?, problem_list = ?, start_time = ?, end_time = ?, open_time = ?, close_time = ?, freeze_time = ?, release_time = ? WHERE `uid` = ? AND `domain` = ?"

	tDomain := senrenrpc.Domain(req.Contest.Domain)
	tDomain.ConvertDomain()
	req.Contest.Domain = string(tDomain)

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
	res.Domain = senrenrpc.Domain(req.Contest.Domain)
	res.UID = req.Contest.Uid
}
