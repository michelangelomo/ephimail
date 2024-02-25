package server

import (
	"fmt"
	"io"
	"log"
	"net/mail"
	"strings"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/michelangelomo/ephimail/internal/redis"
	"github.com/urfave/cli/v2"
)

type MailServer struct {
	Address        string
	Port           int
	AllowedDomains cli.StringSlice

	storage redis.Storage
}

func NewMailServer(storage redis.Storage) *MailServer {
	return &MailServer{
		storage: storage,
	}
}

func (m *MailServer) Run() {
	b := &Backend{
		allowed: m.isAllowed,
		storage: m.storage,
	}

	s := smtp.NewServer(b)

	s.Addr = fmt.Sprintf("%s:%d", m.Address, m.Port)
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("starting mail server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (m MailServer) isAllowed(rcpt string) error {
	// check if address is valid
	address, err := mail.ParseAddress(rcpt)
	if err != nil {
		return err
	}
	// extract domain from address
	comp := strings.Split(address.Address, "@")

	for _, d := range m.AllowedDomains.Value() {
		if d == comp[1] {
			return nil
		}
	}
	return fmt.Errorf("%s is not allowed", comp[1])
}

// The Backend implements SMTP server methods.
type Backend struct {
	allowed func(string) error
	storage redis.Storage
}

// A Session is returned after successful login.
type Session struct {
	Recipient string
	Backend   *Backend
}

// NewSession is called after client greeting (EHLO, HELO).
func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		Backend: b,
	}, nil
}

// AuthPlain implements authentication using SASL PLAIN.
func (s *Session) AuthPlain(username, password string) error {
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.Recipient = to
	if err := s.Backend.allowed(s.Recipient); err != nil {
		fmt.Printf("Recipient error: %v\n", err)
		return err
	}
	return nil
}

func (s *Session) Data(r io.Reader) error {
	m, err := mail.ReadMessage(r)
	if err != nil {
		return fmt.Errorf("can't decode mail")
	}

	header := m.Header
	if header.Get("To") != s.Recipient {
		return fmt.Errorf("recipient mismatch")
	}

	body, err := io.ReadAll(m.Body)
	if err != nil {
		log.Fatal(err)
	}
	// save on redis
	err = s.Backend.storage.StoreEmail(
		s.Recipient,
		string(body),
	)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
