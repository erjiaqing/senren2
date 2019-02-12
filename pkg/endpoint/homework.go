package endpoint

import (
	"context"
	"time"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/erjiaqing/senren2/pkg/util"
)

func getHomework(ctx context.Context, req *senrenrpc.GetHomeworkRequest, state map[string]string, res *senrenrpc.GetHomeworkResponse) {
	hw := &base.Homework{}
	if req.UID == "" || req.UID == noUID {
		hw.StartTime = time.Now()
		hw.EndTime = hw.StartTime.Add(7 * 24 * time.Hour)
		res.Success = true
		res.Homework = hw
		return
	}
	qry := db.DB.QueryRowContext(ctx, "SELECT uid, domain, title, content, attachments, start_time, end_time FROM homework WHERE domain = ? AND uid = ?", req.Domain, req.UID)
	if err := qry.Scan(&hw.Uid, &hw.Domain, &hw.Title, &hw.Description, &hw.Attachments, &hw.StartTime, &hw.EndTime); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	res.Success = true
	res.Homework = hw
}

func getHomeworks(ctx context.Context, req *senrenrpc.GetHomeworksRequest, state map[string]string, res *senrenrpc.GetHomeworksResponse) {
	res.Homeworks = make([]*base.Homework, 0)
	rows, err := db.DB.QueryContext(ctx, "SELECT uid, domain, title, start_time, end_time FROM homework WHERE `domain` = ?", req.Domain)

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Success = true

	for rows.Next() {
		tHomework := &base.Homework{}
		if err := rows.Scan(&tHomework.Uid, &tHomework.Domain, &tHomework.Title, &tHomework.StartTime, &tHomework.EndTime); err != nil {
			res.Success = false
			res.Error = err.Error()
			rows.Close()
			break
		}
		res.Homeworks = append(res.Homeworks, tHomework)
	}
}

func getHomeorkSubmission(ctx context.Context, req *senrenrpc.GetHomeworkSubmissionRequest, state map[string]string, res *senrenrpc.GetHomeworkSubmissionResponse) {

}

func getHomeorkSubmissions(ctx context.Context, req *senrenrpc.GetHomeworkSubmissionsRequest, state map[string]string, res *senrenrpc.GetHomeworkSubmissionsResponse) {

}

func createHomework(ctx context.Context, req *senrenrpc.CreateHomeworkRequest, state map[string]string, res *senrenrpc.CreateHomeworkResponse) {
	dbExec := "UPDATE homework SET title = ? , content = ? , attachments = ?, start_time = ?, end_time = ? WHERE uid = ? AND domain = ?"

	// TODO: check domain permission

	tDomain := senrenrpc.Domain(req.Homework.Domain)
	tDomain.ConvertDomain()
	req.Homework.Domain = string(tDomain)

	if req.Homework.Uid == "" || req.Homework.Uid == noUID {
		req.Homework.Uid = util.GenUid()

		dbExec = "INSERT INTO homework (title, content, attachments, start_time, end_time, uid, domain) VALUES (?, ?, ?, ?, ?, ?, ?)"
	}

	if _, err := db.DB.Exec(dbExec, req.Homework.Title, req.Homework.Description, req.Homework.Attachments, req.Homework.StartTime, req.Homework.EndTime, req.Homework.Uid, req.Homework.Domain); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Domain = senrenrpc.Domain(req.Homework.Domain)
	res.UID = req.Homework.Uid
	res.Success = true
}

func createHomeworkSubmission(ctx context.Context, req *senrenrpc.CreateHomeworkSubmissionRequest, state map[string]string, res *senrenrpc.CreateHomeworkSubmissionResponse) {

}
