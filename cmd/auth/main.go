package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err_godotenv := godotenv.Load("database.env")
	if err_godotenv != nil {
		log.Fatalf("Can`t get env for database: %v", err_godotenv.Error())
	}
	u := os.Getenv("POSTGRES_USER")
	fmt.Println("START", u)
}
