package main

import (
	"os"

	toolkit "github.com/HarryBird/alfred-dev-toolkit/toolkit"
	alfred "github.com/HarryBird/alfred-toolkit-go"
	"github.com/urfave/cli"
)

func main() {
	al, err := alfred.NewAlfred("alfred dev toolkit")
	if err != nil {
		os.Stdout.WriteString("Error: alfred toolkit init fail\n")
		os.Stdout.WriteString("Reason: " + err.Error())
		os.Exit(-1)
	}

	app := cli.NewApp()
	app.Name = "alfred-dev-toolkit"
	app.Usage = "Alfred Workflow To Help Developers' Daily Works"
	app.Action = func(ctx *cli.Context) {
		os.Stdout.WriteString(`NAME:
  alfred-dev-toolkit - Alfred Workflow To Help Developers' Daily Work

  Enter "alfred-dev-toolkit help" for more information`)
	}

	app.Commands = []cli.Command{
		{
			Name:        "ping",
			Usage:       "alfred-dev-toolkit ping <address>",
			Description: "ICMP Ping",
			Action: func(ctx *cli.Context) {
				toolkit.PingAction(ctx, al)
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		os.Stdout.WriteString("Error: console run fail\n")
		os.Stdout.WriteString("Reason: " + err.Error())
		os.Exit(-1)
	}
}
