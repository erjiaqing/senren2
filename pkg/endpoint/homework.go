package endpoint

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/erjiaqing/senren2/pkg/publicchan"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/erjiaqing/senren2/pkg/util"
	"github.com/sirupsen/logrus"
)

func getHomework(ctx context.Context, req *senrenrpc.GetHomeworkRequest, state map[string]string, res *senrenrpc.GetHomeworkResponse) {
	hw := &base.Homework{}
	if req.UID == "" || req.UID == noUID {
		hw.StartTime = time.Now()
		hw.EndTime = hw.StartTime.Add(7 * 24 * time.Hour)
		res.Success = true
		res.Homework = hw
		return
	}
	qry := db.DB.QueryRowContext(ctx, "SELECT uid, domain, title, content, attachments, start_time, end_time FROM homework WHERE domain = ? AND uid = ?", req.Domain, req.UID)
	if err := qry.Scan(&hw.Uid, &hw.Domain, &hw.Title, &hw.Description, &hw.Attachments, &hw.StartTime, &hw.EndTime); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	res.Success = true
	res.Homework = hw
}

func setHomeworkScore(ctx context.Context, req *senrenrpc.SetHomeworkScoreRequest, state map[string]string, res *senrenrpc.SetHomeworkScoreResponse) {
	if !(state["role"] == "ADMIN" || state["role"] == "ROOT") {
		res.Success = false
		res.Error = "Forbidden"
		return
	}
	score, err := strconv.Atoi(req.Filter)
	if err != nil {
		score = -1
	}
	_, err = db.DB.ExecContext(ctx, "UPDATE homework_submission SET score = ? WHERE uid = ? AND domain = ?", score, req.UID, req.Domain)
	if err != nil {
		res.Success = false
		res.Error = err.Error()
	}
	res.Success = true
}

func getHomeworks(ctx context.Context, req *senrenrpc.GetHomeworksRequest, state map[string]string, res *senrenrpc.GetHomeworksResponse) {
	res.Homeworks = make([]*base.Homework, 0)
	rows, err := db.DB.QueryContext(ctx, "SELECT uid, domain, title, start_time, end_time FROM homework WHERE `domain` = ?", req.Domain)

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Success = true

	for rows.Next() {
		tHomework := &base.Homework{}
		if err := rows.Scan(&tHomework.Uid, &tHomework.Domain, &tHomework.Title, &tHomework.StartTime, &tHomework.EndTime); err != nil {
			res.Success = false
			res.Error = err.Error()
			rows.Close()
			break
		}
		res.Homeworks = append(res.Homeworks, tHomework)
	}
}

func getHomeworkSubmission(ctx context.Context, req *senrenrpc.GetHomeworkSubmissionRequest, state map[string]string, res *senrenrpc.GetHomeworkSubmissionResponse) {
	sub := &base.HomeworkSubmission{}
	if req.Filter != "" && (state["role"] == "ADMIN" || state["role"] == "ROOT") {
		row := db.DB.QueryRowContext(ctx, "SELECT homework_submission.uid, homework_submission.domain, homework_submission.useruid, homework_submission.homework_uid, homework_submission.attachments, homework_submission.create_time, homework_submission.score, user.nickname FROM homework_submission, user WHERE homework_submission.uid = ? AND homework_submission.domain = ? AND homework_submission.useruid = user.uid", req.Filter, req.Domain)
		if err := row.Scan(&sub.Uid, &sub.Domain, &sub.UserUid, &sub.HomeworkUid, &sub.Attachments, &sub.CreateTime, &sub.Score, &sub.Nick); err != nil {
			res.Success = false
			res.Error = "Not found"
			return
		}
	} else {
		row := db.DB.QueryRowContext(ctx, "SELECT uid, domain, useruid, homework_uid, attachments, create_time, score FROM homework_submission WHERE useruid = ? AND homework_uid = ?", state["uid"], req.UID)
		if err := row.Scan(&sub.Uid, &sub.Domain, &sub.UserUid, &sub.HomeworkUid, &sub.Attachments, &sub.CreateTime, &sub.Score); err != nil {
			res.Success = false
			res.Error = "Not found"
			return
		}
	}
	res.Success = true
	res.HomeworkSubmission = sub
}

