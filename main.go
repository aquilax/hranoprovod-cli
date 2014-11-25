package main

import (
	"github.com/codegangsta/cli"
	"os"
)

const (
	appName    = "hranoprovod-cli"
	appUsage   = "Lifestyle tracker"
	appVersion = "2.0.0"
	appAuthor  = "aquilax"
	appEmail   = "aquilax@gmail.com"
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = appUsage
	app.Version = appVersion
	app.Author = appAuthor
	app.Email = appEmail
	app.Commands = []cli.Command{
		{
			Name:      "register",
			ShortName: "reg",
			Usage:     "Shows the register",
			Action: func(c *cli.Context) {
				handleExit(NewHranoprovod().Register())
			},
		},
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "Adds new item to the log",
			Action: func(c *cli.Context) {
				handleExit(NewHranoprovod().Add(c.Args().First(), c.Args().Get(1)))
			},
		},
		{
			Name:  "search",
			Usage: "Search for food online",
			Action: func(c *cli.Context) {
				handleExit(NewHranoprovod().Search(c.Args().First()))
			},
		},
		{
			Name:  "lint",
			Usage: "Lints file",
			Action: func(c *cli.Context) {
				handleExit(NewHranoprovod().Lint(c.Args().First()))
			},
		},
	}
	app.Run(os.Args)
}

func handleExit(err error) {

}
