package main

import (
	"os"

	"github.com/urfave/cli"
)

const (
	appName    = "hranoprovod-cli"
	appUsage   = "Lifestyle tracker"
	appVersion = "2.1.3"
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
				cli.BoolFlag{
					Name:  "group-food, g",
					Usage: "Single element grouped by food",
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
					Name:  "totals-only",
					Usage: "Disable totals",
				},
				cli.BoolFlag{
					Name:  "unresolved",
					Usage: "Show unresolved elements only",
				},
				cli.IntFlag{
					Name:   "maxdepth",
					Value:  defaultResolverMaxDepth,
					Usage:  "Resolve depth",
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
			Name:      "balance",
			ShortName: "bal",
			Usage:     "Shows food balance as tree",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "begin, b",
					Usage: "Beginning of period",
				},
				cli.StringFlag{
					Name:  "end, e",
					Usage: "End of period",
				},
				cli.IntFlag{
					Name:   "maxdepth",
					Value:  defaultResolverMaxDepth,
					Usage:  "Resolve depth",
					EnvVar: "HR_MAXDEPTH",
				},
				cli.BoolFlag{
					Name:  "collapse-last",
					Usage: "Collapses last dimension",
				},
				cli.BoolFlag{
					Name:  "collapse, c",
					Usage: "Collapses sole branches",
				},
				cli.StringFlag{
					Name:  "single-element, s",
					Usage: "Show only single element",
				},
			},
			Action: func(c *cli.Context) {
				o := NewOptions()
				if err := o.Load(c); err != nil {
					handleExit(err)
				}
				handleExit(NewHranoprovod(o).Balance())
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
		{
			Name:  "report",
			Usage: "Generates report from the database file",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "desc",
					Usage: "Descending order",
				},
			},
			Action: func(c *cli.Context) {
				o := NewOptions()
				if err := o.Load(c); err != nil {
					handleExit(err)
				}
				handleExit(NewHranoprovod(o).Report(c.Args().First(), c.IsSet("desc")))
			},
		},
		{
			Name:  "csv",
			Usage: "Generates csv log export",
			Action: func(c *cli.Context) {
				o := NewOptions()
				if err := o.Load(c); err != nil {
					handleExit(err)
				}
				handleExit(NewHranoprovod(o).CSV())
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
