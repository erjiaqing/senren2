package senrenrpc

import "github.com/erjiaqing/senren2/pkg/types/base"

type GetHomeworkRequest GetDomainObjectRequest

type GetHomeworkResponse struct {
	SuccessError `json:"result"`
	Homework     *base.Homework `json:"homework"`
}

type GetHomeworksRequest GetDomainObjectsRequest

type GetHomeworksResponse struct {
	SuccessError `json:"result"`
	Homeworks    []*base.Homework `json:"homeworks"`
}

type CreateHomeworkRequest struct {
	Session  `json:"session"`
	Homework base.Homework `json:"homework"`
}

type CreateHomeworkResponse CreateDomainObjectResponse

type CreateHomeworkSubmissionRequest struct {
	Session            `json:"session"`
	HomeworkSubmission base.HomeworkSubmission `json:"homeworksubmission"`
}

type CreateHomeworkSubmissionResponse CreateDomainObjectResponse

type GetHomeworkSubmissionRequest GetDomainObjectRequest

type GetHomeworkSubmissionResponse struct {
	SuccessError       `json:"result"`
	HomeworkSubmission *base.HomeworkSubmission `json:"homeworksubmission"`
}

type GetHomeworkSubmissionsRequest GetDomainObjectsRequest

type GetHomeworkSubmissionsResponse struct {
	SuccessError        `json:"result"`
	HomeworkSubmissions []*base.HomeworkSubmission `json:"homeworksubmissions"`
}
