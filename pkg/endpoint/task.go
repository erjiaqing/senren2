package endpoint

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
)

const signStr = "417MD2y3nUNI7e5r1CZ8x6LBLzwkMk4y8FHmrC0srJP5YHXPXUXrcYNj4Bfr417MD2y3nUNI7e5r1CZ8FLZ2vk3iTbUaUQ9v0wBk"

func taskOutput(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	taskUID := r.URL.Query().Get("task")
	signTime := r.URL.Query().Get("time")
	sign := r.URL.Query().Get("sign")

	current := time.Now()
	var signTimeUTC int64
	fmt.Sscanf(signTime, "%d", &signTimeUTC)
	signTimeT := time.Unix(signTimeUTC, 0)
	if signTimeT.After(current) {
		w.WriteHeader(403)
		return
	} else if signTimeT.Before(current.Add(-24 * time.Hour)) {
		w.WriteHeader(403)
		return
	}

	sigSignSrc := fmt.Sprintf("task=%s&time=%s&key=%s", taskUID, signTime, signStr)
	sigSignRes := fmt.Sprintf("%x", sha1.Sum([]byte(sigSignSrc)))

	if sigSignRes != sign {
		w.WriteHeader(403)
		return
	}

	var desc, state string
	row := db.DB.QueryRow("SELECT `desc`, `state` FROM task WHERE uid = ?", taskUID)
	if err := row.Scan(&desc, &state); err != nil {
		w.WriteHeader(500)
		return
	}

	if state != "SUCCESS" {
		w.WriteHeader(404)
		w.Write([]byte("Still pending..."))
		return
	}

	task := &base.ServerTask{}
	json.Unmarshal([]byte(desc), task)

	fp, err := os.Open("output/" + taskUID + "/" + task.OutputFileName)

	if err != nil {
		w.WriteHeader(500)
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+task.OutputFileName+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")

	io.Copy(w, fp)
}

func getTask(ctx context.Context, req *senrenrpc.GetTaskRequest, state map[string]string, res *senrenrpc.GetTaskResponse) {
	res.Success = false
	res.Task = &base.ServerTask{}
	var creator = ""
	row := db.DB.QueryRow("SELECT `creator`, `state` FROM task WHERE uid = ?", req.UID)
	if err := row.Scan(&creator, &res.Task.State); err != nil {
		res.Success = false
		res.Error = "Not found"
		return
	}

	if creator != state["guid"] {
		res.Success = false
		res.Error = "Not found"
		return
	}

	if res.Task.State == "SUCCESS" {
		current := time.Now()
		sigSignSrc := fmt.Sprintf("task=%s&time=%d&key=%s", req.UID, current.Unix(), signStr)
		sigSignRes := fmt.Sprintf("%x", sha1.Sum([]byte(sigSignSrc)))
		res.Task.DownloadURL = "/rpc/attachments/taskOutput?" + fmt.Sprintf("task=%s&time=%d&sign=%s", req.UID, current.Unix(), sigSignRes)
	}
	res.Success = true
}
