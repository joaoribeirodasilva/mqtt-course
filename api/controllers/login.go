package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/mqtt-course/api/configuration"
	"github.com/joaoribeirodasilva/mqtt-course/api/database"
	"github.com/joaoribeirodasilva/mqtt-course/api/password"
	"github.com/joaoribeirodasilva/mqtt-course/api/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"id"`
	AccountID primitive.ObjectID `json:"accountId" bson:"accountId"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"_" bson:"password"`
	Name      string             `json:"name" bson:"name"`
	Surename  string             `json:"surename" bson:"surename"`
	Active    bool               `json:"active" bson:"active"`
}

func Login(c *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	tempConf := c.MustGet("conf")
	conf := tempConf.(*configuration.Configuration)

	d := c.MustGet("db")
	db := d.(*database.Database)

	collUsers := db.GetCollection("users")

	username, passwd, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	user := User{}

	err := collUsers.FindOne(context.TODO(), bson.M{"email": username}).Decode(&user)
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

	auth := token.New(conf)

	tokenUser := &token.User{
		ID:       user.ID.Hex(),
		Account:  user.AccountID.Hex(),
		Name:     user.Name,
		Surename: user.Surename,
	}

	if err := auth.Create(tokenUser); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"token": auth.TokenString})
}

func Signup(c *gin.Context) {

	c.Status(http.StatusOK)
}

func Logout(c *gin.Context) {

	// tempConf := c.MustGet("conf")
	// conf := tempConf.(*configuration.Configuration)

	d := c.MustGet("db")
	db := d.(*database.Database)
	db.GetCollection("metrics")

	c.Status(http.StatusOK)
}

func MetricsGet(c *gin.Context) {

	// tempConf := c.MustGet("conf")
	// conf := tempConf.(*configuration.Configuration)

	d := c.MustGet("db")
	db := d.(*database.Database)
	db.GetCollection("metrics")

	c.Status(http.StatusOK)
}
