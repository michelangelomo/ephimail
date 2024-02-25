package main

import (
	"log"
	"os"
	"sync"

	"github.com/michelangelomo/ephimail/internal/redis"
	"github.com/michelangelomo/ephimail/server"

	"github.com/urfave/cli/v2"
)

var mail *server.MailServer
var storage *redis.RedisStorage
var web *server.WebServer

func main() {
	// initialize redis storage
	storage = redis.NewStorage()
	// initialize mail
	mail = server.NewMailServer(storage)
	// initialize web
	web = server.NewWebServer(storage)

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
				Destination: &mail.Address,
			},
			&cli.IntFlag{
				Name:        "mail-port",
				Value:       25,
				Usage:       "Port where to bind to",
				Category:    "Mail server",
				Destination: &mail.Port,
			},
			&cli.StringSliceFlag{
				Name:        "allow-domain",
				Required:    true,
				Usage:       "Allowed recipient domains list",
				Category:    "Mail server",
				Destination: &mail.AllowedDomains,
			},
			&cli.StringFlag{
				Name:        "web-address",
				Value:       "127.0.0.1",
				Usage:       "Address where to listen to",
				Category:    "Web server",
				Destination: &web.Address,
			},
			&cli.IntFlag{
				Name:        "web-port",
				Value:       80,
				Usage:       "Port where to bind to",
				Category:    "Web server",
				Destination: &web.Port,
			},
		},
		Action: func(*cli.Context) error {
			var wg sync.WaitGroup
			// connect to redis
			storage.Connect()
			// start mail server
			wg.Add(1)
			go func() {
				defer wg.Done()
				mail.Run()
			}()
			// start web server
			wg.Add(1)
			go func() {
				defer wg.Done()
				web.Run()
			}()

			wg.Wait()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
