// server/mail_encryption.go
package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/michelangelomo/ephimail/internal/encryption"
	"github.com/michelangelomo/ephimail/internal/redis"
)

// EncryptingBackend extends the Backend to handle encrypted emails
type EncryptingBackend struct {
	Backend
	storageWithEncryption *redis.RedisStorage
	webSocketHub          *WebSocketHub
}

// NewEncryptingBackend creates a new encrypting backend
func NewEncryptingBackend(allowed func(string) error, storage *redis.RedisStorage, hub *WebSocketHub) *EncryptingBackend {
	return &EncryptingBackend{
		Backend: Backend{
			allowed: allowed,
			storage: storage,
		},
		storageWithEncryption: storage,
		webSocketHub:          hub,
	}
}

// EncryptingSession extends the Session to handle encrypted emails
type EncryptingSession struct {
	Session
	storageWithEncryption *redis.RedisStorage
	webSocketHub          *WebSocketHub
}

// NewSession creates a new session with encryption support
func (b *EncryptingBackend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &EncryptingSession{
		Session: Session{
			Backend: &b.Backend,
		},
		storageWithEncryption: b.storageWithEncryption,
		webSocketHub:          b.webSocketHub,
	}, nil
}

// Data handles incoming email data with encryption support
func (s *EncryptingSession) Data(r io.Reader) error {
	// Read the email data
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Check if the recipient's mailbox is reserved and encrypted
	reservation, err := s.storageWithEncryption.GetReservation(s.Recipient)
	if err != nil {
		// Fall back to standard behavior if error occurs
		return s.Session.Data(bytes.NewReader(b))
	}

	// If mailbox is reserved and uses encryption, encrypt the email
	if reservation != nil && reservation.Encrypted && reservation.PublicKey != "" {
		// Encrypt the email body with the recipient's public key
		encryptedBody, err := encryption.EncryptWithSymmetricKey(string(b), reservation.PublicKey)
		if err != nil {
			// If encryption fails, still deliver the email unencrypted
			// This is a trade-off between security and reliability
			return s.Session.Data(bytes.NewReader(b))
		}

		// Save the encrypted email
		err = s.Backend.storage.StoreEmail(s.Recipient, encryptedBody)
		if err != nil {
			return err
		}

		// Notify WebSocket clients if available
		if s.webSocketHub != nil {
			// We'd extract a message ID here in a real implementation
			s.webSocketHub.NotifyNewEmail(s.Recipient, "")
		}

		return nil
	}

	// If not encrypted, use standard behavior
	err = s.Session.Data(bytes.NewReader(b))

	// Notify WebSocket clients after successful delivery
	if err == nil && s.webSocketHub != nil {
		s.webSocketHub.NotifyNewEmail(s.Recipient, "")
	}

	return err
}

// RunWithEncryption starts a mail server with encryption support
func (m *MailServer) RunWithEncryption(storage *redis.RedisStorage, wsHub *WebSocketHub) {
	b := NewEncryptingBackend(m.isAllowed, storage, wsHub)

	s := smtp.NewServer(b)

	s.Addr = fmt.Sprintf("%s:%d", m.Address, m.Port)
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting mail server with encryption at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
