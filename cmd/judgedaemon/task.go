package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/erjiaqing/senren2/pkg/httpreq"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/pcirpc"
	"github.com/sirupsen/logrus"
)

func getGeneralTask() *base.PCITaskItem {
	req := &pcirpc.GetPCITaskRequest{}
	req.Key = pciSession

	res, code, err := httpreq.POSTJson(fmt.Sprintf("%s/rpc/pci/getTask", pciServer), req)

	if err != nil {
		logrus.Errorf("Failed to do getTask from ProblemCI: %v", err)
		return nil
	}

	if code >= 300 {
		logrus.Errorf("Unexpected return code: 200 expected, %d get", code)
		return nil
	}

	ret := &pcirpc.GetPCITaskResponse{}
	json.Unmarshal(res, ret)

	if !ret.Success {
		return nil
	}

	return ret.Task
}

func processGeneralTask(taskjson string) interface{} {
	if taskjson == "" {
		return nil
	}

	t := &base.PCITask{}
	var r2 interface{}

	json.Unmarshal([]byte(taskjson), t)

	switch t.Type {
	case "judge":
		r2 = &base.PCIJudgeTaskDesc{}
	case "build":
		r2 = &base.PCIBuildTaskDesc{}
	default:
		return nil
	}

	json.Unmarshal([]byte(taskjson), r2)

	return r2
}

func doGeneralTask(task *base.PCITaskItem, desc interface{}) {
	switch desc.(type) {
	case *base.PCIJudgeTaskDesc:
		task.Result = judge(task, desc.(*base.PCIJudgeTaskDesc))
		task.Status = "FINISHED"
	case *base.PCIBuildTaskDesc:
		task.Status = "NOT_IMPL"
		task.Result = ""
	default:
		task.Status = "FAIL"
		task.Result = `{"error": "unknown task"}`
	}
}

func updateTask(task *base.PCITaskItem) {
	req := &pcirpc.UpdatePCITaskRequest{}
	ret := &pcirpc.UpdatePCITaskResponse{}
	req.Key = pciSession
	req.Task = task

	for i := uint(0); i < 5; i++ {
		res, code, err := httpreq.POSTJson(fmt.Sprintf("%s/rpc/pci/updateTask", pciServer), req)

		if err != nil {
			logrus.Warningf("Failed to do getTask from ProblemCI: %v", err)
		} else if code >= 300 {
			logrus.Warningf("Unexpected return code: 200 expected, %d get", code)
		}

		json.Unmarshal(res, ret)

		if !ret.Success {
			time.Sleep((1 << i) * time.Second / 2)
		}

		return
	}

	logrus.Errorf("Failed to update task!")
}
