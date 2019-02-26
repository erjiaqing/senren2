package util

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

const sigSign = "417MD2y3nUNI7e5r1CZ8x6LBLzwkMk4y8FHmrC0srJP5YHXPXUXrcYNj4Bfr417MD2y3nUNI7e5r1CZ8FLZ2vk3iTbUaUQ9v0wBk"

func SignSession(uid string) string {
	current := time.Now()
	ruid := base64.URLEncoding.EncodeToString([]byte(uid))
	sigSignSrc := fmt.Sprintf("%s:%16x:%s", ruid, current.UnixNano(), sigSign)
	sigSignRes := fmt.Sprintf("%x", sha1.Sum([]byte(sigSignSrc)))
	return fmt.Sprintf("%s:%16x:%s", ruid, current.UnixNano(), sigSignRes)
}

func CheckSession(sid string) string {
	return CheckSessionTime(sid, 14*24*time.Hour)
}

func CheckSessionTime(sid string, limit time.Duration) string {
	parts := strings.Split(sid, ":")
	if len(parts) != 3 {
		return ""
	}
	var sigtime int64
	fmt.Sscanf(parts[1], "%x", &sigtime)
	current := time.Now()
	if sigtime < current.UnixNano()-limit.Nanoseconds() {
		return ""
	}
	sigSignSrc := fmt.Sprintf("%s:%s:%s", parts[0], parts[1], sigSign)
	sigSignRes := fmt.Sprintf("%x", sha1.Sum([]byte(sigSignSrc)))
	if sigSignRes != parts[2] {
		return ""
	}
	uid, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return ""
	}
	return string(uid)
}
