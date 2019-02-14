package base

import "time"

type User struct {
	Uid string `json:"uid"`
	// user have different UID in different domains, however, they share a single GlobalUID
	// in group 0000000000000000 or root group, GUid = Uid
	// for Domain users, GUid = Uid, and no other records in database
	// user can bind a domain user to a exist global user, thus set GUid to Uid in root group
	GUid       string `json:"guid"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Domain     string `json:"domain"`
	Role       string `json:"role"`
	Score      int    `json:"score"`
	Status     string `json:"status"`
	AuthSource string `json:"authsource"`
}

type Problem struct {
	Uid           string    `json:"uid"`
	RootUid       string    `json:"root_uid"`
	Domain        string    `json:"domain"`
	Alias         string    `json:"alias"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ReleaseTime   time.Time `json:"release"`
	ProblemCI     string    `json:"problem_ci"`
	Score         int       `json:"score"`
	LanguageLimit string    `json:"languagelimit"`
}

type Homework struct {
	Uid         string    `json:"uid"`
	Domain      string    `json:"domain"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Attachments string    `json:"attachments"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type HomeworkSubmission struct {
	Uid         string    `json:"uid"`
	Domain      string    `json:"domain"`
	UserUid     string    `json:"user_uid"`
	HomeworkUid string    `json:"homework_uid"`
	Attachments string    `json:"attachments"`
	Nick        string    `json:"nick"`
	CreateTime  time.Time `json:"create_time"`
}

type HomeworkArchiveDescriptor struct {
	Type    string                       `json:"type"`
	Name    string                       `json:"name"`
	Source  string                       `json:"source"`
	Content []*HomeworkArchiveDescriptor `json:"content"`
}

type HomeworkArchiveTask struct {
	ServerTask
	Desc *HomeworkArchiveDescriptor `json:"desc"`
}

type ServerTask struct {
	Type           string `json:"type"`
	State          string `json:"state,omitempty"`
	DownloadURL    string `json:"downloadurl,omitempty"`
	OutputFileName string `json:"output_filename"`
}

type Contest struct {
	Uid         string    `json:"uid"`
	Domain      string    `json:"domain"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	ProblemList string    `json:"problem_list"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	OpenTime    time.Time `json:"open_time"`
	CloseTime   time.Time `json:"close_time"`
	FreezeTime  time.Time `json:"freeze_time"`
	ReleaseTime time.Time `json:"release_time"`
}

type Discussion struct {
	Uid      string `json:"uid"`
	Domain   string `json:"domain"`
	RootUid  string `json:"root_uid"`
	Title    string `json:"title"`
	PostUser string `json:"post_user"`
	Content  string `json:"content"`
	IsPublic bool   `json:"is_public"`
}

type Submission struct {
	Uid            string    `json:"uid"`
	UserUid        string    `json:"user_uid"`
	ProblemUid     string    `json:"problem_uid"`
	ContestUid     string    `json:"contest_uid"`
	Domain         string    `json:"domain"`
	Language       string    `json:"language"`
	Code           string    `json:"code"`
	FileName       string    `json:"file_name"`
	ExecuteTime    int64     `json:"execute_time"`
	ExecuteMemory  int64     `json:"execute_memory"`
	Status         string    `json:"status"`
	Verdict        string    `json:"verdict"`
	Testcase       int       `json:"test_case"`
	Score          int       `json:"score"`
	JudgerResponse string    `json:"judger_response"`
	CEMessage      string    `json:"ce_message"`
	SubmitTime     time.Time `json:"submit_time"`
	JudgeTime      time.Time `json:"judge_time"`

	ProblemTitle string `json:"problem_title"`
	UserName     string `json:"user_name"`
}

type DomainInfo struct {
	Uid         string `json:"uid"`
	ShortName   string `json:"alias"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    string `json:"is_public"`
}

type DomainInvite struct {
	Uid         string    `json:"invite_uid"`
	Description string    `json:"description"`
	Domain      string    `json:"domain"`
	Password    string    `json:"password"`
	ValidTo     time.Time `json:"valid_to"`
	InviteRole  string    `json:"invite_role"`
	InviteState string    `json:"invite_state"`
}

type Message struct {
	Uid      string    `json:"uid"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	SendTime time.Time `json:"send_time"`
	ReadTime time.Time `json:"read_time"`
}
