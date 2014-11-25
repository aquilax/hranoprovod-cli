package main

import (
	"github.com/codegangsta/cli"
	"os"
)

const (
	appName            = "hranoprovod-cli"
	appUsage           = "Lifestyle tracker"
	appVersion         = "2.0.0"
	appAuthor          = "aquilax"
	appEmail           = "aquilax@gmail.com"
	defaultDbFilename  = "food.yaml"
	defaultLogFilename = "log.yaml"
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
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "database, d",
					Value:  defaultDbFilename,
					Usage:  "database file name",
					EnvVar: "HR_DATABASE",
				},
				cli.StringFlag{
					Name:   "logfile, l",
					Value:  defaultLogFilename,
					Usage:  "log file name",
					EnvVar: "HR_LOGFILE",
				},
			},
			Action: func(c *cli.Context) {

				handleExit(NewHranoprovod().Register(c.String("database")))
			},
		},
		{
			Name:  "add",
			Usage: "Adds new item to the log",
			Action: func(c *cli.Context) {
				handleExit(NewHranoprovod().Add(c.Args().First(), c.Args().Get(1)))
			},
		},
		{
			Name:  "api",
			Usage: "Service API commands",
			Subcommands: []cli.Command{
				{
					Name:  "search",
					Usage: "Search for food online",
					Action: func(c *cli.Context) {
						handleExit(NewHranoprovod().Search(c.Args().First()))
					},
				},
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
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
