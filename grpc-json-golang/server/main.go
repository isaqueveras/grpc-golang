package main

import (
	"context"
	"log"
	"net"

	pb "github.com/isaqueveras/grpc-golang/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAddServiceServer
}

const port = ":4040"

func (s *server) Add(_ context.Context, in *pb.Request) (*pb.Response, error) {
	a, b := in.GetA(), in.GetB()
	return &pb.Response{Result: a + b}, nil
}

func (s *server) Multiply(_ context.Context, in *pb.Request) (*pb.Response, error) {
	a, b := in.GetA(), in.GetB()
	return &pb.Response{Result: a * b}, nil
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	srv := grpc.NewServer()
	pb.RegisterAddServiceServer(srv, &server{})
	log.Printf("Server listening at %v", listener.Addr())

	if err = srv.Serve(listener); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
