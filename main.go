package main

import "github.com/gin-gonic/gin"

func HttpRouter() *gin.Engine {
	router := gin.Default()
	{
		run := router.Group("/run")
		run.POST("/", RunAgent)
	}
	return router
}

func RunAgent(c *gin.Context) {}

func main() {
	router := HttpRouter()
	router.Run()
}