func getHomeworkSubmissions(ctx context.Context, req *senrenrpc.GetHomeworkSubmissionsRequest, state map[string]string, res *senrenrpc.GetHomeworkSubmissionsResponse) {
	sub := make([]*base.HomeworkSubmission, 0)
	row, err := db.DB.QueryContext(ctx, "SELECT homework_submission.uid, homework_submission.domain, homework_submission.useruid, homework_submission.homework_uid, homework_submission.attachments, homework_submission.create_time, homework_submission.score, user.nickname FROM homework_submission LEFT JOIN user ON homework_submission.domain = user.domain AND homework_submission.useruid = user.uid WHERE homework_submission.homework_uid = ? AND homework_submission.domain = ?", req.Filter, state["domain"])
	if err != nil {
		res.Success = false
		res.Error = err.Error()
	}
	for row.Next() {
		t := &base.HomeworkSubmission{}
		tNick := ""
		row.Scan(&t.Uid, &t.Domain, &t.UserUid, &t.HomeworkUid, &t.Attachments, &t.CreateTime, &t.Score, &tNick)
		t.Attachments = tNick + "!!" + t.Attachments
		sub = append(sub, t)
	}
	res.Success = true
	res.HomeworkSubmissions = sub
}

func packHomeworkSubmissions(ctx context.Context, req *senrenrpc.PackHomeworkSubmissionsRequest, state map[string]string, res *senrenrpc.PackHomeworkSubmissionsResponse) {
	// TODO
	hw := &base.Homework{}
	qry := db.DB.QueryRowContext(ctx, "SELECT uid, domain, title, content, attachments, start_time, end_time FROM homework WHERE domain = ? AND uid = ?", req.Domain, req.UID)
	if err := qry.Scan(&hw.Uid, &hw.Domain, &hw.Title, &hw.Description, &hw.Attachments, &hw.StartTime, &hw.EndTime); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	packTask := &base.HomeworkArchiveTask{}
	packTask.Type = "HomeworkArchiveTask"
	packTask.OutputFileName = fmt.Sprintf("%s.zip", strings.Map(func(r rune) rune {
		if r <= 32 {
			return '_'
		} else if r > 127 && r <= 255 {
			return '_'
		} else {
			return r
		}
	}, hw.Title))

	row, err := db.DB.QueryContext(ctx, "SELECT homework_submission.uid, homework_submission.domain, homework_submission.useruid, homework_submission.homework_uid, homework_submission.attachments, homework_submission.create_time, user.nickname FROM homework_submission LEFT JOIN user ON homework_submission.domain = user.domain AND homework_submission.useruid = user.guid WHERE homework_submission.homework_uid = ? AND homework_submission.domain = ?", req.UID, state["domain"])
	if err != nil {
		res.Success = false
		res.Error = err.Error()
	}
	packTask.Desc = &base.HomeworkArchiveDescriptor{
		Type:    "root",
		Content: make([]*base.HomeworkArchiveDescriptor, 0),
	}
	for row.Next() {
		t := &base.HomeworkSubmission{}
		row.Scan(&t.Uid, &t.Domain, &t.UserUid, &t.HomeworkUid, &t.Attachments, &t.CreateTime, &t.Nick)
		t2 := &base.HomeworkArchiveDescriptor{
			Type:    "archive",
			Name:    t.Nick,
			Content: make([]*base.HomeworkArchiveDescriptor, 0),
		}
		submitted := strings.Split(t.Attachments, ";")
		for _, v := range submitted {
			if v == "" {
				continue
			}

			v2 := strings.Split(v, ",")
			if len(v2) < 5 {
				continue
			}

			t2.Content = append(t2.Content, &base.HomeworkArchiveDescriptor{
				Type:   "file",
				Name:   v2[0] + v2[4],
				Source: "upload/" + v2[2][13:] + "/" + v2[2][:13] + "/" + v2[0] + v2[4],
			})
			logrus.Infof("%s", "upload/"+v2[2][13:]+"/"+v2[2][:13]+"/"+v2[0]+v2[4])
		}
		packTask.Desc.Content = append(packTask.Desc.Content, t2)
	}
	res.UID = util.GenUid()

	packTaskString, _ := json.Marshal(packTask)

	if _, err := db.DB.ExecContext(ctx, "INSERT INTO `task` (uid, create_time, state, creator, `desc`) VALUES (?, ?, ?, ?, ?)", res.UID, time.Now(), "PENDING", state["guid"], string(packTaskString)); err != nil {
		res.Success = false
		res.Error = err.Error()
	}

	publicchan.ChanTask <- res.UID

	res.Success = true
}

