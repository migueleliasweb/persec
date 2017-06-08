package persec

import (
	"testing"
	"time"
)

func TestFailGetRedisConn(t *testing.T) {
	_, err := GetRedisConn("foo")

	if err == nil {
		t.Error("This should be an error.")
	}
}

func TestSimpleGetAvgRequests(t *testing.T) {
	avg := GetAvgRequests(
		1000,
		time.Second*100,
		time.Second)

	if avg != 10 {
		t.Errorf("Should be 10 but was %d", avg)
	}
}

func TestSimple2GetAvgRequests(t *testing.T) {
	avg := GetAvgRequests(
		1000,
		time.Second*100,
		time.Second*1000)

	if avg != 10000 {
		t.Error("Should be 10000.")
	}
}

func TestZeroDivGetAvgRequests(t *testing.T) {
	avg := GetAvgRequests(
		1000,
		time.Second*100,
		0)

	if avg != 10 {
		t.Errorf("Should be 10 but was %d", avg)
	}
}
