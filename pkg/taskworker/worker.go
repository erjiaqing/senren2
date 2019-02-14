package taskworker

import (
	"encoding/json"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/publicchan"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/sirupsen/logrus"
)

func Work() {
	for {
		taskuid, ok := <-publicchan.ChanTask
		if !ok {
			break
		}
		var desc string
		row := db.DB.QueryRow("SELECT `desc` FROM task WHERE uid = ?", taskuid)
		if err := row.Scan(&desc); err != nil {
			logrus.Errorf("Failed to get task: %v", err)
			continue
		}

		logrus.Infof("INCOMING TASK: %s", taskuid)

		task := &base.ServerTask{}
		json.Unmarshal([]byte(desc), task)

		result := "UNDEF"

		switch task.Type {
		case "HomeworkArchiveTask":
			t2 := &base.HomeworkArchiveTask{}
			json.Unmarshal([]byte(desc), t2)
			result = "SUCCESS"
			if err := homeworkArchiveTask(taskuid, t2); err != nil {
				result = "ERROR"
			}
		}

		db.DB.Exec("UPDATE task SET state = ? WHERE uid = ?", result, taskuid)
		logrus.Infof("INCOMING TASK: %s [%s]", taskuid, result)
	}
	logrus.Errorf("Failed to get tasks")
}
