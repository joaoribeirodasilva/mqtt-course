package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Device struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	Name           string             `json:"name" bson:"name"`
	LastMetricTime *time.Time         `json:"lastMetricTime" bson:"lastMetricTime"`
	Active         bool               `json:"active" bson:"active"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func DeviceList(c *gin.Context) {

	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	collDevices := ptrs.Db.GetCollection("devices")

	query, err := listQuery(c)
	if err != nil {
		return
	}

	sort := bson.D{{Key: "updatedAt", Value: 1}}
	if query.Sort != "" {
		sort = bson.D{{Key: query.Sort, Value: query.Dir}}
	}

	limit := int64(query.PageSize)
	skip := int64(query.Page*query.PageSize - query.PageSize)
	options := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  sort,
	}

	filter := bson.D{{Key: "userId", Value: ptrs.User.ID}}
	if ptrs.User.Admin {
		filter = bson.D{}
	}

	var cursor *mongo.Cursor
	cursor, err = collDevices.Find(context.TODO(), filter, &options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	devices := make([]Device, 0)
	for cursor.Next(context.TODO()) {
		device := Device{}
		if err := cursor.Decode(&device); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		devices = append(devices, device)
	}

	c.JSON(http.StatusOK, &devices)
}

func DeviceGet(c *gin.Context) {

	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	id := idQuery(c)
	if id == nil {
		return
	}

	device := &Device{}
	collDevices := ptrs.Db.GetCollection("devices")

	filter := bson.D{{Key: "_id", Value: id}, {Key: "userId", Value: ptrs.User.ID}}
	if ptrs.User.Admin {
		filter = bson.D{{Key: "_id", Value: id}}
	}

	err = collDevices.FindOne(context.TODO(), filter).Decode(&device)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &device)
}

func DeviceAdd(c *gin.Context) {

	// get all service pointers from middleware
	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	// alocates the device struct
	device := &Device{}

	if err := c.ShouldBind(device); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// gets the device new data from the payload
	collDevices := ptrs.Db.GetCollection("devices")

	// if the logged user is not an Admin
	// assigns the logged user id to the
	// device's user id
	if !ptrs.User.Admin {
		device.UserID = ptrs.User.ID
	}

	now := time.Now().UTC()
	device.LastMetricTime = nil
	device.CreatedAt = now
	device.UpdatedAt = now

	// check if the device id exists in the
	// issueddevices collection
	collIssuedDevices := ptrs.Db.GetCollection("issueddevices")

	filter := bson.D{{Key: "_id", Value: device.ID}}

	issuedDevice := &IssuedDevice{}

	err = collIssuedDevices.FindOne(context.TODO(), filter).Decode(issuedDevice)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// check if the device already belongs and
	// it's active on another account
	dbDevice := &Device{}

	filter = bson.D{{Key: "_id", Value: device.ID}, {Key: "active", Value: true}}

	err = collDevices.FindOne(context.TODO(), filter).Decode(dbDevice)
	if err != nil && err != mongo.ErrNoDocuments {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	} else if err == nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	// inserts the device
	_, err = collDevices.InsertOne(context.TODO(), device)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// return the updated device id
	c.JSON(http.StatusCreated, map[string]string{"id": device.ID.Hex()})
}

func DeviceUpdate(c *gin.Context) {

	// get all service pointers from middleware
	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	// gets the device id from query string
	id := idQuery(c)
	if id == nil {
		return
	}

	// alocates the device struct
	device := &Device{}

	// gets the device new data from the payload
	if err := c.ShouldBind(device); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// sets the device id to the id in
	// the query string
	device.ID = *id

	// gets the devices collection
	collDevices := ptrs.Db.GetCollection("devices")

	// sets the device update time to now
	device.UpdatedAt = time.Now().UTC()
	if !ptrs.User.Admin {
		device.UserID = ptrs.User.ID
	}

	filter := bson.D{{Key: "_id", Value: id}, {Key: "userId", Value: ptrs.User.ID}}
	if ptrs.User.Admin {
		filter = bson.D{{Key: "_id", Value: id}}
	}

	// declares the result pointer
	var result *mongo.UpdateResult

	// updated the device according to the
	// filter
	result, err = collDevices.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: device}})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// if the device wasn't found returns
	// not found
	if result.MatchedCount == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// return the updated device id
	c.JSON(http.StatusOK, map[string]string{"id": id.Hex()})
}

func DeviceDelete(c *gin.Context) {

	// get all service pointers from middleware
	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	// gets the device id from query string
	id := idQuery(c)
	if id == nil {
		return
	}

	// gets the devices collection
	collDevices := ptrs.Db.GetCollection("devices")

	// filter by user id and device admin
	// if the logged user in not Admin
	filter := bson.D{{Key: "_id", Value: id}, {Key: "userId", Value: ptrs.User.ID}}
	if ptrs.User.Admin {
		filter = bson.D{{Key: "_id", Value: id}}
	}

	// deletes the device according to the
	// filter
	result, err := collDevices.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// if no document was found in the database
	// return not found
	if result.DeletedCount == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// returns ok
	c.Status(http.StatusOK)
}
