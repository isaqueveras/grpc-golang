package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	pb "github.com/isaqueveras/auth-microservice/proto"
	"google.golang.org/grpc"
)

const address = "localhost:50051"

func main() {
	router := gin.New()

	router.POST("/register", func(c *gin.Context) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}

		var data map[string]string
		if err := c.ShouldBindJSON(&data); err != nil {
			log.Printf("Could not register user: %v", err)
			return
		}

		client := pb.NewUserAuthClient(conn)
		res, err := client.RegisterUser(context.Background(), &pb.Register{
			Name:     data["name"],
			Email:    data["email"],
			Password: data["password"],
		})

		if err != nil {
			log.Printf("Could not register user: %v", err)
		}

		c.JSON(200, gin.H{"message": res.GetMessage()})
	})

	router.POST("/login", func(c *gin.Context) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}

		var data map[string]string
		if err := c.ShouldBindJSON(&data); err != nil {
			log.Printf("Could not register user: %v", err)
			return
		}

		client := pb.NewUserAuthClient(conn)
		res, err := client.LoginUser(context.Background(), &pb.Login{
			Email:    data["email"],
			Password: data["password"],
		})

		if err != nil {
			log.Printf("Could not register user: %v", err)
		}

		if res.Token != "" {
			c.SetCookie("token", res.Token, 3600, "/", "127.0.0.1", false, true)
		} else {
			log.Printf("Could not set cookie")
		}

		c.JSON(200, gin.H{
			"name":    res.GetName(),
			"email":   res.GetEmail(),
			"message": res.GetMessage(),
			"token":   res.GetToken(),
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
