package pciend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/erjiaqing/senren2/pkg/router"
	"github.com/erjiaqing/senren2/pkg/types/pcirpc"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var perms map[string]string
var senrenServ string

func init() {
	perms = make(map[string]string)

	r2 := router.R.PathPrefix("/pci_problem").Subrouter()
	r2.HandleFunc("/problemUpdate/{problem}", problemUpdate)

	r := router.R.PathPrefix("/pci").Subrouter()
	r.HandleFunc("/{method}", endpointsRouter)

	logrus.Info("Init routes of pci")

	perms["getProblem"] = ".PROBLEM.READ"
	perms["getProblemDescription"] = ".PROBLEM.READ"
	perms["getProblemAccessKeys"] = "."
	perms["createProblemAccessKey"] = "."
	perms["createProblemEditSession"] = ".PROBLEM"
	perms["closeProblemEditSession"] = ".PROBLEM"
	perms["createSubmissionTask"] = ".PROBLEM.SUBMISSION.WRITE"
	perms["createProblemTestTask"] = ".PROBLEM.TASK.WRITE"

	senrenServ = os.Getenv("SENREN_SERV")
	if senrenServ == "" {
		senrenServ = "http://127.0.0.1:8080"
	}
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
	case "getProblems":
		req = &pcirpc.GetProblemsRequest{}
		res = &pcirpc.GetProblemsResponse{}
	case "getProblemDescription":
		req = &pcirpc.GetProblemDescriptionRequest{}
		res = &pcirpc.GetProblemDescriptionResponse{}
	case "getProblemVersions":
		req = &pcirpc.GetProblemVersionsRequest{}
		res = &pcirpc.GetProblemVersionsResponse{}
	case "getProblemAccessKeys":
		req = &pcirpc.GetProblemAccessKeysRequest{}
		res = &pcirpc.GetProblemAccessKeysResponse{}
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
	case "getTask":
		req = &pcirpc.GetPCITaskRequest{}
		res = &pcirpc.GetPCITaskResponse{}
	case "updateTask":
		req = &pcirpc.UpdatePCITaskRequest{}
		res = &pcirpc.UpdatePCITaskResponse{}

	case "loginBySenrenSid":
		req = &senrenrpc.GetPCISidRequest{}
		res = &senrenrpc.GetPCISidResponse{}
	case "loginToSenren":
		req = &senrenrpc.AuthRequest{}
		res = &senrenrpc.GetPCISidResponse{}
	case "isLogin":
		req = &pcirpc.Session{}
		res = &senrenrpc.WhoAmIResponse{}
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

	if pr, ok := req.(pcirpc.HasSession); ok {
		checkLogin(ctx, pr, state)
	}

	if pr, ok := req.(pcirpc.HasProblemAccessKey); ok {
		puid := int64(0)
		if pr2, ok := req.(pcirpc.HasProblemId); ok {
			puid = pr2.GetId()
		}
		resolveProblemAccessKey(ctx, pr, puid, state)

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
	}

	switch params["method"] {
	case "getProblem":
		getProblem(ctx, req.(*pcirpc.GetProblemRequest), state, res.(*pcirpc.GetProblemResponse))
	case "getProblems":
		getProblems(ctx, req.(*pcirpc.GetProblemsRequest), state, res.(*pcirpc.GetProblemsResponse))
	case "getProblemVersions":
		getProblemVersions(ctx, req.(*pcirpc.GetProblemVersionsRequest), state, res.(*pcirpc.GetProblemVersionsResponse))
	case "getProblemAccessKeys":
		getProblemAccessKeys(ctx, req.(*pcirpc.GetProblemAccessKeysRequest), state, res.(*pcirpc.GetProblemAccessKeysResponse))
	case "getProblemDescription":
		getProblemDescription(ctx, req.(*pcirpc.GetProblemDescriptionRequest), state, res.(*pcirpc.GetProblemDescriptionResponse))
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
	case "getTask":
		getTask(ctx, req.(*pcirpc.GetPCITaskRequest), state, res.(*pcirpc.GetPCITaskResponse))
	case "updateTask":
		updateTask(ctx, req.(*pcirpc.UpdatePCITaskRequest), state, res.(*pcirpc.UpdatePCITaskResponse))
	case "loginBySenrenSid":
		loginBySenrenSid(ctx, req.(*senrenrpc.GetPCISidRequest), state, res.(*senrenrpc.GetPCISidResponse))
	case "loginToSenren":
		loginToSenren(ctx, req.(*senrenrpc.AuthRequest), state, res.(*senrenrpc.GetPCISidResponse))
	case "isLogin":
		isLogin(ctx, req.(*pcirpc.Session), state, res.(*senrenrpc.WhoAmIResponse))
	}

	wbody, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(wbody)
}
