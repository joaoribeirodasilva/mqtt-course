package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/mqtt-course/api/configuration"
	"github.com/joaoribeirodasilva/mqtt-course/api/database"
	"github.com/joaoribeirodasilva/mqtt-course/api/token"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListQuery struct {
	Page     int    `query:"p"`
	PageSize int    `query:"ps"`
	Sort     string `query:"s"`
	Dir      int    `query:"d"`
}

type Variables struct {
	Conf *configuration.Configuration
	Db   *database.Database
	User *token.User
}

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

func mustGetAll(c *gin.Context) (*Variables, error) {

	v := Variables{}
	ok := false

	co := c.MustGet("conf")
	v.Conf, ok = co.(*configuration.Configuration)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, errors.New("invalid configuration pointer")
	}

	d := c.MustGet("db")
	v.Db, ok = d.(*database.Database)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, errors.New("invalid database pointer")
	}

	v.User = nil
	a, exists := c.Get("auth")
	if exists {
		v.User = a.(*token.User)
	}

	return &v, nil

}

func listQuery(c *gin.Context) (*ListQuery, error) {

	query := &ListQuery{}

	if err := c.ShouldBind(query); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return nil, err
	}

	if query.PageSize == 0 {
		query.PageSize = defaultPageSize
	} else if query.PageSize > maxPageSize {
		query.PageSize = maxPageSize
	}

	query.Sort = strings.TrimSpace(query.Sort)

	if query.Dir > 0 {
		query.Dir = 1
	}

	return query, nil
}

func idQuery(c *gin.Context) *primitive.ObjectID {

	strId := c.Params.ByName("id")
	if strId == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return nil
	}
	id, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return nil
	}
	return &id
}
