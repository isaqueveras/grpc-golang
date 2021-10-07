package user

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isaqueveras/grpc-golang/grpc-postgresql-gin/client/config"
	pb "github.com/isaqueveras/grpc-golang/grpc-postgresql-gin/proto"
)

// cadastrarUsuario insert user in database
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
		log.Printf("Could not create user: %v", err)
	}

	c.JSON(200, r)
}

// listarUsuarios list all users in database
func listarUsuarios(c *gin.Context) {
	client := config.Conn()

	r, err := client.GetUsers(c, &pb.GetUsersParams{})
	if err != nil {
		log.Printf("Could not retrieve users: %v", err)
	}

	c.JSON(200, r.GetUsers())
}

// deleteUser delete a user in database
func deleteUser(c *gin.Context) {
	idUser, err := strconv.ParseInt(c.Param("id_user"), 10, 64)
	if err != nil {
		log.Printf("Could not get id user")
		c.JSON(500, gin.H{"message": "Erro! inform the identifier of user"})
		return
	}

	client := config.Conn()
	response, err := client.DeleteUser(c, &pb.DeleteUserReq{Id: int32(idUser)})
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not delete user"})
		return
	}

	c.JSON(200, gin.H{
		"message": response.GetMessage(),
	})
}
