package endpoint

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

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
	row := db.DB.QueryRowContext(ctx, "SELECT uid, domain, useruid, homework_uid, attachments, create_time FROM homework_submission WHERE useruid = ? AND homework_uid = ?", state["guid"], req.UID)
	if err := row.Scan(&sub.Uid, &sub.Domain, &sub.UserUid, &sub.HomeworkUid, &sub.Attachments, &sub.CreateTime); err != nil {
		res.Success = false
		res.Error = "Not found"
		return
	}
	res.Success = true
	res.HomeworkSubmission = sub
}

func getHomeworkSubmissions(ctx context.Context, req *senrenrpc.GetHomeworkSubmissionsRequest, state map[string]string, res *senrenrpc.GetHomeworkSubmissionsResponse) {
	sub := make([]*base.HomeworkSubmission, 0)
	row, err := db.DB.QueryContext(ctx, "SELECT homework_submission.uid, homework_submission.domain, homework_submission.useruid, homework_submission.homework_uid, homework_submission.attachments, homework_submission.create_time, user.nickname FROM homework_submission LEFT JOIN user ON homework_submission.domain = user.domain AND homework_submission.useruid = user.guid WHERE homework_submission.homework_uid = ? AND homework_submission.domain = ?", req.Filter, state["domain"])
	if err != nil {
		res.Success = false
		res.Error = err.Error()
	}
	for row.Next() {
		t := &base.HomeworkSubmission{}
		tNick := ""
		row.Scan(&t.Uid, &t.Domain, &t.UserUid, &t.HomeworkUid, &t.Attachments, &t.CreateTime, &tNick)
		t.Attachments = tNick + "!!" + t.Attachments
		sub = append(sub, t)
	}
	res.Success = true
	res.HomeworkSubmissions = sub
}

func packHomeworkSubmissions(ctx context.Context, req *senrenrpc.PackHomeworkSubmissionsRequest, state map[string]string, res *senrenrpc.PackHomeworkSubmissionsResponse) {
	// TODO
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

func createHomeworkSubmission(ctx context.Context, req *senrenrpc.CreateHomeworkSubmissionRequest, state map[string]string, res *senrenrpc.CreateHomeworkSubmissionResponse) {

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
	row2 := db.DB.QueryRowContext(ctx, "SELECT uid, domain, useruid, homework_uid, attachments, create_time FROM homework_submission WHERE useruid = ? AND homework_uid = ?", state["guid"], uhomework)
	if err := row2.Scan(&homeworkSub.Uid, &homeworkSub.Domain, &homeworkSub.UserUid, &homeworkSub.HomeworkUid, &homeworkSub.Attachments, &homeworkSub.CreateTime); err != nil {
		// not submitted before
		doCreate = true
		homeworkSub.Uid = util.GenUid()
		homeworkSub.Domain = homework.Domain
		homeworkSub.UserUid = state["guid"]
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

	queryStr := "UPDATE homework_submission SET attachments = ?, create_time = ? WHERE uid = ? AND domain = ? AND useruid = ? AND homework_uid = ?"
	if doCreate {
		queryStr = "INSERT INTO homework_submission (attachments, create_time, uid, domain, useruid, homework_uid) VALUES (?, ?, ?, ?, ?, ?)"
	}

	if _, err := db.DB.ExecContext(ctx, queryStr, homeworkSub.Attachments, homeworkSub.CreateTime, homeworkSub.Uid, homeworkSub.Domain, homeworkSub.UserUid, homeworkSub.HomeworkUid); err != nil {
		res.Error = err.Error()
		res.Success = false
		return
	}
	res.Success = true
	res.UID = homeworkSub.Uid
	res.Domain = senrenrpc.Domain(homeworkSub.Domain)
}
