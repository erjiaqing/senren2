package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	gitUserName = ""
	gitUserPass = ""
	gitServer   = ""

	pciServer  = ""
	pciSession = ""
)

func init() {
	gitUserName = os.Getenv("PCI_GIT_USER")
	if gitUserName == "" {
		gitUserName = "pci"
	}
	gitUserPass = os.Getenv("PCI_GIT_PASS")
	if gitUserPass == "" {
		gitUserPass = "pcipcipci"
	}
	gitServer = os.Getenv("PCI_GIT_SERV")
	if gitServer == "" {
		gitServer = "http://127.0.0.1:3000"
	}
	pciServer = os.Getenv("PCI_SERV")
	if pciServer == "" {
		pciServer = "http://127.0.0.1:8079"
	}
	pciSession = os.Getenv("PCI_SESSION")
	if pciSession == "" || len(pciSession) != 32 {
		logrus.Fatal("PCI Session is requires")
	}

	logrus.Infof("Git User: %s", gitUserName)
	logrus.Infof("Git Pass: %s", gitUserPass[0:1]+("********")+gitUserPass[len(gitUserPass)-1:len(gitUserPass)])
	logrus.Infof("Git Serv: %s", gitServer)
	logrus.Infof("PCI Serv: %s", pciServer)
	logrus.Infof("PCI Sess: %s", pciSession[0:8]+("****************")+pciSession[24:32])
}

// TODO: websocket long connection

func main() {
	for {
		t := getGeneralTask()
		if t == nil {
			time.Sleep(5 * time.Second)
			continue
		}
		logrus.Infof("Get task!")
		d := processGeneralTask(t.Desc)
		doGeneralTask(t, d)
		updateTask(t)
	}
}
