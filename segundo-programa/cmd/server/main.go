package main

import (
	"context"
	"log"
	"net"

	"github.com/isaqueveras/segundo-programa/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSendMessageServer
}

func (service *Server) RequestMessage(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Println("Usuario: ", req.GetNameUser())
	log.Println("Mensagem: ", req.GetMessage())

	response := &pb.Response{
		Message: "Mensagem enviado com sucesso",
	}

	return response, nil
}

func main() {
	grpcSever := grpc.NewServer()
	pb.RegisterSendMessageServer(grpcSever, &Server{})

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal(err)
	}

	err = grpcSever.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
