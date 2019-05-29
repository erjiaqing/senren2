package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/util"
	"github.com/sirupsen/logrus"
)

func judge(task *base.PCITaskItem, desc *base.PCIJudgeTaskDesc) string {
	// TODO: save file in temporary folder, and then run PCI Judger 1.5
	// Clone problem first
	prob := checkProblem(task.Problem)
	probpath := cloneProblemVersion(task.Problem, prob.CurrentVersion)

	cloneLock.RLock()
	defer cloneLock.RUnlock()
	currpath, err := os.Getwd()

	if err != nil {
		return `{"error": "failed to judge"}`
	}

	defer os.Chdir(currpath)

	tempdir := filepath.Join(currpath, "temp", util.GenUid())
	rtempdir := filepath.Join(currpath, "temp", util.GenUid())

	//defer os.RemoveAll(tempdir)
	//defer os.RemoveAll(rtempdir)

	os.MkdirAll(tempdir, os.ModePerm)
	os.MkdirAll(rtempdir, os.ModePerm)
	if err := os.Chdir(tempdir); err != nil {
		return `{"error": "failed to judge"}`
	}

	if err := ioutil.WriteFile(filepath.Join(tempdir, "code"), []byte(desc.Code), os.ModePerm); err != nil {
		return `{"error": "failed to judge"}`
	}

	// "docker", "run", "--privileged", "--mount", "type=bind,source={},target=/problem,readonly".format(problem), "    --mount", "type=bind,source={},target=/code/code,readonly".format(code), "--mount", "type=bind,source={},target=/fj_tmp".format(tmpdir)    , "fjudger2", "--language", language, "--tempdir", "/fj_tmp", "--docker=1"
	logrus.Infof("Judge")
	logrus.Info(strings.Join([]string{"docker", "run", "--privileged", "--mount", fmt.Sprintf("type=bind,source=%s,target=/problem,readonly", probpath), "--mount", fmt.Sprintf("type=bind,source=%s,target=/code/code,readonly", filepath.Join(tempdir, "code")), "--mount", fmt.Sprintf("type=bind,source=%s,target=/fj_tmp", rtempdir), "--", "fj2-judger", "--language", desc.Lang, "--tempdir", "/fj_tmp", "--docker=1"}, " "))
	cmd := exec.Command("docker", "run", "--privileged", "--mount", fmt.Sprintf("type=bind,source=%s,target=/problem,readonly", probpath), "--mount", fmt.Sprintf("type=bind,source=%s,target=/code/code,readonly", filepath.Join(tempdir, "code")), "--mount", fmt.Sprintf("type=bind,source=%s,target=/fj_tmp", rtempdir), "--", "fj2-judger", "--language", desc.Lang, "--tempdir", "/fj_tmp", "--docker=1")
	cmdOutBuff := &bytes.Buffer{}
	cmdErrBuff := &bytes.Buffer{}
	cmd.Stdout = cmdOutBuff
	cmd.Stderr = cmdErrBuff

	if err := cmd.Run(); err != nil {
		logrus.Error(err)
		logrus.Error(cmdErrBuff.String())
	}

	ioutil.WriteFile(filepath.Join(tempdir, "fj.stderr"), cmdErrBuff.Bytes(), os.ModePerm)

	return cmdOutBuff.String()
}