func createHomework(ctx context.Context, req *senrenrpc.CreateHomeworkRequest, state map[string]string, res *senrenrpc.CreateHomeworkResponse) {
	dbExec := "UPDATE homework SET title = ? , content = ? , attachments = ?, start_time = ?, end_time = ? WHERE uid = ? AND domain = ?"

	// TODO: check domain permission

	tDomain := senrenrpc.Domain(req.Homework.Domain)
	tDomain.ConvertDomain()
	req.Homework.Domain = string(tDomain)

	if req.Homework.Uid == "" || req.Homework.Uid == noUID {
		req.Homework.Uid = util.GenUid()

		dbExec = "INSERT INTO homework (title, content, attachments, start_time, end_time, uid, domain) VALUES (?, ?, ?, ?, ?, ?, ?)"
	}

	if _, err := db.DB.Exec(dbExec, req.Homework.Title, req.Homework.Description, req.Homework.Attachments, req.Homework.StartTime, req.Homework.EndTime, req.Homework.Uid, req.Homework.Domain); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	res.Domain = senrenrpc.Domain(req.Homework.Domain)
	res.UID = req.Homework.Uid
	res.Success = true
}

func getHomeworkSubmissionKey(ctx context.Context, req *senrenrpc.GetHomeworkSubmissionKeyRequest, state map[string]string, res *senrenrpc.GetHomeworkSubmissionKeyResponse) {
	sub := &base.HomeworkSubmission{}
	filterStr := strings.Split(req.Filter, ",")
	filterStr = append(filterStr, "", "")
	if filterStr[1] != "" && (state["grole"] == "ADMIN" || state["grole"] == "ROOT") {
		row := db.DB.QueryRowContext(ctx, "SELECT uid, domain, useruid, homework_uid, attachments, create_time FROM homework_submission WHERE uid = ? AND homework_uid = ?", filterStr[1], req.UID)
		if err := row.Scan(&sub.Uid, &sub.Domain, &sub.UserUid, &sub.HomeworkUid, &sub.Attachments, &sub.CreateTime); err != nil {
			res.Success = false
			res.Error = "Not found"
			return
		}
	} else {
		row := db.DB.QueryRowContext(ctx, "SELECT uid, domain, useruid, homework_uid, attachments, create_time FROM homework_submission WHERE useruid = ? AND homework_uid = ?", state["uid"], req.UID)
		if err := row.Scan(&sub.Uid, &sub.Domain, &sub.UserUid, &sub.HomeworkUid, &sub.Attachments, &sub.CreateTime); err != nil {
			res.Success = false
			res.Error = "Not found"
			return
		}
	}

	submitted := strings.Split(sub.Attachments, ";")
	exists := false
	for _, v := range submitted {
		if v == "" {
			continue
		}

		v2 := strings.Split(v, ",")
		if len(v2) < 5 {
			continue
		}

		if v2[2] == filterStr[0] {
			exists = true
			break
		}
	}

	if !exists {
		res.Success = false
		res.Error = "Not exists"
		return
	}

	res.Success = true
	res.UID = util.SignSession(filterStr[0])
}

type uploadHomeworkResponse senrenrpc.CreateHomeworkSubmissionResponse

