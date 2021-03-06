package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

var sigSign = "417MD2y3nUNI7e5r1CZ8x6LBLzwkMk4y8FHmrC0srJP5YHXPXUXrcYNj4Bfr417MD2y3nUNI7e5r1CZ8FLZ2vk3iTbUaUQ9v0wBk"

func init() {
	sigSign2 := os.Getenv("SIG_SIGN")
	if sigSign2 != "" {
		sigSign = sigSign2
	}
}

func Sign(va ...string) string {
	raw := bytes.Buffer{}
	for _, v := range va {
		raw.WriteString(base64.StdEncoding.EncodeToString([]byte(v)))
		raw.WriteString("|")
	}
	raw.WriteString(sigSign)
	return fmt.Sprintf("%x", sha1.Sum(raw.Bytes()))
}

func CheckSign(sig string, va ...string) bool {
	raw := bytes.Buffer{}
	for _, v := range va {
		raw.WriteString(base64.StdEncoding.EncodeToString([]byte(v)))
		raw.WriteString("|")
	}
	raw.WriteString(sigSign)
	return sig == fmt.Sprintf("%x", sha1.Sum(raw.Bytes()))
}

func SignSession(uid string) string {
	return SignSessionDomain(uid, "")
}

func SignSessionDomain(uid string, domain string) string {
	current := time.Now()
	ruid := base64.URLEncoding.EncodeToString([]byte(uid))
	ruid = strings.TrimRight(ruid, "=")
	sigSignSrc := fmt.Sprintf("%s_%s_%16x_%s", ruid, domain, current.UnixNano(), sigSign)
	sigSignRes := fmt.Sprintf("%x", sha1.Sum([]byte(sigSignSrc)))
	return fmt.Sprintf("%s_%s_%16x_%s", ruid, domain, current.UnixNano(), sigSignRes)
}

func CheckSession(sid string) string {
	return CheckSessionTime(sid, 14*24*time.Hour)
}

func CheckSessionTimeDomain(sid string, limit time.Duration, domain string) string {
	parts := strings.Split(sid, "_")
	if len(parts) != 4 {
		return ""
	}
	var sigtime int64
	fmt.Sscanf(parts[2], "%x", &sigtime)
	current := time.Now()
	if sigtime < current.UnixNano()-limit.Nanoseconds() {
		return ""
	}
	sigSignSrc := fmt.Sprintf("%s_%s_%s_%s", parts[0], parts[1], parts[2], sigSign)
	sigSignRes := fmt.Sprintf("%x", sha1.Sum([]byte(sigSignSrc)))
	if sigSignRes != parts[3] {
		return ""
	}
	for len(parts[0])%4 != 0 {
		parts[0] += "="
	}
	uid, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return ""
	} else if parts[1] != domain {
		return ""
	}
	return string(uid)
}

func CheckSessionTime(sid string, limit time.Duration) string {
	return CheckSessionTimeDomain(sid, limit, "")
}
