package main

import (
	"log"
	"net"

	"github.com/isaqueveras/grpc-golang/proto/gen"
	"google.golang.org/grpc"
)

type Server struct {
	gen.UnimplementedUserServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	gen.RegisterUserServiceServer(grpcServer, &Server{})

	log.Println("Escutando na porta :9000")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Falta ao iniciar o servidor: %s", err)
	}
}
