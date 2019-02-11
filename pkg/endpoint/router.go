package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"

	"github.com/gorilla/mux"

	"github.com/erjiaqing/senren2/pkg/router"
	"github.com/sirupsen/logrus"
)

func init() {
	r := router.R.PathPrefix("/class").Subrouter()
	r.HandleFunc("/{method}", endpointsRouter)

	logrus.Info("Init routes of class")
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
	case "createContest":
		req = &senrenrpc.CreateContestRequest{}
		res = &senrenrpc.CreateContestResponse{}
	case "getHomework":
		req = &senrenrpc.GetHomeworkRequest{}
		res = &senrenrpc.GetHomeworkResponse{}
	case "getHomeworks":
		req = &senrenrpc.GetHomeworksRequest{}
		res = &senrenrpc.GetHomeworksResponse{}
	case "getHomeorkSubmission":
		req = &senrenrpc.GetHomeworkSubmissionRequest{}
		res = &senrenrpc.GetHomeworkSubmissionResponse{}
	case "getHomeorkSubmissions":
		req = &senrenrpc.GetHomeworkSubmissionsRequest{}
		res = &senrenrpc.GetHomeworkSubmissionsResponse{}
	case "createHomework":
		req = &senrenrpc.CreateHomeworkRequest{}
		res = &senrenrpc.CreateHomeworkResponse{}
	case "createHomeworkSubmission":
		req = &senrenrpc.CreateHomeworkSubmissionRequest{}
		res = &senrenrpc.CreateHomeworkSubmissionResponse{}
	case "getProblem":
		req = &senrenrpc.GetProblemRequest{}
		res = &senrenrpc.GetProblemResponse{}
	case "getProblems":
		req = &senrenrpc.GetProblemsRequest{}
		res = &senrenrpc.GetProblemsResponse{}
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
	case "createContest":
		createContest(ctx, req.(*senrenrpc.CreateContestRequest), state, res.(*senrenrpc.CreateContestResponse))
	case "getHomework":
		getHomework(ctx, req.(*senrenrpc.GetHomeworkRequest), state, res.(*senrenrpc.GetHomeworkResponse))
	case "getHomeworks":
		getHomeworks(ctx, req.(*senrenrpc.GetHomeworksRequest), state, res.(*senrenrpc.GetHomeworksResponse))
	case "getHomeorkSubmission":
		getHomeorkSubmission(ctx, req.(*senrenrpc.GetHomeworkSubmissionRequest), state, res.(*senrenrpc.GetHomeworkSubmissionResponse))
	case "getHomeorkSubmissions":
		getHomeorkSubmissions(ctx, req.(*senrenrpc.GetHomeworkSubmissionsRequest), state, res.(*senrenrpc.GetHomeworkSubmissionsResponse))
	case "createHomework":
		createHomework(ctx, req.(*senrenrpc.CreateHomeworkRequest), state, res.(*senrenrpc.CreateHomeworkResponse))
	case "createHomeworkSubmission":
		createHomeworkSubmission(ctx, req.(*senrenrpc.CreateHomeworkSubmissionRequest), state, res.(*senrenrpc.CreateHomeworkSubmissionResponse))
	case "getProblem":
		getProblem(ctx, req.(*senrenrpc.GetProblemRequest), state, res.(*senrenrpc.GetProblemResponse))
	case "getProblems":
		getProblems(ctx, req.(*senrenrpc.GetProblemsRequest), state, res.(*senrenrpc.GetProblemsResponse))
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
	}

	wbody, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(wbody)
}
