// internal/redis/reservation.go
package redis

import (
	"fmt"
	"strconv"
	"time"
)

// ReservationDuration represents the possible durations for mailbox reservation
type ReservationDuration string

const (
	OneHour ReservationDuration = "1h"
	OneDay  ReservationDuration = "24h"
	OneWeek ReservationDuration = "168h"
)

// Reservation represents a mailbox reservation
type Reservation struct {
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
	PublicKey string    `json:"public_key,omitempty"` // Optional for E2E encryption
	Encrypted bool      `json:"encrypted"`
}

// ReserveMailbox reserves a mailbox for a specific duration
func (r *RedisStorage) ReserveMailbox(email string, duration ReservationDuration, publicKey string) (*Reservation, error) {
	// Check if mailbox is already reserved
	exists, err := r.IsMailboxReserved(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check reservation: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("mailbox %s is already reserved", email)
	}

	// Parse duration
	parsedDuration, err := time.ParseDuration(string(duration))
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %w", err)
	}

	// Create reservation
	expiresAt := time.Now().Add(parsedDuration)
	reservation := &Reservation{
		Email:     email,
		ExpiresAt: expiresAt,
		PublicKey: publicKey,
		Encrypted: publicKey != "",
	}

	// Store reservation in Redis
	key := fmt.Sprintf("reservation:%s", email)
	err = r.Client.HSet(
		r.GetContext(),
		key,
		map[string]interface{}{
			"email":      reservation.Email,
			"expires_at": reservation.ExpiresAt.Unix(),
			"public_key": reservation.PublicKey,
			"encrypted":  reservation.Encrypted,
		},
	).Err()

	if err != nil {
		return nil, fmt.Errorf("failed to save reservation: %w", err)
	}

	// Set expiration for the reservation
	err = r.Client.Expire(r.GetContext(), key, parsedDuration).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set expiration: %w", err)
	}

	return reservation, nil
}

// IsMailboxReserved checks if a mailbox is reserved
func (r *RedisStorage) IsMailboxReserved(email string) (bool, error) {
	key := fmt.Sprintf("reservation:%s", email)
	exists, err := r.Client.Exists(r.GetContext(), key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// GetReservation returns the reservation for a mailbox
func (r *RedisStorage) GetReservation(email string) (*Reservation, error) {
	key := fmt.Sprintf("reservation:%s", email)
	data, err := r.Client.HGetAll(r.GetContext(), key).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	// Parse Unix timestamp
	expiresAtUnix, err := strconv.ParseInt(data["expires_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration timestamp: %w", err)
	}

	reservation := &Reservation{
		Email:     data["email"],
		ExpiresAt: time.Unix(expiresAtUnix, 0),
		PublicKey: data["public_key"],
		Encrypted: data["encrypted"] == "1" || data["encrypted"] == "true",
	}

	return reservation, nil
}

// DeleteReservation deletes a mailbox reservation
func (r *RedisStorage) DeleteReservation(email string) error {
	key := fmt.Sprintf("reservation:%s", email)
	return r.Client.Del(r.GetContext(), key).Err()
}
