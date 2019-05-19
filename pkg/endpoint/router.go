package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/erjiaqing/senren2/pkg/util"

	"github.com/gorilla/mux"

	"github.com/erjiaqing/senren2/pkg/router"
	"github.com/sirupsen/logrus"
)

var selfURL string
var pciURL string

func init() {
	r := router.R.PathPrefix("/class").Subrouter()
	r.HandleFunc("/{method}", endpointsRouter)

	r2 := router.R.PathPrefix("/attachments").Subrouter()
	r2.HandleFunc("/uploadHomework", uploadHomework)
	r2.HandleFunc("/downloadHomework/{token}/{filename}", downloadHomework)
	r2.HandleFunc("/taskOutput", taskOutput)

	r3 := router.R.PathPrefix("/pcicallback").Subrouter()
	r3.HandleFunc("/taskcallback/{task}", taskCallback)

	logrus.Info("Init routes of class")

	selfURL = os.Getenv("SENREN_SELF_URL")
	if selfURL == "" {
		selfURL = "http://127.0.0.1:8080"
	}

	pciURL = os.Getenv("PCI_SERV")
	if pciURL == "" {
		pciURL = "http://127.0.0.1:8079"
	}
}

func taskCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)

	taskId := util.CheckSessionTime(params["task"], 10*365*24*time.Hour) // almost forever
	if taskId == "" {
		w.WriteHeader(404)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}

	callBacktask := &base.PCITaskItem{}
	//logrus.Info(string(reqBody))

	if err := json.Unmarshal(reqBody, callBacktask); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}

	if callBacktask.Status != "FINISHED" {
		w.Write([]byte("ignored"))
		return
	}

	judge := &base.PCIJudgeResult{}
	if err := json.Unmarshal([]byte(callBacktask.Result), judge); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}

	cemsg := ""

	for _, v := range judge.Detail {
		if v.Name == "compile" {
			cemsg = v.Output
			break
		}
	}

	if _, err := db.DB.ExecContext(ctx, "UPDATE submission SET execute_time = ?, execute_memory = ?, state = ?, verdict = ?, judger_response = ?, ce_message = ?, judge_time = ? WHERE uid = ?",
		int(judge.ExeTime*1000), judge.ExeMemory, "JUDGED", judge.Verdict, callBacktask.Result, cemsg, time.Now(), taskId); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}

	w.Write([]byte("committed"))
	return
}

func uploadHomework(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
		if err := recover(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprint(err)))
		}
	}()

	req := &senrenrpc.CreateHomeworkSubmissionRequest{}
	res := &uploadHomeworkResponse{}

	req.Session.Sid = r.Header.Get("UPLOAD_SESSION")
	req.HomeworkSubmission.Domain = r.Header.Get("UPLOAD_DOMAIN")

	doHomeworkUpload(ctx, r, req, res)
	w.Header().Set("Content-Type", "application/json")
	wbody, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(wbody)
}

