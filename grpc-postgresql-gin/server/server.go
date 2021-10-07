package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/isaqueveras/grpc-golang/grpc-postgresql-gin/proto"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

const port = ":50051"

func NewUserManagenentServer() *Server {
	return &Server{
		userList: &pb.UserList{},
	}
}

type Server struct {
	conn     *pgx.Conn
	userList *pb.UserList
	pb.UnimplementedUserManagenentServer
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterUserManagenentServer(server, s.UnimplementedUserManagenentServer)
	log.Printf("Server listening at %v", lis.Addr())

	return server.Serve(lis)
}

func (s *Server) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())

	createSql := `
		CREATE TABLE IF NOT EXISTS users_test(
			id SERIAL PRIMARY KEY,
			name text,
			age int
		);`

	_, err := s.conn.Exec(ctx, createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
		os.Exit(1)
	}

	user := &pb.User{Name: in.GetName(), Age: in.GetAge()}

	tx, _ := s.conn.Begin(ctx)
	_, err = tx.Exec(ctx, "INSERT INTO users_test(name, age) VALUES ($1, $2)", user.Name, user.Age)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Panic(err)
	}

	return user, nil
}

func (s *Server) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var (
		userList *pb.UserList = &pb.UserList{}
		err      error
		rows     pgx.Rows
	)

	if rows, err = s.conn.Query(ctx, "SELECT id, name, age FROM users_test ORDER BY id DESC"); err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		user := &pb.User{}
		if err = rows.Scan(&user.Id, &user.Name, &user.Age); err != nil {
			log.Print(err)
		}
		userList.Users = append(userList.Users, user)
	}

	return userList, nil
}

func (s *Server) DeleteUser(ctx context.Context, in *pb.DeleteUserReq) (*pb.DeleteUserRes, error) {
	var user *pb.DeleteUserRes = &pb.DeleteUserRes{}

	log.Print("id of user: ", in.GetId())

	tx, _ := s.conn.Begin(ctx)
	res, err := tx.Exec(ctx, "DELETE FROM users_test WHERE id = $1", in.GetId())
	if err != nil {
		log.Fatalf("Could not delete user with id: %v", in.GetId())
	}

	if res.RowsAffected() != 1 {
		return &pb.DeleteUserRes{Message: "User not found to delete"}, nil
	}

	user.Message = "User delete with success!"
	return user, nil
}

func main() {
	var (
		userMgmtServer *Server = NewUserManagenentServer()
		ctx                    = context.Background()
		conn           *pgx.Conn
		err            error
	)

	if conn, err = pgx.Connect(ctx, "postgres://postgres:123456@localhost:5432/postgres"); err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(ctx)

	userMgmtServer.conn = conn

	if err := userMgmtServer.Run(); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
