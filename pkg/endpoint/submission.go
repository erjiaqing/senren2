package endpoint

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/erjiaqing/senren2/pkg/httpreq"
	"github.com/erjiaqing/senren2/pkg/types/pcirpc"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/util"

	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
)

func createSubmission(ctx context.Context, req *senrenrpc.CreateSubmissionsRequest, state map[string]string, res *senrenrpc.CreateSubmissionsResponse) {
	// TODO: 检查用户权限，检查试题存在性，检查语言合法性，推送到pci
	req.Submission.Uid = util.GenUid()
	req.Submission.SubmitTime = time.Now()

	ok := true

	if req.Submission.UserUid, ok = state["uid"]; !ok {
		res.Success = false
		res.Error = "Login Required"
		return
	}

	if state["contest"] != "" {
		req.Submission.ContestUid = state["contest"]
	} else {
		req.Submission.ContestUid = "0000000000000000"
	}

	if _, err := db.DB.ExecContext(ctx, "INSERT INTO submission (uid, user_uid, domain, problem_uid, contest_uid, lang, code, state, verdict, submit_time, judge_time, filename, execute_time, execute_memory, testcase, score, judger_response, ce_message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		req.Submission.Uid, req.Submission.UserUid, req.Submission.Domain, req.Submission.ProblemUid, req.Submission.ContestUid, req.Submission.Language, req.Submission.Code, "PENDING", "PENDING", req.Submission.SubmitTime, req.Submission.SubmitTime, "", -1, -1, -1, -1, "{}", ""); err != nil {
		panic(err)
	}

	probSubkey := ""

	if err := db.DB.QueryRowContext(ctx, "SELECT problemci FROM problem WHERE uid = ?", req.Submission.ProblemUid).Scan(&probSubkey); err != nil {
		panic(err)
	}

	res.Domain = senrenrpc.Domain(req.Submission.Domain)
	res.Success = true
	res.UID = req.Submission.Uid

	pcireq := &pcirpc.CreateSubmissionTaskRequest{
		Callback: fmt.Sprintf("%s/rpc/pcicallback/taskcallback/%s", selfURL, util.SignSession(res.UID)),
		Code: &base.PCIJudgeTaskDesc{
			Lang: req.Submission.Language,
			Code: req.Submission.Code,
		},
	}
	pcireq.Key = probSubkey
	pcireq.Code.Type = "judge"

	for i := uint(0); i < 5; i++ {
		_, code, err := httpreq.POSTJson(fmt.Sprintf("%s/rpc/pci/createSubmissionTask", pciURL), pcireq)
		if code >= 300 || err != nil {
			time.Sleep((1 << i) * time.Second / 10)
		} else {
			return
		}
	}

	res.Success = false
	res.Error = "failed to submit to pci"
}

func getSubmission(ctx context.Context, req *senrenrpc.GetSubmissionRequest, state map[string]string, res *senrenrpc.GetSubmissionResponse) {
	r := &base.Submission{}
	row := db.DB.QueryRowContext(ctx, "SELECT uid, user_uid, domain, problem_uid, contest_uid, lang, code, execute_time, execute_memory, state, verdict, testcase, score, judger_response, ce_message, submit_time, judge_time FROM submission WHERE uid = ? AND domain = ?", req.UID, req.Domain)
	if err := row.Scan(&r.Uid, &r.UserUid, &r.Domain, &r.ProblemUid, &r.ContestUid, &r.Language, &r.Code, &r.ExecuteTime, &r.ExecuteMemory, &r.Status, &r.Verdict, &r.Testcase, &r.Score, &r.JudgerResponse, &r.CEMessage, &r.SubmitTime, &r.JudgeTime); err != nil {
		res.Error = err.Error()
		res.Success = false
		return
	}
	res.Success = true
	res.Submission = r
}

func getSubmissions(ctx context.Context, req *senrenrpc.GetSubmissionsRequest, state map[string]string, res *senrenrpc.GetSubmissionsResponse) {
	r := make([]*base.Submission, 0)
	query := "SELECT submission.uid, submission.user_uid, IFNULL(user.nickname, user.username), submission.domain, submission.problem_uid, problem.title, submission.contest_uid, submission.lang, submission.execute_time, submission.execute_memory, submission.state, submission.verdict, submission.testcase, submission.score, submission.submit_time FROM submission LEFT JOIN problem ON submission.problem_uid = problem.uid LEFT JOIN user ON submission.user_uid = user.uid WHERE submission.domain = ? "

	if state["extra_config"] == "rank" {
		query = "SELECT submission.uid, submission.user_uid, IFNULL(user.nickname, user.username), submission.domain, submission.problem_uid, '', '', '', -1, -1, 'JUDGED', submission.verdict, submission.testcase, submission.score, submission.submit_time FROM submission LEFT JOIN user ON submission.user_uid = user.uid WHERE submission.domain = ? "
	}

	limits := strings.Split(req.Filter, ";")
	limits = append(limits, make([]string, 4)...)
	whereArgs := []interface{}{string(req.Domain)}
	if limits[0] != "" {
		query += " AND submission.problem_uid = ? "
		whereArgs = append(whereArgs, limits[0])
	}
	/*
		if state["enable_contest"] == "U" {
			limits[1] = state["guid"]
		}
	*/
	if limits[1] != "" {
		query += " AND submission.user_uid = ? "
		whereArgs = append(whereArgs, limits[1])
	}
	if limits[2] != "" {
		tquery := strings.Split(limits[2], ",")
		tArg := " 0 = 1 "
		for _, t := range tquery {
			tArg += " OR submission.verdict = ? "
			whereArgs = append(whereArgs, t)
		}
		query += "AND (" + tArg + ")"
	}
	query += " AND submission.contest_uid = ? "
	if state["enable_contest"] == "A" || state["enable_contest"] == "U" {
		whereArgs = append(whereArgs, state["contest_uid"])
	} else {
		whereArgs = append(whereArgs, "0000000000000000")
	}

	if state["extra_config"] != "rank" {
		query = query + fmt.Sprintf(" ORDER BY submit_time DESC LIMIT %d, 50", req.From)
	}

	rows, err := db.DB.QueryContext(ctx, query, whereArgs...)

	if err != nil {
		res.Error = err.Error()
		res.Success = false
		return
	}

	for rows.Next() {
		t := &base.Submission{}
		if err := rows.Scan(&t.Uid, &t.UserUid, &t.UserName, &t.Domain, &t.ProblemUid, &t.ProblemTitle, &t.ContestUid, &t.Language, &t.ExecuteTime, &t.ExecuteMemory, &t.Status, &t.Verdict, &t.Testcase, &t.Score, &t.SubmitTime); err != nil {
			rows.Close()
			res.Error = err.Error()
			res.Success = false
			return
		}
		r = append(r, t)
	}

	res.Success = true
	res.Submissions = r
}
