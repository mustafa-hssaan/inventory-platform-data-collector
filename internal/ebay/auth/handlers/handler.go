package handlers

import (
	"encoding/json"
	"fmt"
	"inventory-platform-data-collector/internal/ebay/auth/models"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) RegisterConfigAndStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Register config
	h.service.RegisterConfig(req.UserID, &req.Config)

	// Get auth URL
	state := req.UserID
	authURL, err := h.service.GetAuthURL(req.UserID, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return auth URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"auth_url": authURL,
		"user_id":  req.UserID,
	})
}
func (h *Handler) CompleteAuth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	if err := h.service.ExchangeCodeForToken(r.Context(), state, code); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get token immediately
	token, err := h.service.GetAccessToken(r.Context(), state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return complete auth info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":       "success",
		"user_id":      state,
		"access_token": token,
	})
}
func (h *Handler) RegisterConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.service.RegisterConfig(req.UserID, &req.Config)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	state := userID // use a secure state generation method
	authURL, err := h.service.GetAuthURL(userID, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}
	completeURL := fmt.Sprintf("/auth/complete?code=%s&state=%s", code, state)
	http.Redirect(w, r, completeURL, http.StatusTemporaryRedirect)

	// if err := h.service.ExchangeCodeForToken(r.Context(), state, code); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"status": "success",
	// 	"userId": state,
	// })
}

func (h *Handler) GetToken(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	token, err := h.service.GetAccessToken(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": token})
}
