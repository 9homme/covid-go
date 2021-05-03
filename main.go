package main

import (
	"example.com/covid-go/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	repository.InitMongoDb()
	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(Logger())
	r.GET("/ping", PingHandler)
	authorized := r.Group("/admin", BasicAuth())
	authorized.GET("/covid19", Covid19Handler)
	authorized.GET("/users", GetUsersHandler)
	return r
}
