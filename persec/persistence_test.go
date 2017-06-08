package persec

import (
	"strconv"
	"testing"
	"time"
)

func TestFailGetRedisConn(t *testing.T) {
	_, err := GetRedisConn("foo")

	if err == nil {
		t.Error("This should be an error.")
	}
}

func TestGetRequestKeyTimestamp(t *testing.T) {
	now := time.Now()
	result := GetRequestKeyTimestamp("FOO", now)

	if result != "FOO_"+strconv.FormatInt(now.Unix(), 10) {
		t.Error("Invalid request key.")
	}
}
