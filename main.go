package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jared-prime/weather/conditions"
	"github.com/urfave/cli"
)

func main() {
	var mode string
	app := cli.NewApp()
	app.Name = "weather"
	app.Usage = "report the weather"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "mode",
			Value:       "client",
			Usage:       "run the app in server or client mode",
			Destination: &mode,
		},
	}

	app.Action = func(c *cli.Context) error {
		switch mode {
		case "client":
			client := conditions.Client{}
			return client.Start()
		case "server":
			return conditions.StartServer()
		default:
			return fmt.Errorf("unknown runner mode: %s", mode)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
