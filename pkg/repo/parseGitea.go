package repo

import (
	"encoding/json"

	"github.com/erjiaqing/senren2/pkg/types/base"
)

type hookCommitPush struct {
	Commits []*payloadCommit `json:"commits"`
}

func ParseCommitPush(problem int64, body []byte) []*base.ProblemVersionState {
	res := &hookCommitPush{}
	err := json.Unmarshal(body, res)
	if err != nil {
		return []*base.ProblemVersionState{}
	}
	ret := make([]*base.ProblemVersionState, 0)

	for _, ti := range res.Commits {
		if ti.ID == "0000000000000000000000000000000000000000" {
			continue
		}
		ret = append(ret, &base.ProblemVersionState{
			Problem: problem,
			Version: ti.ID,
			LogTime: ti.TimeStamp,
			Message: ti.Message,
			State:   "UNKNOWN",
		})
	}
	return ret
}
