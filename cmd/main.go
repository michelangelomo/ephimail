// cmd/main.go
package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/michelangelomo/ephimail/internal/redis"
	"github.com/michelangelomo/ephimail/server"

	"github.com/urfave/cli/v2"
)

var storage *redis.RedisStorage
var mail *server.MailServerWithWebSocket
var web *server.WebServerWithWebSocket

func main() {
	// Initialize redis storage
	storage = redis.NewStorage()

	// Initialize web server with WebSocket support
	web = server.NewWebServerWithWebSocket(storage, []string{})

	// Initialize mail server with WebSocket support
	mail = server.NewMailServerWithWebSocket(storage, web)

	app := &cli.App{
		Name:  "ephimail",
		Usage: "all-in-one disposable email service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "redis-address",
				Value:       "127.0.0.1",
				EnvVars:     []string{"REDIS_ADDRESS"},
				Category:    "Redis",
				Destination: &storage.Address,
			},
			&cli.IntFlag{
				Name:        "redis-port",
				Value:       6379,
				EnvVars:     []string{"REDIS_PORT"},
				Category:    "Redis",
				Destination: &storage.Port,
			},
			&cli.StringFlag{
				Name:        "mail-address",
				Value:       "127.0.0.1",
				EnvVars:     []string{"MAIL_ADDRESS"},
				Usage:       "Address where to listen to",
				Category:    "Mail server",
				Destination: &mail.Address,
			},
			&cli.IntFlag{
				Name:        "mail-port",
				Value:       25,
				EnvVars:     []string{"MAIL_PORT"},
				Usage:       "Port where to bind to",
				Category:    "Mail server",
				Destination: &mail.Port,
			},
			&cli.StringSliceFlag{
				Name:        "allow-domain",
				Required:    true,
				EnvVars:     []string{"ALLOWED_DOMAINS"},
				Usage:       "Allowed recipient domains list (comma-separated)",
				Category:    "Mail server",
				Destination: &mail.AllowedDomains,
			},
			&cli.StringFlag{
				Name:        "web-address",
				Value:       "127.0.0.1",
				EnvVars:     []string{"WEB_ADDRESS"},
				Usage:       "Address where to listen to",
				Category:    "Web server",
				Destination: &web.Address,
			},
			&cli.IntFlag{
				Name:        "web-port",
				Value:       80,
				EnvVars:     []string{"WEB_PORT"},
				Usage:       "Port where to bind to",
				Category:    "Web server",
				Destination: &web.Port,
			},
			&cli.IntFlag{
				Name:     "email-ttl",
				Value:    24,
				EnvVars:  []string{"EMAIL_TTL"},
				Usage:    "Email time-to-live in hours (0 for no expiration)",
				Category: "Storage",
			},
			&cli.IntFlag{
				Name:     "max-reservation-days",
				Value:    7,
				EnvVars:  []string{"MAX_RESERVATION_DAYS"},
				Usage:    "Maximum reservation time in days",
				Category: "Web server",
			},
		},
		Action: func(c *cli.Context) error {
			var wg sync.WaitGroup

			// Set email TTL if provided
			if c.Int("email-ttl") > 0 {
				storage.EmailTTL = time.Duration(c.Int("email-ttl")) * time.Hour
			}

			// Connect to redis
			storage.Connect()

			// Set allowed domains for web server
			web.SetAllowedDomains(mail.AllowedDomains.Value())

			// Start web server with WebSocket support
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := web.Run(); err != nil {
					log.Fatal(err)
				}
			}()

			// Start mail server with WebSocket support
			wg.Add(1)
			go func() {
				defer wg.Done()
				mail.Run()
			}()

			wg.Wait()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
