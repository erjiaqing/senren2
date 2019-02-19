package pciend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/erjiaqing/senren2/pkg/router"
	"github.com/erjiaqing/senren2/pkg/types/pcirpc"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var perms map[string]string

func init() {
	perms = make(map[string]string)

	r := router.R.PathPrefix("/pci").Subrouter()
	r.HandleFunc("/{method}", endpointsRouter)

	logrus.Info("Init routes of pci")

	perms["getProblem"] = ".PROBLEM.READ"
	perms["createProblemEditSession"] = ".PROBLEM"
	perms["closeProblemEditSession"] = ".PROBLEM"
	perms["createSubmissionTask"] = ".PROBLEM.SUBMISSION.WRITE"
	perms["createProblemTestTask"] = ".PROBLEM.TASK.WRITE"
}

// TODO context
func endpointsRouter(w http.ResponseWriter, r *http.Request) {
	var req, res interface{}

	ctx := r.Context()

	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprint(err)))
		}
	}()

	params := mux.Vars(r)
	//
	logrus.Infof("--> %s", params["method"])
	// make req & res
	switch params["method"] {
	case "getProblem":
		req = &pcirpc.GetProblemRequest{}
		res = &pcirpc.GetProblemResponse{}
	case "createProblem":
		req = &pcirpc.CreateProblemRequest{}
		res = &pcirpc.CreateProblemResponse{}
	case "createProblemEditSession":
		req = &pcirpc.CreateProblemEditSessionRequest{}
		res = &pcirpc.CreateProblemEditSessionResponse{}
	case "closeProblemEditSession":
		req = &pcirpc.CloseProblemEditSessionRequest{}
		res = &pcirpc.CloseProblemEditSessionResponse{}
	case "createProblemAccessKey":
		req = &pcirpc.CreateProblemAccessKeyRequest{}
		res = &pcirpc.CreateProblemAccessKeyResponse{}
	case "createSubmissionTask":
		req = &pcirpc.CreateSubmissionTaskRequest{}
		res = &pcirpc.CreateSubmissionTaskResponse{}
	case "createProblemTestTask":
		req = &pcirpc.CreateProblemTestTaskRequest{}
		res = &pcirpc.CreateProblemTestTaskResponse{}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	if err := json.Unmarshal(body, req); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	state := make(map[string]string)

	if pr, ok := req.(pcirpc.HasProblemAccessKey); ok {
		resolveProblemAccessKey(ctx, pr, state)
	}

	permsReq, ok := perms[params["method"]]
	if ok {
		rawPerm := strings.Split(permsReq, ".")
		tPerm := ""
		granted := false
		logrus.Debugf("Requres Permission: %s", permsReq)
		if g, tok := state["PERM_."]; tok && g == "G" {
			granted = true
		}
		for _, v := range rawPerm {
			if v == "" {
				continue
			}
			tPerm = tPerm + "." + v
			if g, tok := state["PERM_"+tPerm]; tok && g == "G" {
				granted = true
			}
		}
		if !granted {
			w.WriteHeader(403)
			return
		}
	}

	switch params["method"] {
	case "getProblem":
		getProblem(ctx, req.(*pcirpc.GetProblemRequest), state, res.(*pcirpc.GetProblemResponse))
	case "createProblem":
		createProblem(ctx, req.(*pcirpc.CreateProblemRequest), state, res.(*pcirpc.CreateProblemResponse))
	case "createProblemEditSession":
		createProblemEditSession(ctx, req.(*pcirpc.CreateProblemEditSessionRequest), state, res.(*pcirpc.CreateProblemEditSessionResponse))
	case "closeProblemEditSession":
		closeProblemEditSession(ctx, req.(*pcirpc.CloseProblemEditSessionRequest), state, res.(*pcirpc.CloseProblemEditSessionResponse))
	case "createProblemAccessKey":
		createProblemAccessKey(ctx, req.(*pcirpc.CreateProblemAccessKeyRequest), state, res.(*pcirpc.CreateProblemAccessKeyResponse))
	case "createSubmissionTask":
		createSubmissionTask(ctx, req.(*pcirpc.CreateSubmissionTaskRequest), state, res.(*pcirpc.CreateSubmissionTaskResponse))
	case "createProblemTestTask":
		createProblemTestTask(ctx, req.(*pcirpc.CreateProblemTestTaskRequest), state, res.(*pcirpc.CreateProblemTestTaskResponse))
	}

	wbody, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(wbody)
}
