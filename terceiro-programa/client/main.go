package main

import (
	"context"
	"fmt"
	"log"

	user "github.com/isaqueveras/grpc-golang/proto/gen"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	usr := user.NewUserServiceClient(conn)

	fmt.Print("Digite o username: ")
	var inputUserName string
	fmt.Scanln(&inputUserName)

	response, err := usr.GetUser(context.Background(), &user.UserRequest{
		Username: inputUserName,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println(response)
}
