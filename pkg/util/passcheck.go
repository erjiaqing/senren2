package util

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authResult struct {
	Errno int    `json:"errno"`
	Error string `json:"error"`
	Data  struct {
		Result string `json:"result"`
	} `json:"data"`
}

func AuthOlive(username, password string) (bool, error) {
	resp, err := http.PostForm("http://acm.whu.edu.cn/olive/login",
		url.Values{"user": {username}, "pass": {password}})
	if err != nil {
		return false, errors.New("network error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, errors.New("network error")
	}
	authRes := &authResult{}
	json.Unmarshal(body, authRes)
	if authRes.Data.Result != "success" {
		return false, errors.New("username and password mismatch")
	}
	return true, nil
}

func CheckPass(pass, upass string) bool {
	passField := strings.Split(pass, ":")
	switch passField[0] {
	case "bcrypt":
		if err := bcrypt.CompareHashAndPassword([]byte(passField[1]), []byte(upass)); err != nil {
			return false
		}
		return true
	case "plain":
		if passField[1] != base64.StdEncoding.EncodeToString([]byte(upass)) {
			sleepTime := rand.Intn(150) + 150
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
			// avoid time attack
			return false
		}
		return true
	default:
		return false
	}
}
