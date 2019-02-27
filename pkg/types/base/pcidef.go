package base

import "time"

type PCIProblem struct {
	Uid            int64  `json:"uid"`
	Title          string `json:"title"`
	RemoteURL      string `json:"remote"`
	CurrentVersion string `json:"version"`
	EditSession    string `json:"edit_session"`
	Owner          string `json:"owner"`
	State          string `json:"state"`
}

type ProblemVersionState struct {
	Problem int64     `json:"problem"`
	Version string    `json:"version"`
	State   string    `json:"state"`
	Message string    `json:"message"`
	LogTime time.Time `json:"created"`
}

type PCIACL struct {
	Key           string    `json:"key"`
	PrivateKey    string    `json:"pkey"`
	ProblemUID    int64     `json:"puid"`
	CreateTime    time.Time `json:"create_time"`
	AccessControl string    `json:"access_control"`
}

type PCIACLLog struct {
	LogUID string `json:"uid"`
	Key    string `json:"key"`
	Time   string `json:"time"`
	IP     string `json:"ip"`
	API    string `json:"api"`
}

type PCITaskItem struct {
	Uid        int64     `json:"uid"`
	Problem    int64     `json:"problem"`
	Creator    string    `json:"creator"`
	Status     string    `json:"status"`
	Desc       string    `json:"desc"`
	Result     string    `json:"result"`
	Callback   string    `json:"callback"`
	CreateTime time.Time `json:"create_time"`
	FinishTime time.Time `json:"finish_time"`
}

type PCITask struct {
	Type string `json:"type"`
}

type PCIBuildTaskDesc struct {
	PCITask
	Version     string `json:"lang"`
	ProblemRepo string `json:"repo"`
}

type PCIBuildTaskResult struct {
	Success bool `json:"success"`
}

type PCIJudgeTaskDesc struct {
	PCITask
	Lang string `json:"lang"`
	Code string `json:"code"`
}

type PCIJudgeResult struct {
	Success   bool              `json:"success"`
	Verdict   string            `json:"verdict"`
	ExeTime   float32           `json:"exe_time"`
	ExeMemory uint64            `json:"exe_memory"`
	ExitCode  int32             `json:"exit_code"`
	UsedTime  uint64            `json:"used_time"`
	Score     int               `json:"score"`
	FullScore int               `json:"full_score"`
	Detail    []*PCIJudgeDetail `json:"detail"`
}

type PCIJudgeDetail struct {
	Name       string  `json:"name"`
	Input      string  `json:"input,omitempty"`
	Output     string  `json:"output,omitempty"`
	Answer     string  `json:"answer,omitempty"`
	Comment    string  `json:"comment,omitempty"`
	Score      int     `json:"score"`
	Verdict    string  `json:"verdict"`
	ExeTime    float32 `json:"exe_time"`
	ExeMemory  uint64  `json:"exe_memory"`
	ExitCode   int32   `json:"exit_code"`
	ExitSignal int32   `json:"exit_signal"`
}