func doHomeworkUpload(ctx context.Context, r *http.Request, req *senrenrpc.CreateHomeworkSubmissionRequest, res *uploadHomeworkResponse) {
	uhomework := r.Header.Get("UPLOAD_HOMEWORK")

	udomain := senrenrpc.Domain(r.Header.Get("UPLOAD_DOMAIN"))

	res.UID = util.GenUid()
	res.Domain = udomain
	res.Domain.ConvertDomain()

	if res.Domain == "0000000000000000" {
		res.Error = "Cannot submit homework to domain WOJ"
		res.Success = false
		return
	}

	state := make(map[string]string)
	checkLogin(ctx, req, res.Domain.GetDomain(), state)

	//-----//
	if state["role"] == "NONE" || state["role"] == "" {
		res.Error = "You should join the domain first"
		res.Success = false
		return
	}

	//-----//
	homework := base.Homework{}
	homeworkSub := base.HomeworkSubmission{}

	row := db.DB.QueryRowContext(ctx, "SELECT uid, domain, attachments, start_time, end_time FROM homework WHERE domain = ? AND uid = ?", res.Domain, uhomework)
	if err := row.Scan(&homework.Uid, &homework.Domain, &homework.Attachments, &homework.StartTime, &homework.EndTime); err != nil {
		res.Error = "Homework not found (domain = " + string(res.Domain) + ", uid = " + uhomework + ")"
		res.Success = false
		return
	}

	current := time.Now()
	if current.Before(homework.StartTime) || current.After(homework.EndTime) {
		res.Error = "Illegal submission time, (maybe overdue?)"
		res.Success = false
		return
	}

	doCreate := false
	row2 := db.DB.QueryRowContext(ctx, "SELECT uid, domain, useruid, homework_uid, attachments, create_time FROM homework_submission WHERE useruid = ? AND homework_uid = ?", state["uid"], uhomework)
	if err := row2.Scan(&homeworkSub.Uid, &homeworkSub.Domain, &homeworkSub.UserUid, &homeworkSub.HomeworkUid, &homeworkSub.Attachments, &homeworkSub.CreateTime); err != nil {
		// not submitted before
		doCreate = true
		homeworkSub.Uid = util.GenUid()
		homeworkSub.Domain = homework.Domain
		homeworkSub.UserUid = state["uid"]
		homeworkSub.Attachments = ""
		homeworkSub.HomeworkUid = homework.Uid
	}
	homeworkSub.CreateTime = time.Now()

	logrus.Infof("SUBMIT BY %s --> %s.%s", homeworkSub.UserUid, homeworkSub.Domain, homeworkSub.HomeworkUid)

	r.ParseMultipartForm(16 << 20)

	itemname := r.FormValue("name")
	if itemname == "" {
		res.Error = "You should specific one homework item to submit."
		res.Success = false
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		res.Error = err.Error()
		res.Success = false
		return
	}

	defer file.Close()

	fileExt := ""
	//-----//
	isValidSubmission := false
	validSlot := strings.Split(homework.Attachments, ";")
	for _, v := range validSlot {
		v2 := strings.Split(v, ",")
		if len(v2) < 3 {
			continue
		}
		if v2[0] != itemname {
			continue
		}
		v2s := strings.Split(v2[2], "|")
		for _, vv := range v2s {
			if vv == "*" {
				isValidSubmission = true
				fileExt = ""
				break
			}
			if strings.HasSuffix(handler.Filename, "."+vv) {
				isValidSubmission = true
				fileExt = "." + vv
				break
			}
		}
	}

	if !isValidSubmission {
		res.Error = "No slot found. (Maybe file type mismatch.)"
		res.Success = false
		return
	}

	newSubmitted := make([]string, 0)
	submitted := strings.Split(homeworkSub.Attachments, ";")
	for _, v := range submitted {
		if v == "" {
			continue
		}

		v2 := strings.Split(v, ",")
		if len(v2) < 5 {
			continue
		}

		if v2[0] == itemname {
			continue
		}

		newSubmitted = append(newSubmitted, v)
	}

	os.MkdirAll("upload/"+res.UID[13:]+"/"+res.UID[:13], os.ModePerm)
	wrfile, err := os.OpenFile("upload/"+res.UID[13:]+"/"+res.UID[:13]+"/"+itemname+fileExt, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		res.Error = err.Error()
		res.Success = false
		return
	}
	defer wrfile.Close()

	if _, err := io.Copy(wrfile, file); err != nil {
		res.Error = err.Error()
		res.Success = false
		return
	}

	newSubmitted = append(newSubmitted, strings.Join([]string{itemname, base64.StdEncoding.EncodeToString([]byte(handler.Filename)), res.UID, fmt.Sprint(handler.Size), fileExt}, ","))
	homeworkSub.Attachments = strings.Join(newSubmitted, ";")

	queryStr := "UPDATE homework_submission SET attachments = ?, score = ?, create_time = ? WHERE uid = ? AND domain = ? AND useruid = ? AND homework_uid = ?"
	if doCreate {
		queryStr = "INSERT INTO homework_submission (attachments, score, create_time, uid, domain, useruid, homework_uid) VALUES (?, ?, ?, ?, ?, ?)"
	}

	if _, err := db.DB.ExecContext(ctx, queryStr, homeworkSub.Attachments, -1, homeworkSub.CreateTime, homeworkSub.Uid, homeworkSub.Domain, homeworkSub.UserUid, homeworkSub.HomeworkUid); err != nil {
		res.Error = err.Error()
		res.Success = false
		return
	}
	res.Success = true
	res.UID = homeworkSub.Uid
	res.Domain = senrenrpc.Domain(homeworkSub.Domain)
}

func downloadHomework(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	token := params["token"]
	validUid := util.CheckSessionTime(token, 30*time.Second)
	if validUid == "" {
		w.WriteHeader(404)
		return
	}
	logrus.Infof("Download: %s / %s", validUid, params["filename"])
	f, err := os.Open(filepath.Join("upload", validUid[13:], validUid[:13], params["filename"]))
	if err != nil {
		w.WriteHeader(404)
	}
	io.Copy(w, f)
}
