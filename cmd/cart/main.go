package main

import (
	"log"
	"net"

	"github.com/fdg312/ecommerce-microservices/internal/repository"
	"github.com/fdg312/ecommerce-microservices/internal/service"
	"github.com/fdg312/ecommerce-microservices/pkg/auth_v1"
	"github.com/fdg312/ecommerce-microservices/pkg/cart_v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	defer rdb.Close()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("could not listen on port 50052: %v", err)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to auth grpc server: %v", err)
	}
	defer conn.Close()

	authClient := auth_v1.NewAuthServiceClient(conn)

	repo := repository.NewCartRepository(rdb)
	s := service.NewCartService(authClient, repo)

	grpcServer := grpc.NewServer()
	cart_v1.RegisterCartServiceServer(grpcServer, s)

	log.Println("Cart service is running on port 50052...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
