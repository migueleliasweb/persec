package main

import (
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

//RedisConn Returns the connection to Redis
func RedisConn(redisURL string) *redis.Conn {
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
func GetTotalRequests(conn *redis.Conn, key, timestampStart int64, timestampEnd int64) int {
	total := 0
	for timestampStart < timestampEnd {

	}

	return total
}

func main() {
	conn := RedisConn("redis://172.17.0.2:6379")
	for {
		key := GetRequestKeyTimestamp("BARRRR")
		IncrementRequestKey(conn, key)
	}
}
