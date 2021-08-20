package main

import (
	"context"
	"log"

	"github.com/isaqueveras/segundo-programa/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewSendMessageClient(conn)
	request := &pb.Request{
		NameUser: "Isaque Véras",
		Message:  "Olá mundo!",
	}

	res, err := client.RequestMessage(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.GetMessage())
}
