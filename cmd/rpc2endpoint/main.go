package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("Failed to get working directory!")
	}

	logrus.Infof("Source code path: %s", path)

	data, err := ioutil.ReadFile("api.txt")

	if err != nil {
		logrus.Fatalf("Failed to read api.txt: %v", err)
	}

	ds := string(data)

	lines := strings.Split(ds, "\n")

	for _, api := range lines {
		spec := strings.Split(api, " ")
		if len(spec) != 3 {
			continue
		}
		fmt.Printf(`case "%s":
			req = &%s{}
			res = &%s{}
			`, spec[0], spec[1][1:], spec[2][1:])
	}

	for _, api := range lines {
		spec := strings.Split(api, " ")
		if len(spec) != 3 {
			continue
		}
		fmt.Printf(`case "%s":
			%s(req.(%s), state, res.(%s))
			`, spec[0], spec[0], spec[1], spec[2])
	}
}
