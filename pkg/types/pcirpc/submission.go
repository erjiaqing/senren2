package pcirpc

import "github.com/erjiaqing/senren2/pkg/types/base"

type CreateSubmissionTaskRequest struct {
	ProblemAccessKey
	Callback string                 `json:"callback"`
	Code     *base.PCIJudgeTaskDesc `json:"code"`
}

type CreateSubmissionTaskResponse struct {
	SuccessError
	Uid int64 `json:"uid"`
}

type CreateProblemTestTaskRequest struct {
	ProblemAccessKey
	Callback string                 `json:"callback"`
	Desc     *base.PCIBuildTaskDesc `json:"desc"`
}

type CreateProblemTestTaskResponse struct {
	SuccessError
	Uid int64 `json:"uid"`
}

type GetPCITaskRequest struct {
	ProblemAccessKey
	TaskId int64 `json:"uid"`
	// if problem access key is for problem 0, and task id is 0, then one unjudged task will be returned
}

type GetPCITaskResponse struct {
	SuccessError
	Task *base.PCITaskItem `json:"task"`
}

type UpdatePCITaskRequest struct {
	ProblemAccessKey
	Task *base.PCITaskItem `json:"task"`
	// callback will be invoked whenever this is called
}

type UpdatePCITaskResponse SuccessErrorOnly
