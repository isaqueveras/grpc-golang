package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/isaqueveras/grpc-golang/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pb.NewAddServiceClient(conn)
	g := gin.Default()
	g.GET("/add/:a/:b", func(c *gin.Context) {
		a, _ := strconv.ParseUint(c.Param("a"), 10, 64)
		b, _ := strconv.ParseUint(c.Param("b"), 10, 64)

		response, err := client.Add(c, &pb.Request{A: int64(a), B: int64(b)})
		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{"result": response.Result})
	})

	g.GET("/mult/:a/:b", func(c *gin.Context) {
		a, _ := strconv.ParseUint(c.Param("a"), 10, 64)
		b, _ := strconv.ParseUint(c.Param("b"), 10, 64)

		response, err := client.Multiply(c, &pb.Request{A: int64(a), B: int64(b)})
		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{"result": response.Result})
	})

	if err = g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run servert: %v", err)
	}
}
