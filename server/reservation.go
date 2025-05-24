// server/reservation.go
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/michelangelomo/ephimail/internal/redis"
)

const (
	reservationEnabled bool = false
)

// ReservationRequest represents the request to reserve a mailbox
type ReservationRequest struct {
	Email     string `json:"email"`
	Duration  string `json:"duration"`
	PublicKey string `json:"public_key,omitempty"`
}

// ReservationResponse represents the response for a mailbox reservation
type ReservationResponse struct {
	Email      string    `json:"email"`
	ExpiresAt  time.Time `json:"expires_at"`
	Encrypted  bool      `json:"encrypted"`
	ReservedAt time.Time `json:"reserved_at"`
	URL        string    `json:"url,omitempty"`
}

// RegisterReservationHandlers registers the reservation handlers
func (w *WebServer) RegisterReservationHandlers(router *mux.Router) {
	router.HandleFunc("/api/inbox/reserve", w.reserveMailbox).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/inbox/{email}/reservation", w.getReservation).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/inbox/{email}/reservation", w.deleteReservation).Methods("DELETE", "OPTIONS")
}

// reserveMailbox handles the reservation of a mailbox
func (w *WebServer) reserveMailbox(rw http.ResponseWriter, r *http.Request) {
	// Enable CORS
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
	if r.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusOK)
		return
	}

	if !reservationEnabled {
		http.Error(rw, "Mailbox reservation is disabled", http.StatusNotImplemented)
		return
	}

	// Parse request
	var req ReservationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Email == "" {
		http.Error(rw, "Email is required", http.StatusBadRequest)
		return
	}

	// Check if duration is valid
	var duration time.Duration
	switch req.Duration {
	case "1h":
		duration = 1 * time.Hour
	case "24h":
		duration = 24 * time.Hour
	case "168h":
		duration = 168 * time.Hour
	default:
		http.Error(rw, "Invalid duration. Allowed values: 1h, 24h, 168h", http.StatusBadRequest)
		return
	}

	// Reserve mailbox
	storage, ok := w.storage.(*redis.RedisStorage)
	if !ok {
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if mailbox is already reserved
	reserved, err := w.isMailboxReserved(req.Email)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to check mailbox status: %s", err), http.StatusInternalServerError)
		return
	}

	if reserved {
		http.Error(rw, "Mailbox is already reserved", http.StatusConflict)
		return
	}

	// Create reservation key
	key := fmt.Sprintf("reservation:%s", req.Email)

	// Set reservation data
	now := time.Now()
	expiresAt := now.Add(duration)
	encrypted := req.PublicKey != ""

	// Store reservation in Redis
	err = storage.Client.HSet(
		storage.GetContext(),
		key,
		map[string]interface{}{
			"email":      req.Email,
			"expires_at": expiresAt.Unix(),
			"public_key": req.PublicKey,
			"encrypted":  encrypted,
			"created_at": now.Unix(),
		},
	).Err()

	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to reserve mailbox: %s", err), http.StatusInternalServerError)
		return
	}

	// Set key expiration
	err = storage.Client.Expire(
		storage.GetContext(),
		key,
		duration,
	).Err()

	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to set reservation expiration: %s", err), http.StatusInternalServerError)
		return
	}

	// Create response
	resp := ReservationResponse{
		Email:      req.Email,
		ExpiresAt:  expiresAt,
		Encrypted:  encrypted,
		ReservedAt: now,
	}

	// If encrypted, create URL with private key placeholder
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	if encrypted {
		resp.URL = fmt.Sprintf("%s://%s/inbox/%s#private_key_goes_here", scheme, r.Host, req.Email)
	} else {
		resp.URL = fmt.Sprintf("%s://%s/inbox/%s", scheme, r.Host, req.Email)
	}

	// Return response
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// isMailboxReserved checks if a mailbox is reserved
func (w *WebServer) isMailboxReserved(email string) (bool, error) {
	storage, ok := w.storage.(*redis.RedisStorage)
	if !ok {
		return false, fmt.Errorf("invalid storage type")
	}

	key := fmt.Sprintf("reservation:%s", email)
	exists, err := storage.Client.Exists(storage.GetContext(), key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

// getReservation handles getting a mailbox reservation
func (w *WebServer) getReservation(rw http.ResponseWriter, r *http.Request) {
	// Enable CORS
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	// Handle preflight request
	if r.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusOK)
		return
	}

	if !reservationEnabled {
		http.Error(rw, "Mailbox reservation is disabled", http.StatusNotImplemented)
		return
	}

	vars := mux.Vars(r)
	email := vars["email"]

	// Check if mailbox is reserved
	storage, ok := w.storage.(*redis.RedisStorage)
	if !ok {
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

	key := fmt.Sprintf("reservation:%s", email)
	data, err := storage.Client.HGetAll(storage.GetContext(), key).Result()
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to get reservation: %s", err), http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		http.Error(rw, "Reservation not found", http.StatusNotFound)
		return
	}

	// Parse data
	encrypted := data["encrypted"] == "1" || data["encrypted"] == "true"
	expiresAtUnix, err := storage.Client.HGet(storage.GetContext(), key, "expires_at").Int64()
	if err != nil {
		http.Error(rw, fmt.Sprintf("Invalid expiration timestamp: %s", err), http.StatusInternalServerError)
		return
	}

	// Create response
	resp := ReservationResponse{
		Email:      email,
		ExpiresAt:  time.Unix(expiresAtUnix, 0),
		Encrypted:  encrypted,
		ReservedAt: time.Now(), // This would ideally come from the stored data
	}

	// Return reservation
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}

// deleteReservation handles deleting a mailbox reservation
func (w *WebServer) deleteReservation(rw http.ResponseWriter, r *http.Request) {
	// Enable CORS
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")

	// Handle preflight request
	if r.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusOK)
		return
	}

	if !reservationEnabled {
		http.Error(rw, "Mailbox reservation is disabled", http.StatusNotImplemented)
		return
	}

	vars := mux.Vars(r)
	email := vars["email"]

	storage, ok := w.storage.(*redis.RedisStorage)
	if !ok {
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

	key := fmt.Sprintf("reservation:%s", email)
	deleted, err := storage.Client.Del(storage.GetContext(), key).Result()
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to delete reservation: %s", err), http.StatusInternalServerError)
		return
	}

	if deleted == 0 {
		http.Error(rw, "Reservation not found", http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
