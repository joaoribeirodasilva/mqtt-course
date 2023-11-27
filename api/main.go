package main

import (
	"log"

	"github.com/joaoribeirodasilva/mqtt-course/api/configuration"
	"github.com/joaoribeirodasilva/mqtt-course/api/database"
)

func main() {

	conf := configuration.NewConfiguration()
	if err := conf.Read(); err != nil {
		panic(err)
	}

	db := database.NewDatabase(conf)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}

	http := NewServer(conf)
	router := NewRouter(http.Router, conf, db)

	router.SetRoutes()

	if err := http.Listen(); err != nil {
		panic(err)
	}

	db.Disconnect()

}
