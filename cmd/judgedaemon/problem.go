package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/erjiaqing/senren2/pkg/httpreq"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/pcirpc"
	"github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
)

func checkProblem(uid int64) *base.PCIProblem {
	req := &pcirpc.GetProblemRequest{}
	req.Key = pciSession
	req.ProblemId = uid

	res, code, err := httpreq.POSTJson(fmt.Sprintf("%s/rpc/pci/getProblem", pciServer), req)

	if code >= 300 {
		logrus.Errorf("Failed to get problem: HTTP%d", code)
		return nil
	} else if err != nil {
		logrus.Errorf("Failed to get problem: %v", err.Error())
		return nil
	}

	ret := &pcirpc.GetProblemResponse{}
	if err := json.Unmarshal(res, ret); err != nil {
		logrus.Errorf("Failed to get problem: %v", err.Error())
		return nil
	}

	return ret.Problem
}

func checkProblemVersion(uid int64, expect string) bool {
	wd, err := os.Getwd()
	if err != nil {
		return false
	}
	r, err := git.PlainOpen(filepath.Join(wd, "problem", strconv.FormatInt(uid, 10)))
	if err != nil {
		return false
	}
	h, err := r.Head()
	if err != nil {
		return false
	}
	return h.Hash().String() == expect
}

func cloneProblemVersion(uid int64, expect string) string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	target := filepath.Join(wd, "problem", strconv.FormatInt(uid, 10))

	if checkProblemVersion(uid, expect) {
		return target
	}

	cloneurl, err := url.Parse(fmt.Sprintf("%s/%s/p%06d.git", gitServer, gitUserName, uid))
	if err != nil {
		return ""
	}

	cloneurl.User = url.UserPassword(gitUserName, gitUserPass)

	exec.Command("rm", "-rf", target).Run()

	logrus.Infof("Clone problem: %d", uid)
	clonecmd := exec.Command("git", "clone", "-q", "--", cloneurl.String(), target)
	terr := &bytes.Buffer{}
	clonecmd.Stderr = terr
	if err := clonecmd.Run(); err != nil {
		logrus.Errorf("Failed to clone problem: %v", err)
		logrus.Error(terr.String())
		return ""
	}

	cmd := exec.Command("git", "checkout", expect)
	cmd.Dir = target
	cmd.Run()

	// docker run --privileged --mount type=bind,source=/hofioeg,target=/problem --mount type=bind,source=/home/ejq/rrtmp,target=/fj_tmp fj2-builder --tempdir /fj_tmp
	logrus.Infof("Build problem: %d", uid)
	exec.Command("docker", "run", "--rm", "--privileged", "--mount", fmt.Sprintf("type=bind,source=%s,target=/problem", target), "--mount", fmt.Sprintf("type=bind,source=%s,target=/fj_tmp", filepath.Join(wd, "temp")), "fj2-builder", "--tempdir", "/fj_tmp").Run()
	return target
}
