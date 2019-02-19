package repo

import (
	"encoding/json"
	"fmt"

	"github.com/erjiaqing/senren2/pkg/httpreq"
	"github.com/erjiaqing/senren2/pkg/types/base"
)

func GetRepoInfo(problem *base.PCIProblem) (string, string, error) {
	basicinfo, code, err := httpreq.GETAuth(fmt.Sprintf("%s/api/v1/repos/%s", gitServer, problem.RemoteURL), gitUserName, gitUserPass)
	repoinfo := &repo{}

	if err := json.Unmarshal(basicinfo, repoinfo); err != nil {
		return "", "", err
	}

	resp, code, err := httpreq.GETAuth(fmt.Sprintf("%s/api/v1/repos/%s/branch/%s", gitServer, problem.RemoteURL, "master"), gitUserName, gitUserPass)
	if err != nil {
		return repoinfo.DefaultBranch, "", err
	} else if code >= 300 {
		return repoinfo.DefaultBranch, "", fmt.Errorf("Unexpected http response code http 200 expected, %d received (ERRID: HTTP.2xx.300)", code)
	}

	branchinfo := &branch{}

	if err := json.Unmarshal(resp, branchinfo); err != nil {
		return repoinfo.DefaultBranch, "", err
	}

	if branchinfo.Commit == nil {
		return repoinfo.DefaultBranch, "", fmt.Errorf("Cannot get branch info (ERRID: GIT.NULL.COMMIT)")
	}

	return repoinfo.DefaultBranch, branchinfo.Commit.ID, nil
}
