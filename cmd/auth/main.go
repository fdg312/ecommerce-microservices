package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fdg312/ecommerce-microservices/internal/grpc"
	"github.com/fdg312/ecommerce-microservices/internal/repository"
	"github.com/fdg312/ecommerce-microservices/internal/service"
	"github.com/fdg312/ecommerce-microservices/pkg/auth_v1"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	google_grpc "google.golang.org/grpc"
)

func main() {
	err := godotenv.Load("database.env")
	if err != nil {
		log.Fatalf("Can`t get env for database: %v\n", err)
	}
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPass := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DB")
	postgresPort := 5432
	postgresHost := "localhost"
	postgresUrl := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", postgresUser, postgresPass, postgresHost, postgresPort, postgresDb)
	conn, err := pgx.Connect(context.Background(), postgresUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	repo := repository.NewUserRepository(conn)
	s := service.NewAuthService(repo)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("couldn`t connect to tcp: %v", err)
	}

	ser := grpc.NewAuthServer(s)

	inter := grpc.NewInterceptorManager(*s)

	grpcServer := google_grpc.NewServer(google_grpc.UnaryInterceptor(inter.AuthInterceptor))
	auth_v1.RegisterAuthServiceServer(grpcServer, ser)
	grpcServer.Serve(lis)
}
