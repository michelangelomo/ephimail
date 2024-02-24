package main

import (
	"log"
	"os"
	"temp/server"

	"github.com/urfave/cli/v2"
)

func main() {
	// create mail server
	m := server.NewMailServer()

	app := &cli.App{
		Name:  "disposable emails",
		Usage: "we",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "address",
				Value:       "127.0.0.1",
				Usage:       "Address where to listen to",
				Category:    "Mail server",
				Destination: &m.Address,
			},
			&cli.IntFlag{
				Name:        "port",
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
			m.Run()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
