package user

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/isaqueveras/grpc-golang/grpc-postgresql-gin/proto"
	"google.golang.org/grpc"
)

const address = "localhost:50051"

func cadastrarUsuario(c *gin.Context) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	connClient := pb.NewUserManagenentClient(conn)
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var new_users = make(map[string]int32)
	new_users["Jonathans"] = 25

	for name, age := range new_users {
		r, err := connClient.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}
		log.Printf(`%d, %s, tem %d anos`, r.GetId(), r.GetName(), r.GetAge())
	}
}

func listarUsuarios(c *gin.Context) {
	var r *pb.UserList

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	connClient := pb.NewUserManagenentClient(conn)
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	if r, err = connClient.GetUsers(ctx, &pb.GetUsersParams{}); err != nil {
		log.Fatalf("Could not retrieve users: %v", err)
	}

	c.JSON(200, r.GetUsers())
}
