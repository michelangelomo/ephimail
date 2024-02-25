package main

import (
	"log"
	"os"

	"github.com/michelangelomo/ephimail/internal/redis"
	"github.com/michelangelomo/ephimail/server"

	"github.com/urfave/cli/v2"
)

var m *server.MailServer
var storage *redis.RedisStorage

func main() {
	// initialize redis storage
	storage = redis.NewStorage()

	// initialize mail
	m = server.NewMailServer(storage)

	app := &cli.App{
		Name:  "ephimail",
		Usage: "all-in-one disposable email service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "redis-address",
				Value:       "127.0.0.1",
				Category:    "Redis",
				Destination: &storage.Address,
			},
			&cli.IntFlag{
				Name:        "redis-port",
				Value:       6379,
				Category:    "Redis",
				Destination: &storage.Port,
			},
			&cli.StringFlag{
				Name:        "mail-address",
				Value:       "127.0.0.1",
				Usage:       "Address where to listen to",
				Category:    "Mail server",
				Destination: &m.Address,
			},
			&cli.IntFlag{
				Name:        "mail-port",
				Value:       25,
				Usage:       "Port where to bind to",
				Category:    "Mail server",
				Destination: &m.Port,
			},
			&cli.StringSliceFlag{
				Name:        "allow-domain",
				Required:    true,
				Usage:       "Allowed recipient domains list",
				Category:    "Mail server",
				Destination: &m.AllowedDomains,
			},
		},
		Action: func(*cli.Context) error {
			// connect to redis
			storage.Connect()
			// start mail server
			m.Run()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
