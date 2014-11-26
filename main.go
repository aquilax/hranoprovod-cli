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

	defaultDbFilename       = "food.yaml"
	defaultLogFilename      = "log.yaml"
	defaultResolverMaxDepth = 10
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
			Usage:     "Shows the log register report",
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
				cli.StringFlag{
					Name:  "beginning, b",
					Usage: "Beginning of period",
					Value: "",
				},
				cli.StringFlag{
					Name:  "end, e",
					Usage: "End of period",
					Value: "",
				},
				cli.BoolFlag{
					Name:  "csv",
					Usage: "Export as CSV",
				},
				cli.BoolFlag{
					Name:  "no-color",
					Usage: "Disable color output",
				},
				cli.BoolFlag{
					Name:  "no-totals",
					Usage: "Disable totals",
				},
				cli.IntFlag{
					Name:   "maxdepth",
					Value:  defaultResolverMaxDepth,
					Usage:  "resolve depth",
					EnvVar: "HR_MAXDEPTH",
				},
			},
			Action: func(c *cli.Context) {
				ro := &RegisterOptions{
					DateFormat:       "2006/01/02",
					DbFileName:       c.String("database"),
					LogFileName:      c.String("logfile"),
					ResolverMaxDepth: c.Int("maxdepth"),
					CSV:              c.Bool("csv"),
					Color:            !c.Bool("no-color"),
					Totals:           !c.Bool("no-totals"),
					Beginning:        c.String("beginning"),
					End:              c.String("end"),
				}
				err := ro.Validate()
				if err != nil {
					handleExit(err)
				}
				handleExit(NewHranoprovod().Register(ro))
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
