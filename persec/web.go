package persec

import (
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
)

//AddRequestResponse Add response struct
type AddRequestResponse struct {
	total int
	key   string
}

//GetRequestResponse Get response struct
type GetRequestResponse struct {
	total         int
	key           string
	optimisticAvg int64
	realAvg       int64
}

func handleAddRequest(context echo.Context) error {

	redisConn := context.Get("redisConn").(redis.Conn)
	key := context.Param("key")
	requestKeyTimestamp := GetRequestKeyTimestamp(
		key,
		time.Now())

	incr, err := IncrementRequestKey(redisConn, requestKeyTimestamp)

	if err != nil {
		context.Logger().Errorf("Redis error: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return context.JSON(http.StatusOK, AddRequestResponse{
		total: incr,
		key:   key})
}

func handleGetRequest(context echo.Context) error {

	redisConn := context.Get("redisConn").(redis.Conn)
	startTimestamp := context.Get("startTimestamp").(int64)
	endTimestamp := context.Get("endTimestamp").(int64)
	key := context.Param("key")
	desiredDuration := time.Unix(startTimestamp, 0).Sub(time.Unix(endTimestamp, 0))

	total, duration, err := GetTotalRequests(
		redisConn,
		key,
		startTimestamp,
		endTimestamp)

	if err != nil {
		context.Logger().Errorf("Redis error: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return context.JSON(http.StatusOK, GetRequestResponse{
		total: total,
		key:   key,
		optimisticAvg: GetOptimisticAvgRequests(
			total,
			duration,
			desiredDuration,
		),
		realAvg: GetRealAvgRequests(total, desiredDuration),
	})
}
