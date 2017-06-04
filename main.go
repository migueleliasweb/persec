package main

import (
	"fmt"
	"time"
)

func main() {
	conn := GetRedisConn("redis://127.0.0.1:6379")

	go func() {
		c := time.Tick(1 * time.Second)
		for now := range c {
			key := GetRequestKeyTimestamp("BARRRR")
			total, duration := GetTotalRequests(
				conn,
				key,
				now.Unix(),
				now.Unix()+int64(10),
			)
			fmt.Println("AVG: ", GetAvgRequests(
				total,
				duration,
				time.Duration(5),
			))
		}
	}()

	for {
		key := GetRequestKeyTimestamp("BARRRR")
		IncrementRequestKey(conn, key)
	}
}
