package persec

import (
	"log"
	"testing"
	"time"
)

//FakeConn Fake redis connection
type FakeConn struct{}

func (conn FakeConn) Close() error                                                    { return nil }
func (conn FakeConn) Err() error                                                      { return nil }
func (conn FakeConn) Send(commandName string, args ...interface{}) error              { return nil }
func (conn FakeConn) Flush() error                                                    { return nil }
func (conn FakeConn) Receive() (interface{}, error)                                   { return nil, nil }
func (conn FakeConn) Do(commandName string, args ...interface{}) (interface{}, error) { return nil, nil }

type FakeConnGetter struct {
	FakeConn
}

func (conn FakeConnGetter) Do(commandName string, args ...interface{}) (interface{}, error) {
	return int64(100), nil
}

type FakeConnIncr struct {
	FakeConn
}

func (conn FakeConnIncr) Do(commandName string, args ...interface{}) (interface{}, error) {
	if commandName == "INCR" {
		// couldn't find a better way to test this for now...
		log.Panic("INCR called !")
	}
	return nil, nil
}

func TestGetTotalRequests(t *testing.T) {
	conn := FakeConnGetter{}
	timestampStart := time.Now().Unix()
	timestampEnd := timestampStart + 10
	requestKeyWithoutTimestamp := "FOO"

	total, duration := GetTotalRequests(
		conn,
		requestKeyWithoutTimestamp,
		timestampStart,
		timestampEnd,
	)

	correctTotal := 1000
	correctDuration := int64(timestampEnd - timestampStart)

	if total != correctTotal {
		t.Errorf("Wrong total, got %d but it should be %d", total, correctTotal)
	}

	if int64(duration.Seconds()) != correctDuration {
		t.Errorf("Wrong duration, got %d but it should be %d", duration, correctDuration)
	}
}

func TestFailGetTotalRequests(t *testing.T) {
	conn := FakeConn{}
	timestampStart := time.Now().Unix()
	timestampEnd := timestampStart - 10
	requestKeyWithoutTimestamp := "FOO"

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic but it was supposed to do so.")
		}
	}()

	GetTotalRequests(
		conn,
		requestKeyWithoutTimestamp,
		timestampStart,
		timestampEnd,
	)
}

func TestIncrementRequestKey(t *testing.T) {
	conn := FakeConnIncr{}
	requestKeyTimestamp := GetRequestKeyTimestamp(
		"FOO",
		time.Now())

	defer func() {
		if r := recover(); r != "INCR called !" {
			t.Errorf("Wrong panic on TestIncrementRequestKey()")
		}
	}()

	IncrementRequestKey(conn, requestKeyTimestamp)
}
