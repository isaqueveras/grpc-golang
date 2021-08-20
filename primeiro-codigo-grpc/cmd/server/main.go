package main

import (
	"context"
	"log"
	"net"

	"github.com/isaqueveras/primeiro-codigo-grpc/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSendMessageServer
}

func (service *Server) RequestMessage(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Println("Mensagem recebida: ", req.GetMessage())
	response := &pb.Response{
		Status: 1,
	}

	return response, nil
}

func (service *Server) mustEmbedUnimplementedSendMessageServer() {}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterSendMessageServer(grpcServer, &Server{})

	port := ":9000"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
