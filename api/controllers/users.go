package controllers

import (
	"context"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/mqtt-course/api/password"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Email     string             `json:"email,omitempty" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Name      string             `json:"name" bson:"name"`
	Surename  string             `json:"surename" bson:"surename"`
	Admin     bool               `json:"admin" bson:"admin"`
	Active    bool               `json:"active" bson:"active"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func UserList(c *gin.Context) {

	// gets all service pointers from middleware
	ptrs, err := mustGetAll(c)
	if err != nil {

		return
	}

	// failsafe. only Admin users are authorized to gat
	// a users list. although this rule is also set in
	// the router middleware
	if !ptrs.User.Admin {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// gets the users collection
	collUsers := ptrs.Db.GetCollection("users")

	// gets the list query parameters
	query, err := listQuery(c)
	if err != nil {

		return
	}

	// by defaul surts the records bu the first name ASC and surename ASC
	sort := bson.D{
		{
			Key:   "name",
			Value: 0,
		},
		{
			Key:   "surename",
			Value: 0,
		},
	}

	// if sort data is present on the query string
	// sets it
	if query.Sort != "" {
		sort = bson.D{
			{
				Key:   query.Sort,
				Value: query.Dir,
			},
		}
	}

	// page to be fetched
	limit := int64(query.PageSize)

	// number of documents to be skipped
	skip := int64(query.Page*query.PageSize - query.PageSize)

	// set the find options
	options := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  sort,
	}

	// no filter (gets all users)
	filter := bson.D{}

	// alocates the result cursor
	var cursor *mongo.Cursor
	cursor, err = collUsers.Find(context.TODO(), filter, &options)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// alocates the users list
	users := make([]User, 0)

	// for each document
	for cursor.Next(context.TODO()) {

		user := User{}
		if err := cursor.Decode(&user); err != nil {

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		users = append(users, user)
	}

	// returns the users list
	c.JSON(http.StatusOK, &users)

}

func UserGet(c *gin.Context) {

	// gets all service pointers from middleware
	ptrs, err := mustGetAll(c)
	if err != nil {

		return
	}

	// gets the user id from query string
	id := idQuery(c)
	if id == nil {

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check if the logged user is an Admin user
	// or if the logged user is the account owner
	if !ptrs.User.Admin && id.Hex() != ptrs.User.ID.Hex() {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// gets the users collection
	collUsers := ptrs.Db.GetCollection("users")

	// filter by user id
	filter := bson.D{{Key: "_id", Value: id}}

	user := User{}

	// tries to find the user
	if err := collUsers.FindOne(context.TODO(), filter).Decode(&user); err != nil {

		if err == mongo.ErrNoDocuments {

			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// don't send the password information
	user.Password = ""

	// returns the user
	c.JSON(http.StatusOK, &user)
}

func UserAdd(c *gin.Context) {

	// gets all service pointers from middleware
	ptrs, err := mustGetAll(c)
	if err != nil {

		return
	}

	user := User{}

	// gets the users collection
	collUsers := ptrs.Db.GetCollection("users")

	// gets the new user data from the payload
	if err := c.BindJSON(&user); err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	// the user being created is an Admin user and
	// only loged admin users may set this
	if !ptrs.User.Admin && user.Admin {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// trims any spaces
	user.ID = primitive.NewObjectID()
	user.Email = strings.TrimSpace(user.Email)
	user.Name = strings.TrimSpace(user.Name)
	user.Surename = strings.TrimSpace(user.Surename)
	user.Password = strings.TrimSpace(user.Password)

	// parses the email address
	_, err = mail.ParseAddress(user.Email)
	if err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid email"})
		return
	}

	// check the user first name
	if user.Name == "" {

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid name"})
		return
	}

	// check the user surename
	if user.Surename == "" {

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid surename"})
		return
	}

	// check the user password
	if len(user.Password) < 6 {

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid password"})
		return
	}

	// hashes the user pasword
	user.Password, err = password.Hash(user.Password)
	if err != nil {

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// updates times metadata
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt

	// check if the email is already registred
	err = collUsers.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&user)
	if err != nil {

		if err != mongo.ErrNoDocuments {

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {

		c.AbortWithStatusJSON(http.StatusConflict, map[string]string{"error": "user already exists"})
		return
	}

	// inserts the new user
	_, err = collUsers.InsertOne(context.TODO(), &user)
	if err != nil {

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// returns the new user id
	c.JSON(http.StatusCreated, map[string]string{"id": user.ID.Hex()})
}

func UserUpdate(c *gin.Context) {

	// gets all service pointers from middleware
	ptrs, err := mustGetAll(c)
	if err != nil {

		return
	}

	// gets the user id from query string
	id := idQuery(c)
	if id == nil {

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := User{}

	// gets the users collection
	collUsers := ptrs.Db.GetCollection("users")

	// gets the new user data from the payload
	if err := c.BindJSON(&user); err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	// trims any spaces
	user.Email = strings.TrimSpace(user.Email)
	user.Name = strings.TrimSpace(user.Name)
	user.Surename = strings.TrimSpace(user.Surename)
	user.Password = strings.TrimSpace(user.Password)

	// check if the user making the change is an Admin user
	// or if the user making the change is the account owner
	if !ptrs.User.Admin && id.Hex() != ptrs.User.ID.Hex() {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// the user being changed is an Admin user and
	// only loged admin user may change this status
	if !ptrs.User.Admin && user.Admin {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// filters by user id
	filter := bson.D{{Key: "_id", Value: id}}

	dbUser := &User{}

	// tries to find the user
	err = collUsers.FindOne(context.TODO(), filter).Decode(&dbUser)
	if err != nil {

		// if the user is no found
		if err == mongo.ErrNoDocuments {

			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// any other errors
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// parses the email address
	_, err = mail.ParseAddress(user.Email)
	if err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid email"})
		return
	}

	// checks the user first name
	if user.Name == "" {

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid name"})
		return
	}

	// checks the user first surename
	if user.Surename == "" {

		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid surename"})
		return
	}

	// check if password was changed
	if user.Password != "" {

		// validates the password
		if len(user.Password) < 6 {

			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "invalid password"})
			return
		}
		// hashes the new password
		dbUser.Password, err = password.Hash(user.Password)
		if err != nil {

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	}

	// Update changes email
	if dbUser.Email != user.Email {

		emailUser := &User{}

		// Tries to find the same email in another
		// user id.
		err := collUsers.FindOne(
			context.TODO(),
			bson.D{
				{
					Key:   "email",
					Value: user.Email,
				},
				{
					Key: "_id",
					Value: bson.D{
						{
							Key:   "$ne",
							Value: user.ID,
						},
					},
				},
			},
		).Decode(&emailUser)

		if err != nil && err != mongo.ErrNoDocuments {

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		} else if err == nil {

			// query returned a user, so the email
			// address already belongs to another
			// user account
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		dbUser.Email = user.Email
	}

	dbUser.Name = user.Name
	dbUser.Surename = user.Surename

	// sets the new update metadata
	dbUser.UpdatedAt = time.Now().UTC()

	_, err = collUsers.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: dbUser.ID}}, bson.D{{Key: "$set", Value: dbUser}})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// returns the updated user id
	c.JSON(http.StatusOK, map[string]string{"id": id.Hex()})
}

func UserDelete(c *gin.Context) {

	// get all service pointers from middleware
	ptrs, err := mustGetAll(c)
	if err != nil {

		return
	}

	// gets the user id from query string
	id := idQuery(c)
	if id == nil {

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// gets the users collection
	collUsers := ptrs.Db.GetCollection("users")

	// check if the user making the change is an Admin user
	// or if the user making the change is the account owner
	if !ptrs.User.Admin && id.Hex() != ptrs.User.ID.Hex() {

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// filter by user id
	filter := bson.D{{Key: "_id", Value: id}}

	// deleted the user
	result, err := collUsers.DeleteOne(context.TODO(), filter)
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
