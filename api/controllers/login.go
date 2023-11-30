package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/mqtt-course/api/password"
	"github.com/joaoribeirodasilva/mqtt-course/api/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Session struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	StartTime time.Time          `json:"startTime" bson:"startTime"`
	EndTime   *time.Time         `json:"endTime" bson:"endTime"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// Login
func Login(c *gin.Context) {

	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	collUsers := ptrs.Db.GetCollection("users")

	username, passwd, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	user := User{}

	err = collUsers.FindOne(context.TODO(), bson.M{"email": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !user.Active {
		c.AbortWithStatus(http.StatusLocked)
		return
	}

	if !password.Check(passwd, user.Password) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	auth := token.New(ptrs.Conf)

	tokenUser := &token.User{
		ID:       user.ID,
		Name:     user.Name,
		Surename: user.Surename,
	}

	if err := auth.Create(tokenUser); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	collSessions := ptrs.Db.GetCollection("sessions")

	now := time.Now().UTC()

	session := Session{
		ID:        primitive.NewObjectID(),
		UserID:    user.ID,
		StartTime: now,
		EndTime:   nil,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = collSessions.InsertOne(context.TODO(), &session)
	if err != nil {
		// fmt.Printf("ERROR: %+v\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"token": auth.TokenString})
}

// Logout
func Logout(c *gin.Context) {

	ptrs, err := mustGetAll(c)
	if err != nil {
		return
	}

	collSessions := ptrs.Db.GetCollection("sessions")

	session := Session{}

	err = collSessions.FindOne(context.TODO(), bson.D{{Key: "userId", Value: ptrs.User.ID}, {Key: "endTime", Value: nil}}).Decode(&session)
	if err != nil && err != mongo.ErrNoDocuments {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	} else if err == mongo.ErrNoDocuments {
		c.Status(http.StatusOK)
		return
	}

	filter := bson.D{{Key: "_id", Value: session.ID}}

	now := time.Now().UTC()
	session.EndTime = &now
	session.UpdatedAt = now

	_, err = collSessions.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: session}})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
