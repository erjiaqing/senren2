package util

import (
	"math/rand"
	"time"
)

const rand_runes = "23456789abcdefghijkmnpqrstuvwxyz"

var sessionId = ""

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = rand_runes[rand.Intn(len(rand_runes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	sessionId = RandStringBytes(3)
}

func GenUid() string {
	ret := ""
	t := time.Now()
	upper10 := t.UnixNano() / 1048576
	for i := 0; i < 10; i++ {
		ret = string(rand_runes[upper10%32]) + ret
		upper10 /= 32
	}
	lower3 := RandStringBytes(3)
	return ret + sessionId + lower3
}
