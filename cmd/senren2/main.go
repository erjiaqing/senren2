package main

import (
	"net/http"
	"os"

	_ "github.com/erjiaqing/senren2/pkg/endpoint"
	"github.com/sirupsen/logrus"
)

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	logrus.Fatal(http.ListenAndServe(listenAddr, nil))
}
