package main

import (
	"fmt"
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

	fmt.Println("router variables")
	c.Set("db", r.db)
	c.Set("conf", r.conf)
}

func (r *Router) IsLogged(c *gin.Context) {

	fmt.Println("router is logged")
	auth := token.New(r.conf)
	if !auth.IsValid(c.GetHeader("Authorization")) {
		c.AbortWithStatus(http.StatusForbidden)
		c.Abort()
		return
	}

	c.Set("auth", auth.User)

	c.Next()
}

func (r *Router) IsAdmin(c *gin.Context) {

	fmt.Println("router is admin")
	a := c.MustGet("auth")
	user, ok := a.(*token.User)
	if !ok || user.Admin {
		c.AbortWithStatus(http.StatusUnauthorized)
		c.Abort()
		return
	}

	c.Next()
}

func (r *Router) SetRoutes() {

	// Login related
	r.gin.POST("/login", r.Variables, controllers.Login)
	r.gin.POST("/signup", r.Variables, controllers.UserAdd)
	r.gin.DELETE("/logout", r.Variables, r.IsLogged, controllers.Logout)

	r.gin.GET("/devices", r.Variables, r.IsLogged, controllers.DeviceList)
	r.gin.GET("/device/:id", r.Variables, r.IsLogged, controllers.DeviceGet)
	r.gin.POST("/device", controllers.DeviceAdd)
	r.gin.PUT("/device/:id", r.Variables, r.IsLogged, controllers.DeviceUpdate)
	r.gin.PATCH("/device/:id", r.Variables, r.IsLogged, controllers.DeviceUpdate)
	r.gin.DELETE("/device/:id", r.Variables, r.IsLogged, controllers.DeviceDelete)

	r.gin.GET("/users", r.Variables, r.IsLogged, r.IsAdmin, controllers.UserList)
	r.gin.GET("/user/:id", r.Variables, r.IsLogged, controllers.UserGet)
	r.gin.POST("/user", r.Variables, controllers.UserAdd)
	r.gin.PUT("/user/:id", r.Variables, r.IsLogged, controllers.UserUpdate)
	r.gin.PATCH("/user/:id", r.Variables, r.IsLogged, controllers.UserUpdate)
	r.gin.DELETE("/user/:id", r.Variables, r.IsLogged, controllers.UserDelete)

	r.gin.GET("/issueddevices/", r.Variables, r.IsLogged, r.IsAdmin, controllers.IssuedDeviceList)
	r.gin.GET("/issueddevice/:id", r.Variables, r.IsLogged, r.IsAdmin, controllers.IssuedDeviceList)
	r.gin.POST("/issueddevice", r.Variables, r.IsLogged, r.IsAdmin, controllers.IssuedDeviceList)
	r.gin.PUT("/issueddevice/:id", r.Variables, r.IsLogged, r.IsAdmin, controllers.IssuedDeviceList)
	r.gin.PATCH("/issueddevice/:id", r.Variables, r.IsLogged, r.IsAdmin, controllers.IssuedDeviceList)
	r.gin.DELETE("/issueddevice/:id", r.Variables, r.IsLogged, r.IsAdmin, controllers.IssuedDeviceList)

	// Metrics related
	r.gin.GET("/metrics/:device/:start/:end", r.Variables, r.IsLogged, r.IsLogged, controllers.Login)

}
