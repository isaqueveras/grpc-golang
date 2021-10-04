package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/isaqueveras/grpc-golang/go-grpc-postgresql/proto-user"
	"google.golang.org/grpc"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewUserManagenentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// var new_users = make(map[string]int32)
	// new_users["Rebeca Veras"] = 13

	// for name, age := range new_users {
	// 	r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
	// 	if err != nil {
	// 		log.Fatalf("could not create user: %v", err)
	// 	}
	// 	log.Printf(`%d, %s, tem %d anos`, r.GetId(), r.GetName(), r.GetAge())
	// }

	var r *pb.UserList
	if r, err = c.GetUsers(ctx, &pb.GetUsersParams{}); err != nil {
		log.Fatalf("Could not retrieve users: %v", err)
	}

	fmt.Printf("%v\n", r.GetUsers())
}
