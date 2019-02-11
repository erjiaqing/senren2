package util

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

const sigSign = "417MD2y3nUNI7e5r1CZ8x6LBLzwkMk4y8FHmrC0srJP5YHXPXUXrcYNj4Bfr417MD2y3nUNI7e5r1CZ8FLZ2vk3iTbUaUQ9v0wBk"

func SignSession(uid string) string {
	current := time.Now()
	sigSignSrc := fmt.Sprintf("%s:%16x:%s", uid, current.UnixNano(), sigSign)
	sigSignRes := fmt.Sprintf("%x", sha1.Sum([]byte(sigSignSrc)))
	return fmt.Sprintf("%s:%16x:%s", uid, current.UnixNano(), sigSignRes)
}

func CheckSession(sid string) string {
	parts := strings.Split(sid, ":")
	if len(parts) != 3 {
		return ""
	}
	var sigtime int64
	fmt.Sscanf(parts[1], "%x", &sigtime)
	current := time.Now()
	if sigtime < current.UnixNano()-14*24*3600*1000000000 { // 14 days
		return ""
	}
	sigSignSrc := fmt.Sprintf("%s:%s:%s", parts[0], parts[1], sigSign)
	sigSignRes := fmt.Sprintf("%x", sha1.Sum([]byte(sigSignSrc)))
	if sigSignRes != parts[2] {
		return ""
	}
	return parts[0]
}
