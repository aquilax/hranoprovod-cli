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
	app.Flags = []cli.Flag{
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
			Name:   "config, c",
			Value:  GetDefaultFileName(),
			Usage:  "Configuration file",
			EnvVar: "HR_CONFIG",
		},
		cli.StringFlag{
			Name:   "date-format",
			Value:  "2006/01/02",
			Usage:  "Date format for parsing and printing dates",
			EnvVar: "HR_DATE_FORMAT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "register",
			ShortName: "reg",
			Usage:     "Shows the log register report",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "begin, b",
					Usage: "Beginning of period",
				},
				cli.StringFlag{
					Name:  "end, e",
					Usage: "End of period",
				},
				cli.StringFlag{
					Name:  "single-food, f",
					Usage: "Show only single element",
				},
				cli.StringFlag{
					Name:  "single-element, s",
					Usage: "Show only single element",
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
				cli.BoolFlag{
					Name:  "unresolved",
					Usage: "Show unresolved elements only",
				},
				cli.IntFlag{
					Name:   "maxdepth",
					Value:  defaultResolverMaxDepth,
					Usage:  "resolve depth",
					EnvVar: "HR_MAXDEPTH",
				},
			},
			Action: func(c *cli.Context) {
				o := NewOptions()
				if err := o.Load(c); err != nil {
					handleExit(err)
				}
				handleExit(NewHranoprovod(o).Register())
			},
		},
		{
			Name:  "add",
			Usage: "Adds new item to the log",
			Action: func(c *cli.Context) {
				o := NewOptions()
				if err := o.Load(c); err != nil {
					handleExit(err)
				}
				handleExit(NewHranoprovod(o).Add(c.Args().First(), c.Args().Get(1)))
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
						o := NewOptions()
						if err := o.Load(c); err != nil {
							handleExit(err)
						}
						handleExit(NewHranoprovod(o).Search(c.Args().First()))
					},
				},
			},
		},
		{
			Name:  "lint",
			Usage: "Lints file",
			Action: func(c *cli.Context) {
				o := NewOptions()
				if err := o.Load(c); err != nil {
					handleExit(err)
				}
				handleExit(NewHranoprovod(o).Lint(c.Args().First()))
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
