package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	escutando, erro := net.Listen("tcp", ":9000")
	if erro != nil {
		log.Fatalf("Falha ao escutar na porta 9000: %v", erro)
	}

	server := grpc.NewServer()
	if erro = server.Serve(escutando); erro != nil {
		log.Fatalf("Falha no servidor gRPC usando a porta 9000: %v", erro)
	}
}
