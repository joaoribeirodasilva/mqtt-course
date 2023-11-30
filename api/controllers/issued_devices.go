package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IssuedDevice struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Type      string             `json:"type" bson:"type"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func IssuedDeviceList(c *gin.Context) {

	ptrs, err := mustGetAll(c)
	if err != nil {

		return
	}

	if !ptrs.User.Admin {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	collIssued := ptrs.Db.GetCollection("issueddevices")

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

	var cursor *mongo.Cursor
	cursor, err = collIssued.Find(context.TODO(), collIssued, &options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	issued := make([]IssuedDevice, 0)
	for cursor.Next(context.TODO()) {
		device := IssuedDevice{}
		if err := cursor.Decode(&device); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		issued = append(issued, device)
	}

	c.JSON(http.StatusOK, &issued)

}

func IssuedDeviceGet(c *gin.Context) {

	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	if !ptrs.User.Admin {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id := idQuery(c)
	if id == nil {
		return
	}

	issueddevice := &IssuedDevice{}
	collDevices := ptrs.Db.GetCollection("issueddevices")

	filter := bson.D{{Key: "_id", Value: id}}

	err = collDevices.FindOne(context.TODO(), filter).Decode(&issueddevice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &issueddevice)

}

func IssuedDeviceAdd(c *gin.Context) {

	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	if !ptrs.User.Admin {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	issuedDevice := &IssuedDevice{}

	if err := c.ShouldBind(issuedDevice); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	collDevices := ptrs.Db.GetCollection("issueddevices")

	issuedDevice.Type = strings.TrimSpace(issuedDevice.Type)
	if issuedDevice.Type == "" {

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid type"})
		return
	}

	issuedDevice.ID = primitive.NewObjectID()
	now := time.Now().UTC()
	issuedDevice.CreatedAt = now
	issuedDevice.UpdatedAt = now

	result, err := collDevices.InsertOne(context.TODO(), issuedDevice)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	id := result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, map[string]string{"id": id.Hex()})
}

func IssuedDeviceUpdate(c *gin.Context) {

	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	if !ptrs.User.Admin {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id := idQuery(c)
	if id == nil {
		return
	}

	issuedDevice := &IssuedDevice{}

	if err := c.ShouldBind(issuedDevice); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	issuedDevice.ID = *id

	dbIssueddevice := &IssuedDevice{}

	collDevices := ptrs.Db.GetCollection("issueddevices")

	filter := bson.D{
		{
			Key:   "_id",
			Value: issuedDevice.ID,
		},
	}

	err = collDevices.FindOne(context.TODO(), filter).Decode(dbIssueddevice)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	dbIssueddevice.Type = strings.TrimSpace(issuedDevice.Type)
	if dbIssueddevice.Type == "" {

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid type"})
		return
	}

	dbIssueddevice.UpdatedAt = time.Now().UTC()

	_, err = collDevices.UpdateOne(
		context.TODO(),
		bson.D{
			{
				Key:   "_id",
				Value: dbIssueddevice.ID,
			},
		},
		bson.D{
			{
				Key:   "$set",
				Value: dbIssueddevice,
			},
		},
	)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"_id": dbIssueddevice.ID.Hex()})
}

func IssuedDeviceDelete(c *gin.Context) {

	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	if !ptrs.User.Admin {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id := idQuery(c)
	if id == nil {
		return
	}

	collDevices := ptrs.Db.GetCollection("issueddevices")

	// filter by user id
	filter := bson.D{{Key: "_id", Value: id}}

	// deleted the user
	result, err := collDevices.DeleteOne(context.TODO(), filter)
	if err != nil {

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// checks if a document was deleted
	if result.DeletedCount == 0 {

		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}
