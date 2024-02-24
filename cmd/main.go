package main

import (
	"log"
	"os"

	r "github.com/michelangelomo/ephimail/internal/redis"
	"github.com/michelangelomo/ephimail/server"

	"github.com/urfave/cli/v2"
)

var m server.MailServer
var redis r.Client

var mailFlags = []cli.Flag{
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
}

var redisFlag = []cli.Flag{
	&cli.StringFlag{
		Name:        "redis-address",
		Value:       "127.0.0.1",
		Category:    "Redis",
		Destination: &redis.Address,
	},
	&cli.IntFlag{
		Name:        "redis-port",
		Value:       6379,
		Category:    "Redis",
		Destination: &redis.Port,
	},
}

func main() {
	app := &cli.App{
		Name:  "ephimail",
		Usage: "all-in-one disposable email service",
		Flags: append(mailFlags, redisFlag...),
		Action: func(*cli.Context) error {
			// connect to redis
			cr, err := redis.Connect()
			if err != nil {
				return err
			}
			m.RedisClient = cr

			// start mail server
			m.Run()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
