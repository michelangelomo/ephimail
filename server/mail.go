package server

import (
	"fmt"
	"io"
	"log"
	"net/mail"
	"strings"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/urfave/cli/v2"
)

type MailServer struct {
	Address        string
	Port           int
	AllowedDomains cli.StringSlice
}

func NewMailServer() *MailServer {
	return &MailServer{}
}

func (m *MailServer) Run() {
	b := &Backend{
		Allowed: m.isAllowed,
	}

	s := smtp.NewServer(b)

	s.Addr = fmt.Sprintf("%s:%d", m.Address, m.Port)
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
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
	Allowed func(string) error
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
	if err := s.Backend.Allowed(s.Recipient); err != nil {
		fmt.Printf("Recipient error: %v\n", err)
		return err
	}
	return nil
}

func (s *Session) Data(r io.Reader) error {
	m, err := mail.ReadMessage(r)
	if err != nil {
		log.Fatal(err)
	}

	header := m.Header
	fmt.Println("Date:", header.Get("Date"))
	fmt.Println("From:", header.Get("From"))
	fmt.Println("To:", header.Get("To"))
	fmt.Println("Subject:", header.Get("Subject"))

	// body, err := io.ReadAll(m.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s", body)
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