// TODO context
func endpointsRouter(w http.ResponseWriter, r *http.Request) {
	var req, res interface{}

	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
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
	case "authUser":
		req = &senrenrpc.AuthRequest{}
		res = &senrenrpc.AuthResponse{}
	case "logoutUser":
		req = &senrenrpc.LogoutRequest{}
		res = &senrenrpc.LogoutResponse{}
	case "whoami":
		req = &senrenrpc.WhoAmIRequest{}
		res = &senrenrpc.WhoAmIResponse{}
	case "getContest":
		req = &senrenrpc.GetContestRequest{}
		res = &senrenrpc.GetContestResponse{}
	case "getContests":
		req = &senrenrpc.GetContestsRequest{}
		res = &senrenrpc.GetContestsResponse{}
	case "getContestProblem":
		req = &senrenrpc.GetContestProblemRequest{}
		res = &senrenrpc.GetContestProblemResponse{}
	case "getContestSubmissions":
		req = &senrenrpc.GetContestSubmissionsRequest{}
		res = &senrenrpc.GetContestSubmissionsResponse{}
	case "getContestSubmission":
		req = &senrenrpc.GetContestSubmissionRequest{}
		res = &senrenrpc.GetContestSubmissionResponse{}
	case "createContestSubmission":
		req = &senrenrpc.CreateContestSubmissionRequest{}
		res = &senrenrpc.CreateContestSubmissionResponse{}
	case "createContest":
		req = &senrenrpc.CreateContestRequest{}
		res = &senrenrpc.CreateContestResponse{}
	case "getHomework":
		req = &senrenrpc.GetHomeworkRequest{}
		res = &senrenrpc.GetHomeworkResponse{}
	case "getHomeworks":
		req = &senrenrpc.GetHomeworksRequest{}
		res = &senrenrpc.GetHomeworksResponse{}
	case "getHomeworkSubmission":
		req = &senrenrpc.GetHomeworkSubmissionRequest{}
		res = &senrenrpc.GetHomeworkSubmissionResponse{}
	case "getHomeworkSubmissions":
		req = &senrenrpc.GetHomeworkSubmissionsRequest{}
		res = &senrenrpc.GetHomeworkSubmissionsResponse{}
	case "createHomework":
		req = &senrenrpc.CreateHomeworkRequest{}
		res = &senrenrpc.CreateHomeworkResponse{}
	case "createHomeworkSubmission":
		req = &senrenrpc.CreateHomeworkSubmissionRequest{}
		res = &senrenrpc.CreateHomeworkSubmissionResponse{}
	case "getHomeworkSubmissionKey":
		req = &senrenrpc.GetHomeworkSubmissionKeyRequest{}
		res = &senrenrpc.GetHomeworkSubmissionKeyResponse{}
	case "setHomeworkScore":
		req = &senrenrpc.SetHomeworkScoreRequest{}
		res = &senrenrpc.SetHomeworkScoreResponse{}
	case "packHomeworkSubmissions":
		req = &senrenrpc.PackHomeworkSubmissionsRequest{}
		res = &senrenrpc.PackHomeworkSubmissionsResponse{}
	case "getProblem":
		req = &senrenrpc.GetProblemRequest{}
		res = &senrenrpc.GetProblemResponse{}
	case "getProblems":
		req = &senrenrpc.GetProblemsRequest{}
		res = &senrenrpc.GetProblemsResponse{}
	case "getPCIDescription":
		req = &senrenrpc.GetPCIDescriptionRequest{}
		res = &senrenrpc.GetPCIDescriptionResponse{}
	case "createProblem":
		req = &senrenrpc.CreateProblemRequest{}
		res = &senrenrpc.CreateProblemResponse{}
	case "createSubmission":
		req = &senrenrpc.CreateSubmissionsRequest{}
		res = &senrenrpc.CreateSubmissionsResponse{}
	case "getSubmission":
		req = &senrenrpc.GetSubmissionRequest{}
		res = &senrenrpc.GetSubmissionResponse{}
	case "getSubmissions":
		req = &senrenrpc.GetSubmissionsRequest{}
		res = &senrenrpc.GetSubmissionsResponse{}
	case "getDomain":
		req = &senrenrpc.GetDomainRequest{}
		res = &senrenrpc.GetDomainResponse{}
	case "getDomains":
		req = &senrenrpc.GetDomainsRequest{}
		res = &senrenrpc.GetDomainsResponse{}
	case "createDomain":
		req = &senrenrpc.CreateDomainRequest{}
		res = &senrenrpc.CreateDomainResponse{}
	case "getDomainInvite":
		req = &senrenrpc.GetDomainInviteRequest{}
		res = &senrenrpc.GetDomainInviteResponse{}
	case "getDomainInvites":
		req = &senrenrpc.GetDomainInvitesRequest{}
		res = &senrenrpc.GetDomainInvitesResponse{}
	case "createDomainInvite":
		req = &senrenrpc.CreateDomainInviteRequest{}
		res = &senrenrpc.CreateDomainInviteResponse{}
	case "joinDomain":
		req = &senrenrpc.JoinDomainRequest{}
		res = &senrenrpc.JoinDomainResponse{}
	case "getDomainUser":
		req = &senrenrpc.GetDomainUserRequest{}
		res = &senrenrpc.GetDomainUserResponse{}
	case "getDomainUsers":
		req = &senrenrpc.GetDomainUsersRequest{}
		res = &senrenrpc.GetDomainUsersResponse{}
	case "updateDomainUser":
		req = &senrenrpc.UpdateDomainUserRequest{}
		res = &senrenrpc.UpdateDomainUserResponse{}
	case "getPCISid":
		req = &senrenrpc.GetPCISidRequest{}
		res = &senrenrpc.GetPCISidResponse{}
	case "getTask":
		req = &senrenrpc.GetTaskRequest{}
		res = &senrenrpc.GetTaskResponse{}
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

	domain := "0000000000000000"
	// do middleware

	if dom, ok := req.(senrenrpc.HasDomain); ok {
		dom.ConvertDomain()
		domain = dom.GetDomain()
		state["domain"] = domain
	}

	if session, ok := req.(senrenrpc.HasSession); ok {
		checkLogin(ctx, session, domain, state)
	}

	switch params["method"] {
	case "authUser":
		authUser(ctx, req.(*senrenrpc.AuthRequest), state, res.(*senrenrpc.AuthResponse))
	case "logoutUser":
		logoutUser(ctx, req.(*senrenrpc.LogoutRequest), state, res.(*senrenrpc.LogoutResponse))
	case "whoami":
		whoami(ctx, req.(*senrenrpc.WhoAmIRequest), state, res.(*senrenrpc.WhoAmIResponse))
	case "getContest":
		getContest(ctx, req.(*senrenrpc.GetContestRequest), state, res.(*senrenrpc.GetContestResponse))
	case "getContests":
		getContests(ctx, req.(*senrenrpc.GetContestsRequest), state, res.(*senrenrpc.GetContestsResponse))
	case "getContestProblem":
		getContestProblem(ctx, req.(*senrenrpc.GetContestProblemRequest), state, res.(*senrenrpc.GetContestProblemResponse))
	case "getContestSubmissions":
		getContestSubmissions(ctx, req.(*senrenrpc.GetContestSubmissionsRequest), state, res.(*senrenrpc.GetContestSubmissionsResponse))
	case "getContestSubmission":
		getContestSubmission(ctx, req.(*senrenrpc.GetContestSubmissionRequest), state, res.(*senrenrpc.GetContestSubmissionResponse))
	case "createContestSubmission":
		createContestSubmission(ctx, req.(*senrenrpc.CreateContestSubmissionRequest), state, res.(*senrenrpc.CreateContestSubmissionResponse))
	case "createContest":
		createContest(ctx, req.(*senrenrpc.CreateContestRequest), state, res.(*senrenrpc.CreateContestResponse))
	case "getHomework":
		getHomework(ctx, req.(*senrenrpc.GetHomeworkRequest), state, res.(*senrenrpc.GetHomeworkResponse))
	case "getHomeworks":
		getHomeworks(ctx, req.(*senrenrpc.GetHomeworksRequest), state, res.(*senrenrpc.GetHomeworksResponse))
	case "getHomeworkSubmission":
		getHomeworkSubmission(ctx, req.(*senrenrpc.GetHomeworkSubmissionRequest), state, res.(*senrenrpc.GetHomeworkSubmissionResponse))
	case "getHomeworkSubmissions":
		getHomeworkSubmissions(ctx, req.(*senrenrpc.GetHomeworkSubmissionsRequest), state, res.(*senrenrpc.GetHomeworkSubmissionsResponse))
	case "getHomeworkSubmissionKey":
		getHomeworkSubmissionKey(ctx, req.(*senrenrpc.GetHomeworkSubmissionKeyRequest), state, res.(*senrenrpc.GetHomeworkSubmissionKeyResponse))
	case "setHomeworkScore":
		setHomeworkScore(ctx, req.(*senrenrpc.SetHomeworkScoreRequest), state, res.(*senrenrpc.SetHomeworkScoreResponse))
	case "packHomeworkSubmissions":
		packHomeworkSubmissions(ctx, req.(*senrenrpc.PackHomeworkSubmissionsRequest), state, res.(*senrenrpc.PackHomeworkSubmissionsResponse))
	case "createHomework":
		createHomework(ctx, req.(*senrenrpc.CreateHomeworkRequest), state, res.(*senrenrpc.CreateHomeworkResponse))
	case "getProblem":
		getProblem(ctx, req.(*senrenrpc.GetProblemRequest), state, res.(*senrenrpc.GetProblemResponse))
	case "getProblems":
		getProblems(ctx, req.(*senrenrpc.GetProblemsRequest), state, res.(*senrenrpc.GetProblemsResponse))
	case "getPCIDescription":
		getPCIDescription(ctx, req.(*senrenrpc.GetPCIDescriptionRequest), state, res.(*senrenrpc.GetPCIDescriptionResponse))
	case "createProblem":
		createProblem(ctx, req.(*senrenrpc.CreateProblemRequest), state, res.(*senrenrpc.CreateProblemResponse))
	case "createSubmission":
		createSubmission(ctx, req.(*senrenrpc.CreateSubmissionsRequest), state, res.(*senrenrpc.CreateSubmissionsResponse))
	case "getSubmission":
		getSubmission(ctx, req.(*senrenrpc.GetSubmissionRequest), state, res.(*senrenrpc.GetSubmissionResponse))
	case "getSubmissions":
		getSubmissions(ctx, req.(*senrenrpc.GetSubmissionsRequest), state, res.(*senrenrpc.GetSubmissionsResponse))
	case "getDomain":
		getDomain(ctx, req.(*senrenrpc.GetDomainRequest), state, res.(*senrenrpc.GetDomainResponse))
	case "getDomains":
		getDomains(ctx, req.(*senrenrpc.GetDomainsRequest), state, res.(*senrenrpc.GetDomainsResponse))
	case "createDomain":
		createDomain(ctx, req.(*senrenrpc.CreateDomainRequest), state, res.(*senrenrpc.CreateDomainResponse))
	case "getDomainInvite":
		getDomainInvite(ctx, req.(*senrenrpc.GetDomainInviteRequest), state, res.(*senrenrpc.GetDomainInviteResponse))
	case "getDomainInvites":
		getDomainInvites(ctx, req.(*senrenrpc.GetDomainInvitesRequest), state, res.(*senrenrpc.GetDomainInvitesResponse))
	case "createDomainInvite":
		createDomainInvite(ctx, req.(*senrenrpc.CreateDomainInviteRequest), state, res.(*senrenrpc.CreateDomainInviteResponse))
	case "joinDomain":
		joinDomain(ctx, req.(*senrenrpc.JoinDomainRequest), state, res.(*senrenrpc.JoinDomainResponse))
	case "getDomainUser":
		getDomainUser(ctx, req.(*senrenrpc.GetDomainUserRequest), state, res.(*senrenrpc.GetDomainUserResponse))
	case "getDomainUsers":
		getDomainUsers(ctx, req.(*senrenrpc.GetDomainUsersRequest), state, res.(*senrenrpc.GetDomainUsersResponse))
	case "updateDomainUser":
		updateDomainUser(ctx, req.(*senrenrpc.UpdateDomainUserRequest), state, res.(*senrenrpc.UpdateDomainUserResponse))
	case "getPCISid":
		getPCISid(ctx, req.(*senrenrpc.GetPCISidRequest), state, res.(*senrenrpc.GetPCISidResponse))
	case "getTask":
		getTask(ctx, req.(*senrenrpc.GetTaskRequest), state, res.(*senrenrpc.GetTaskResponse))
	}

	wbody, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(wbody)
}
