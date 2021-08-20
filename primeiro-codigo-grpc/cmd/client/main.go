package main

import (
	"context"
	"log"

	"github.com/isaqueveras/primeiro-codigo-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewSendMessageClient(conn)
	request := &pb.Request{
		Message: "Hello gRPC",
	}

	res, err := client.RequestMessage(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.GetStatus())
}
