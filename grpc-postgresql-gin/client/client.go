package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/isaqueveras/grpc-golang/grpc-postgresql-gin/client/user"
)

func main() {
	router := gin.New()

	v1 := router.Group("v1")
	user.Router(v1.Group("users"))

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
