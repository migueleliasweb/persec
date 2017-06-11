package persec

import (
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

//GetRedisConn Returns the connection to Redis
func GetRedisConn(redisURL string) (*redis.Conn, error) {
	c, err := redis.DialURL(redisURL)

	return &c, err
}

//GetRequestKeyTimestamp Get the complete RequestKey
func GetRequestKeyTimestamp(requestKey string, timestamp time.Time) string {
	return requestKey + "_" + strconv.FormatInt(timestamp.Unix(), 10)
}

//IncrementRequestKey Inscrements the amount of times a certain key was requested
func IncrementRequestKey(conn redis.Conn, requestKeyTimestamp string) {
	conn.Do("INCR", requestKeyTimestamp)
}

//GetTotalRequests Returns the total requests for the give period
func GetTotalRequests(
	conn redis.Conn,
	requestKeyWithoutTimestamp string,
	timestampStart int64,
	timestampEnd int64) (int, time.Duration) {

	if timestampEnd < timestampStart {
		log.Panicln("timestampEnd must be lower than timestampStart.")
	}

	total := 0
	seconds := 0
	for timestampEnd > timestampStart {
		requestKeyWithTimestamp := GetRequestKeyTimestamp(
			requestKeyWithoutTimestamp,
			time.Unix(timestampEnd, 0),
		)
		requestsNum, cmdErr := redis.Int(conn.Do("GET", requestKeyWithTimestamp))

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
	return total, time.Duration(seconds) * time.Second
}

//GetAvgRequests Returns avg requests for the given duration
func GetAvgRequests(totalRequests int, totalRequestsDuration time.Duration, desiredDuration time.Duration) int64 {
	if desiredDuration.Seconds() < time.Second.Seconds() {
		desiredDuration = time.Second
		log.Println("GetAvgRequests() received a desiredDuration too short. Forcing it to be 1 second.")
	}
	perSecondAvg := int64(totalRequests) / int64(totalRequestsDuration.Seconds())
	return perSecondAvg * int64(desiredDuration.Seconds())
}
