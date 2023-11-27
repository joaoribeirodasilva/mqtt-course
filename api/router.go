package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/mqtt-course/api/configuration"
	"github.com/joaoribeirodasilva/mqtt-course/api/controllers"
	"github.com/joaoribeirodasilva/mqtt-course/api/database"
	"github.com/joaoribeirodasilva/mqtt-course/api/token"
)

type Router struct {
	conf *configuration.Configuration
	gin  *gin.Engine
	db   *database.Database
}

func NewRouter(gin *gin.Engine, conf *configuration.Configuration, db *database.Database) *Router {

	r := &Router{}

	r.conf = conf
	r.gin = gin
	r.db = db

	return r
}

func (r *Router) Variables(c *gin.Context) {

	c.Set("db", r.db)
	c.Set("conf", r.conf)
}

func (r *Router) IsLogged(c *gin.Context) {

	auth := token.New(r.conf)
	if !auth.IsValid(c.GetHeader("Authorization")) {
		c.AbortWithStatus(http.StatusForbidden)
		c.Abort()
		return
	}

	c.Set("auth", auth.User)

	c.Next()
}

func (r *Router) SetRoutes() {

	r.gin.POST("/login", controllers.Login)
	r.gin.POST("/signup", controllers.Signup)

	r.gin.Use(r.Variables).DELETE("/logout", r.IsLogged, controllers.Login)

	r.gin.Use(r.Variables, r.IsLogged).GET("/metrics/:clientId/:start/:end", r.IsLogged, controllers.Login)

}
