package pciend

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/erjiaqing/senren2/pkg/pcidb"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/pcirpc"
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

func createProblemTestTask(ctx context.Context, req *pcirpc.CreateProblemTestTaskRequest, state map[string]string, res *pcirpc.CreateProblemTestTaskResponse) {

}
