package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	appName    = "hranoprovod-cli"
	appUsage   = "Lifestyle tracker"
	appVersion = "2.2.2"
	appAuthor  = "aquilax"
	appEmail   = "aquilax@gmail.com"

	defaultDbFilename       = "food.yaml"
	defaultLogFilename      = "log.yaml"
	defaultResolverMaxDepth = 10
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = appUsage
	app.Version = fmt.Sprintf("%v, commit %v, built at %v", version, commit, date)
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "begin",
			Aliases: []string{"b"},
			Usage:   "Beginning of period `DATE`",
		},
		&cli.StringFlag{
			Name:    "end, e",
			Aliases: []string{"e"},
			Usage:   "End of period `DATE`",
		},
		&cli.StringFlag{
			Name:    "database",
			Aliases: []string{"d"},
			Value:   defaultDbFilename,
			Usage:   "database file name `FILE`",
			EnvVars: []string{"HR_DATABASE"},
		},
		&cli.StringFlag{
			Name:    "logfile",
			Aliases: []string{"l"},
			Value:   defaultLogFilename,
			Usage:   "log file name `FILE`",
			EnvVars: []string{"HR_LOGFILE"},
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   GetDefaultFileName(),
			Usage:   "Configuration file `FILE`",
			EnvVars: []string{"HR_CONFIG"},
		},
		&cli.StringFlag{
			Name:    "date-format",
			Value:   "2006/01/02",
			Usage:   "Date format for parsing and printing dates",
			EnvVars: []string{"HR_DATE_FORMAT"},
		},
		&cli.IntFlag{
			Name:    "maxdepth",
			Value:   defaultResolverMaxDepth,
			Usage:   "Resolve depth `DEPTH`",
			EnvVars: []string{"HR_MAXDEPTH"},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "register",
			Aliases: []string{"reg"},
			Usage:   "Shows the log register report",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "begin",
					Aliases: []string{"b"},
					Usage:   "Beginning of period `DATE`",
				},
				&cli.StringFlag{
					Name:    "end",
					Aliases: []string{"e"},
					Usage:   "End of period `DATE`",
				},
				&cli.StringFlag{
					Name:    "single-food",
					Aliases: []string{"f"},
					Usage:   "Show only single food",
				},
				&cli.BoolFlag{
					Name:    "group-food",
					Aliases: []string{"g"},
					Usage:   "Single element grouped by food",
				},
				&cli.StringFlag{
					Name:    "single-element",
					Aliases: []string{"s"},
					Usage:   "Show only single element",
				},
				&cli.BoolFlag{
					Name:  "csv",
					Usage: "Export as CSV",
				},
				&cli.BoolFlag{
					Name:  "no-color",
					Usage: "Disable color output",
				},
				&cli.BoolFlag{
					Name:  "no-totals",
					Usage: "Disable totals",
				},
				&cli.BoolFlag{
					Name:  "totals-only",
					Usage: "Show only totals",
				},
				&cli.BoolFlag{
					Name:  "shorten",
					Usage: "Shorten longer strings",
				},
				&cli.BoolFlag{
					Name:  "new-reg-reporter",
					Usage: "Use the new reg reporter",
				},
				&cli.BoolFlag{
					Name:  "unresolved",
					Usage: "Deprecated: Show unresolved elements only (moved to 'report unresolved')",
				},
			},
			Action: func(c *cli.Context) error {
				o := NewOptions()
				if err := o.Load(c); err != nil {
					return err
				}
				return NewHranoprovod(o).Register()
			},
		},
		{
			Name:    "balance",
			Aliases: []string{"bal"},
			Usage:   "Shows food balance as tree",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "begin",
					Aliases: []string{"b"},
					Usage:   "Beginning of period",
				},
				&cli.StringFlag{
					Name:    "end",
					Aliases: []string{"e"},
					Usage:   "End of period",
				},
				&cli.BoolFlag{
					Name:  "collapse-last",
					Usage: "Collapses last dimension",
				},
				&cli.BoolFlag{
					Name:    "collapse",
					Aliases: []string{"c"},
					Usage:   "Collapses sole branches",
				},
				&cli.StringFlag{
					Name:    "single-element, s",
					Aliases: []string{"s"},
					Usage:   "Show only single element",
				},
			},
			Action: func(c *cli.Context) error {
				o := NewOptions()
				if err := o.Load(c); err != nil {
					return err
				}
				return NewHranoprovod(o).Balance()
			},
		},
		{
			Name:  "add",
			Usage: "Adds new item to the log",
			Action: func(c *cli.Context) error {
				o := NewOptions()
				if err := o.Load(c); err != nil {
					return err
				}
				return NewHranoprovod(o).Add(c.Args().First(), c.Args().Get(1))
			},
		},
		{
			Name:  "api",
			Usage: "Service API commands",
			Subcommands: []*cli.Command{
				{
					Name:  "search",
					Usage: "Search for food online",
					Action: func(c *cli.Context) error {
						o := NewOptions()
						if err := o.Load(c); err != nil {
							return err
						}
						return NewHranoprovod(o).Search(c.Args().First())
					},
				},
			},
		},
		{
			Name:  "lint",
			Usage: "Lints file",
			Action: func(c *cli.Context) error {
				o := NewOptions()
				if err := o.Load(c); err != nil {
					return err
				}
				return NewHranoprovod(o).Lint(c.Args().First())
			},
		},
		{
			Name:  "report",
			Usage: "Generates various reports",
			Subcommands: []*cli.Command{
				{
					Name:  "element-total",
					Usage: "Generates total sum for element grouped by food",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "desc",
							Usage: "Descending order",
						},
					},
					Action: func(c *cli.Context) error {
						o := NewOptions()
						if err := o.Load(c); err != nil {
							return err
						}
						return NewHranoprovod(o).ReportElement(c.Args().First(), c.IsSet("desc"))
					},
				},
				{
					Name:  "unresolved",
					Usage: "Print list of unresolved elements",
					Action: func(c *cli.Context) error {
						o := NewOptions()
						if err := o.Load(c); err != nil {
							return err
						}
						return NewHranoprovod(o).ReportUnresolved()
					},
				},
			},
		},
		{
			Name:  "csv",
			Usage: "Generates csv exports",
			Subcommands: []*cli.Command{
				{
					Name:  "log",
					Usage: "Exports the log file as CSV",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "begin",
							Aliases: []string{"b"},
							Usage:   "Beginning of period `DATE`",
						},
						&cli.StringFlag{
							Name:    "end",
							Aliases: []string{"e"},
							Usage:   "End of period `DATE`",
						},
					},
					Action: func(c *cli.Context) error {
						o := NewOptions()
						if err := o.Load(c); err != nil {
							return err
						}
						return NewHranoprovod(o).CSV()
					},
				},
			},
		},
		{
			Name:  "stats",
			Usage: "Provide summary information",
			Action: func(c *cli.Context) error {
				o := NewOptions()
				if err := o.Load(c); err != nil {
					return err
				}
				return NewHranoprovod(o).Stats()
			},
		},
	}
	app.Run(os.Args)
}
