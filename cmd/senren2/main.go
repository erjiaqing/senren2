package main

import (
	"net/http"
	"os"

	_ "github.com/erjiaqing/senren2/pkg/endpoint"
	"github.com/erjiaqing/senren2/pkg/taskworker"
	"github.com/sirupsen/logrus"
)

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	go taskworker.Work()
	logrus.Fatal(http.ListenAndServe(listenAddr, nil))
}
