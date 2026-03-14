package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	handler "github.com/fdg312/ecommerce-microservices/internal/handlers"
	"github.com/fdg312/ecommerce-microservices/internal/repository"
	"github.com/fdg312/ecommerce-microservices/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
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
	h := handler.NewAuthHandler(s)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /login", h.Login)

	http.ListenAndServe(":8080", mux)
}
