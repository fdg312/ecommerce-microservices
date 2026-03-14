package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fdg312/ecommerce-microservices/internal/service"
)

type Auth struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *Auth {
	return &Auth{s}
}

func (h *Auth) Register(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf(err.Error())
		// http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
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

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "not correct credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
