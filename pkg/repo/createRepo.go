package repo

import (
	"fmt"
	"net/url"
	"os"

	"github.com/erjiaqing/senren2/pkg/httpreq"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/sirupsen/logrus"
)

type createRepoRequest struct {
	AutoInit    bool   `json:"auth_init"`
	Description string `json:"description"`
	GitIgnore   string `json:"gitignore"`
	License     string `json:"license"`
	Name        string `json:"name"`
	Private     bool   `json:"private"`
	ReadMe      string `json:"readme"`
}

var gitUserName = ""
var gitUserPass = ""
var gitServer = ""
var selfURL = ""

func init() {
	gitUserName = os.Getenv("PCI_GIT_USER")
	if gitUserName == "" {
		gitUserName = "pci"
	}
	gitUserPass = os.Getenv("PCI_GIT_PASS")
	if gitUserPass == "" {
		gitUserPass = "pcipcipci"
	}
	gitServer = os.Getenv("PCI_GIT_SERV")
	if gitServer == "" {
		gitServer = "http://127.0.0.1:3000"
	}
	selfURL = os.Getenv("PCI_SELF_URL")
	if selfURL == "" {
		selfURL = "http://127.0.0.1:8079"
	}

	logrus.Infof("Git User: %s", gitUserName)
	logrus.Infof("Git Pass: %s", gitUserPass)
	logrus.Infof("Git Serv: %s", gitServer)
}

// If problem editor does not found this problem, it should
// 1. copy local template
// 2. create init commit
// 3. push it to server
func CreateProblemRepo(problem *base.PCIProblem) (*base.PCIProblem, error) {
	reqBody := &createRepoRequest{
		AutoInit:    false,
		Description: fmt.Sprintf("Problem %s, created by PCI", problem.Title),
		GitIgnore:   "",
		License:     "",
		Name:        fmt.Sprintf("p%06d", problem.Uid),
		Private:     true,
		ReadMe:      fmt.Sprintf("Problem %s, created by PCI", problem.Title),
	}

	if _, code, err := httpreq.POSTJsonAuth(fmt.Sprintf("%s/api/v1/user/repos", gitServer), reqBody, gitUserName, gitUserPass); err != nil {
		return problem, err
	} else if code >= 300 {
		return problem, fmt.Errorf("Unexpected http response code http 201 expected, %d received (ERRID: HTTP.2xx.300)", code)
	}

	if _, code, err := httpreq.POSTJsonAuth(fmt.Sprintf("%s/api/v1/repos/%s/p%06d/hooks", gitServer, gitUserName, problem.Uid), &createRepoHookRequest{
		Active: true,
		Config: map[string]string{
			"content_type": "json",
			"url":          fmt.Sprintf("%s/rpc/pci_problem/problemUpdate/%d", selfURL, problem.Uid),
		},
		Events: []string{"push"},
		Type:   "gitea",
	}, gitUserName, gitUserPass); err != nil {
		return problem, err
	} else if code >= 300 {
		return problem, fmt.Errorf("Unexpected http response code http 201 expected, %d received (ERRID: HTTP.2xx.300)", code)
	}

	problem.RemoteURL = fmt.Sprintf("%s/p%06d", gitUserName, problem.Uid)

	return problem, nil
}

func FetchProblemDescription(problem string, version string) []byte {
	ret, code, err := httpreq.GETAuth(fmt.Sprintf("%s/%s/raw/commit/%s/description.html", gitServer, problem, version), gitUserName, gitUserPass)
	if code >= 300 || err != nil {
		return nil
	}
	return ret
}

func FetchProblemVersions(problem *base.PCIProblem) ([]string, error) {
	return []string{}, nil
}

func GetCloneAddress(remote string) string {
	baseUrl := fmt.Sprintf("%s/%s.git", gitServer, remote)
	tUrl, err := url.Parse(baseUrl)
	if err != nil {
		return ""
	}

	tUrl.User = url.UserPassword(gitUserName, gitUserPass)
	return tUrl.String()
}
