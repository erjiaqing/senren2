package endpoint

import (
	"context"
	"time"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/util"

	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
)

func createSubmission(ctx context.Context, req *senrenrpc.CreateSubmissionsRequest, state map[string]string, res *senrenrpc.CreateSubmissionsResponse) {
	// TODO: 检查用户权限，检查试题存在性，检查语言合法性，推送到pci
	req.Submission.Uid = util.GenUid()
	req.Submission.SubmitTime = time.Now()
	if _, err := db.DB.ExecContext(ctx, "INSERT INTO submission (uid, user_uid, domain, problem_uid, contest_uid, lang, code, state, verdict, submit_time, judge_time, filename, execute_time, execute_memory, testcase, score, judger_response, ce_message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		req.Submission.Uid, req.Submission.UserUid, req.Submission.Domain, req.Submission.ProblemUid, req.Submission.ContestUid, req.Submission.Language, req.Submission.Code, "PENDING", "PENDING", req.Submission.SubmitTime, req.Submission.SubmitTime, "", -1, -1, -1, -1, "{}", ""); err != nil {
		panic(err)
	}
	res.Domain = senrenrpc.Domain(req.Submission.Domain)
	res.Success = true
	res.UID = req.Submission.Uid
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

}
