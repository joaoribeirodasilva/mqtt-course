package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type MetricsQuery struct {
	DeviceId string    `query:"deviceId"`
	Start    time.Time `query:"start" time_format:"2006-01-02T15:04:05Z" time_utc:"0"`
	End      time.Time `query:"start" time_format:"2006-01-02T15:04:05Z" time_utc:"0"`
}

const (
	defaultDataFormat = "2006-01-02T15:04:05.000Z"
)

// MetricsGet
func MetricsGet(c *gin.Context) {

	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	// gets the device id from query string
	id := idQuery(c)
	if id == nil {
		return
	}

	strStart := c.Params.ByName("start")
	if strStart == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "start date/time parameter is required"})
		return
	}

	strEnd := c.Params.ByName("end")
	if strEnd == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "end date/time parameter is required"})
		return
	}

	start, err := time.Parse(defaultDataFormat, strStart)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "start date/time parameter is invalid"})
		return
	}

	end, err := time.Parse(defaultDataFormat, strEnd)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "end date/time parameter is invalid"})
		return
	}

	filter := bson.D{{Key: "metrics.deviceId", Value: id}}

	if end.Sub(start) < 0 {
		tempStart := start
		start = end
		end = tempStart
	}

	if !ptrs.User.Admin {
		filter = append(filter, bson.E{Key: "userId", Value: ptrs.User.ID})
	}

	filter = append(filter, bson.E{Key: "metric.collectedAt", Value: bson.E{Key: "$gte", Value: start}})
	filter = append(filter, bson.E{Key: "metric.collectedAt", Value: bson.E{Key: "$lte", Value: end}})

	c.Status(http.StatusOK)
}
