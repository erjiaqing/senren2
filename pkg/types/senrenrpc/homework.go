package senrenrpc

import "github.com/erjiaqing/senren2/pkg/types/base"

type GetHomeworkRequest GetDomainObjectRequest

type GetHomeworkResponse struct {
	SuccessError
	Homework *base.Homework `json:"homework"`
}

type GetHomeworksRequest GetDomainObjectsRequest

type GetHomeworksResponse struct {
	SuccessError
	Homeworks []*base.Homework `json:"homeworks"`
}

type SetHomeworkScoreRequest GetDomainObjectRequest

type SetHomeworkScoreResponse CreateDomainObjectResponse

type CreateHomeworkRequest struct {
	Session
	Homework base.Homework `json:"homework"`
}

type CreateHomeworkResponse CreateDomainObjectResponse

type CreateHomeworkSubmissionRequest struct {
	Session
	HomeworkSubmission base.HomeworkSubmission `json:"homeworksubmission"`
}

type CreateHomeworkSubmissionResponse CreateDomainObjectResponse

type GetHomeworkSubmissionRequest GetDomainObjectRequest

type GetHomeworkSubmissionResponse struct {
	SuccessError
	HomeworkSubmission *base.HomeworkSubmission `json:"homeworksubmission"`
}

type GetHomeworkSubmissionsRequest GetDomainObjectsRequest

type GetHomeworkSubmissionKeyRequest GetDomainObjectRequest

type GetHomeworkSubmissionKeyResponse CreateDomainObjectResponse

type GetHomeworkSubmissionsResponse struct {
	SuccessError
	HomeworkSubmissions []*base.HomeworkSubmission `json:"homeworksubmissions"`
}

type PackHomeworkSubmissionsRequest GetDomainObjectRequest

type PackHomeworkSubmissionsResponse CreateDomainObjectResponse
