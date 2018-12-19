package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/static"
	"log"
	"github.com/sjeanpierre/rs_input_tracker_go/app"
)


func main() {
	log.Print("-Startup-")
	//InputFunctions()
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./public/rs_audit/", true)))
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	})
	api := router.Group("/api")
	app.RegisterRoutes(api)
	router.Run(":9080")
}

