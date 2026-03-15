package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fdg312/ecommerce-microservices/pkg/auth_v1"
	google_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	client auth_v1.AuthServiceClient
}

func NewGateway(c auth_v1.AuthServiceClient) *Gateway {
	return &Gateway{c}
}

func (g *Gateway) handleLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req loginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	token, err := g.client.Login(r.Context(), &auth_v1.LoginRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		http.Error(w, "not correct credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token.GetToken()})

}

func (g *Gateway) handleRegister(w http.ResponseWriter, r *http.Request) {
	type registerRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req registerRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = g.client.Register(r.Context(), &auth_v1.RegisterRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func main() {
	conn, err := google_grpc.NewClient("localhost:50051", google_grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to grpc server: %v", err)
	}
	defer conn.Close()

	client := auth_v1.NewAuthServiceClient(conn)
	gateway := NewGateway(client)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", gateway.handleRegister)
	mux.HandleFunc("POST /login", gateway.handleLogin)

	http.ListenAndServe(":8080", mux)
}
