package config

import (
	"log"

	pb "github.com/isaqueveras/grpc-golang/grpc-postgresql-gin/proto"
	"google.golang.org/grpc"
)

const address = "localhost:50051"

func Conn() pb.UserManagenentClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	connClient := pb.NewUserManagenentClient(conn)
	return connClient
}
