package httpreq

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func POSTJsonAuth(url string, body interface{}, authname string, authpass string) ([]byte, int, error) {
	bodybyte, err := json.Marshal(body)

	if err != nil {
		return nil, 0, err
	}

	logrus.Infof("POST %s", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodybyte))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	if authname != "" {
		req.SetBasicAuth(authname, authpass)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)

	return respbody, resp.StatusCode, err
}

func POSTJson(url string, body interface{}) ([]byte, int, error) {
	return POSTJsonAuth(url, body, "", "")
}

func GETAuth(url string, authname string, authpass string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	if authname != "" {
		req.SetBasicAuth(authname, authpass)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)

	return respbody, resp.StatusCode, err
}

func GET(url string) ([]byte, int, error) {
	return GETAuth(url, "", "")
}
