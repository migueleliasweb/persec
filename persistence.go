package main

import (
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

//GetRedisConn Returns the connection to Redis
func GetRedisConn(redisURL string) *redis.Conn {
	c, err := redis.DialURL(redisURL)

	if err != nil {
		log.Fatalln(err)
	}

	return &c
}

//GetRequestKeyTimestamp Get the complete RequestKey
func GetRequestKeyTimestamp(requestKey string) string {
	return requestKey + "_" + strconv.FormatInt(time.Now().Unix(), 10)
}

//IncrementRequestKey Inscrements the amount of times a certain key was requested
func IncrementRequestKey(conn *redis.Conn, requestKeyTimestamp string) {
	(*conn).Do("INCR", requestKeyTimestamp)
}

//GetTotalRequests Returns the total requests for the give period
func GetTotalRequests(
	conn *redis.Conn,
	requestKeyTimestamp string,
	timestampStart int64,
	timestampEnd int64) (int, time.Duration) {
	total := 0
	seconds := 0
	for timestampEnd > timestampStart {
		requestsNum, cmdErr := redis.Int((*conn).Do("GET", requestKeyTimestamp))

		if cmdErr != nil {
			if cmdErr.Error() != redis.ErrNil.Error() {
				break
			} else {
				log.Panicf("Got error from Redis: %e", cmdErr)
			}
		}

		total += requestsNum
		seconds++
		timestampEnd--
	}

	//With this approach we can estimate better the throughput
	return total, time.Duration(seconds)
}

//GetAvgRequests Returns avg requests for the given duration
func GetAvgRequests(totalRequests int, totalRequestsDuration time.Duration, desiredDuration time.Duration) int64 {
	perSecondAvg := int64(totalRequests) / int64(totalRequestsDuration.Seconds())
	return perSecondAvg * int64(desiredDuration.Seconds())
}
