package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	pb "github.com/isaqueveras/auth-microservice/proto"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

// port of server
const (
	port      = ":50051"
	SecretKey = "isaque"
)

type Server struct {
	conn *pgx.Conn
	pb.UnimplementedUserAuthServer
}

func NewAuthServer() *Server {
	return &Server{}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	s.UnimplementedUserAuthServer = pb.UnimplementedUserAuthServer{}
	pb.RegisterUserAuthServer(server, s)
	log.Printf("Server listening at %v", lis.Addr())

	return server.Serve(lis)
}

// RegisterUser register user on database
func (s *Server) RegisterUser(ctx context.Context, in *pb.Register) (*pb.Message, error) {
	// Generate password
	password, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), 14)
	if err != nil {
		log.Fatalf("Erro on generate password: %v", err)
	}

	tx, _ := s.conn.Begin(ctx)
	_, err = tx.Exec(ctx, "INSERT INTO users_test_2(name, email, passw) VALUES ($1, $2, $3)", in.GetName(), in.GetEmail(), password)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Panic(err)
	}

	return &pb.Message{Message: "The user has been registered in the database"}, nil
}

func (s *Server) LoginUser(ctx context.Context, in *pb.Login) (*pb.LoginRes, error) {
	tx, _ := s.conn.Begin(ctx)
	res, err := tx.Query(ctx, "SELECT id, name, email, passw FROM users_test_2 WHERE email = $1", in.GetEmail())
	if err != nil {
		log.Fatalf("Unable to fetch the user from the database: %v", err.Error())
	}

	defer res.Close()

	user := &pb.Register{}
	for res.Next() {
		if err = res.Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
			log.Fatalf("Erro in scan: %v", err.Error())
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.GetPassword())); err != nil {
		return &pb.LoginRes{Message: "Incorrect password"}, nil
	}

	clams := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
		Issuer:    strconv.Itoa(int(user.Id)),
	})

	token, err := clams.SignedString([]byte(SecretKey))
	if err != nil {
		return &pb.LoginRes{Message: "Could not login"}, nil
	}

	return &pb.LoginRes{
		Token:   token,
		Name:    user.Name,
		Email:   user.Email,
		Message: "Seja bem vindo(a), " + user.Name,
	}, nil
}

func main() {
	// Initial context
	ctx := context.Background()

	// Connection into database
	conn, err := pgx.Connect(ctx, "postgres://postgres:123456@localhost:5432/postgres")
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}

	defer conn.Close(ctx)

	var userAuth *Server = NewAuthServer()
	userAuth.conn = conn

	if err = userAuth.Run(); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
