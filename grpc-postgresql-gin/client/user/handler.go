package user

import (
	"context"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/isaqueveras/grpc-golang/grpc-postgresql-gin/client/config"
	pb "github.com/isaqueveras/grpc-golang/grpc-postgresql-gin/proto"
)

func cadastrarUsuario(c *gin.Context) {
	client := config.Conn()

	var (
		user UserReq
		r    *pb.User
		err  error
	)

	if err = c.ShouldBindJSON(&user); err != nil {
		if err = errors.New("Informe o nome ou a idade"); err != nil {
			return
		}
		return
	}

	// Passando os dados para cadastrar o usuario no sistema do gRPC
	if r, err = client.CreateNewUser(context.Background(), &pb.NewUser{Name: user.Name, Age: user.Age}); err != nil {
		log.Fatalf("Could not create user: %v", err)
	}

	c.JSON(200, r)
}

func listarUsuarios(c *gin.Context) {
	client := config.Conn()

	r, err := client.GetUsers(context.Background(), &pb.GetUsersParams{})
	if err != nil {
		log.Fatalf("Could not retrieve users: %v", err)
	}

	c.JSON(200, r.GetUsers())
}
