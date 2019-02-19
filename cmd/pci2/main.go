package main

import (
	"net/http"
	"os"

	_ "github.com/erjiaqing/senren2/pkg/pciend"
	"github.com/sirupsen/logrus"
)

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8079"
	}
	logrus.Fatal(http.ListenAndServe(listenAddr, nil))
}
