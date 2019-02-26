package pciend

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/erjiaqing/senren2/pkg/httpreq"

	"github.com/erjiaqing/senren2/pkg/pcidb"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/pcirpc"
	"github.com/sirupsen/logrus"
)

// type PCITaskItem struct {
// 	Uid        int64     `json:"uid"`
// 	Problem    int64     `json:"problem"`
// 	Creator    string    `json:"creator"`
// 	Status     string    `json:"status"`
// 	Desc       string    `json:"desc"`
// 	Result     string    `json:"result"`
// 	Callback   string    `json:"callback"`
// 	CreateTime time.Time `json:"create_time"`
// 	FinishTime time.Time `json:"finish_time"`
// }

func createSubmissionTask(ctx context.Context, req *pcirpc.CreateSubmissionTaskRequest, state map[string]string, res *pcirpc.CreateSubmissionTaskResponse) {
	descBytes, err := json.Marshal(req.Code)
	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	current := time.Now()

	tProb, _ := strconv.Atoi(state["PROB"])
	tTask := &base.PCITaskItem{
		Problem:  int64(tProb),
		Status:   "PENDING",
		Desc:     string(descBytes),
		Result:   "{}",
		Callback: req.Callback,
	}

	r, err := pcidb.PCIDB.ExecContext(ctx, "INSERT INTO task (problem, creator, state, taskdesc, result, create_at, finish_at, callback) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		tTask.Problem, ".api", "PENDING", tTask.Desc, tTask.Result, current, current, tTask.Callback)

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Success = true
	res.Uid, _ = r.LastInsertId()
}

// TODO

func createProblemTestTask(ctx context.Context, req *pcirpc.CreateProblemTestTaskRequest, state map[string]string, res *pcirpc.CreateProblemTestTaskResponse) {
	tProb, _ := strconv.Atoi(state["PROB"])
	req.Desc.ProblemRepo = state["REPO"]
	req.Desc.PCITask.Type = "build"

	descBytes, err := json.Marshal(req.Desc)
	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	current := time.Now()

	tTask := &base.PCITaskItem{
		Problem:  int64(tProb),
		Status:   "PENDING",
		Desc:     string(descBytes),
		Result:   "{}",
		Callback: req.Callback,
	}

	r, err := pcidb.PCIDB.ExecContext(ctx, "INSERT INTO task (problem, creator, state, taskdesc, result, create_at, finish_at, callback) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		tTask.Problem, ".api", "PENDING", tTask.Desc, tTask.Result, current, current, tTask.Callback)

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Success = true
	res.Uid, _ = r.LastInsertId()
}

func getTask(ctx context.Context, req *pcirpc.GetPCITaskRequest, state map[string]string, res *pcirpc.GetPCITaskResponse) {
	tx, err := pcidb.PCIDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	var rowRes *sql.Row

	if state["PROB"] == "-1" {
		logrus.Debug("--> Query New Task")
		rowRes = tx.QueryRowContext(ctx, "SELECT uid, problem, creator, state, taskdesc, result, create_at, finish_at FROM task WHERE state = ? ORDER BY create_at ASC LIMIT 1", "PENDING")
	} else {
		logrus.Debug("--> Query Exist Task")
		rowRes = tx.QueryRowContext(ctx, "SELECT uid, problem, creator, state, taskdesc, result, create_at, finish_at FROM task WHERE problem = ? AND uid = ?", state["problem"], req.TaskId)
	}

	retTask := &base.PCITaskItem{}

	if err := rowRes.Scan(&retTask.Uid, &retTask.Problem, &retTask.Creator, &retTask.Status, &retTask.Desc, &retTask.Result, &retTask.CreateTime, &retTask.FinishTime); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	if state["PROB"] == "-1" {
		if _, err := tx.ExecContext(ctx, "UPDATE task SET state = ? WHERE uid = ?", "DOWNLOADED", retTask.Uid); err != nil {
			res.Success = false
			res.Error = err.Error()
			return
		}
	}

	if err := tx.Commit(); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Task = retTask
	res.Success = true
}

func updateTask(ctx context.Context, req *pcirpc.UpdatePCITaskRequest, state map[string]string, res *pcirpc.UpdatePCITaskResponse) {
	if state["PROB"] != "-1" {
		res.Success = false
		res.Error = "forbidden"
		return
	}

	rowRes := pcidb.PCIDB.QueryRowContext(ctx, "SELECT uid, problem, callback FROM task WHERE uid = ?", req.Task.Uid)

	retTask := &base.PCITaskItem{}

	if err := rowRes.Scan(&retTask.Uid, &retTask.Problem, &retTask.Callback); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	if _, err := pcidb.PCIDB.ExecContext(ctx, "UPDATE task SET state = ?, result = ?, finish_at = ? WHERE uid = ?", req.Task.Status, req.Task.Result, time.Now(), req.Task.Uid); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	retTask.Status = req.Task.Status
	retTask.Result = req.Task.Result
	retTask.FinishTime = time.Now()
	retTask.Uid = req.Task.Uid
	retTask.Problem = req.Task.Problem

	// TODO: call callback
	if retTask.Callback != "" {
		for i := uint(0); i < 5; i++ {
			_, code, err := httpreq.POSTJson(retTask.Callback, retTask)
			if code >= 300 || err != nil {
				time.Sleep((1 << i) * time.Second / 10)
			} else {
				break
			}
		}
	}

	res.Success = true
}
