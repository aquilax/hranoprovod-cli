package cmd

import (
	"fmt"
	"os/user"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type optionLoader = func(*cli.Context) (*app.Options, error)

func GetApp() *cli.App {
	ol := func(c *cli.Context) (*app.Options, error) {
		o := app.NewOptions()
		err := o.Load(c)
		return o, err
	}

	a := &cli.App{
		Name:    app.Name,
		Usage:   app.Usage,
		Version: fmt.Sprintf("%v, commit %v, built at %v", version, commit, date),
	}

	a.Flags = []cli.Flag{
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
			Value:   app.DefaultDbFilename,
			Usage:   "optional database file name `FILE`",
			EnvVars: []string{"HR_DATABASE"},
		},
		&cli.StringFlag{
			Name:    "logfile",
			Aliases: []string{"l"},
			Value:   app.DefaultLogFilename,
			Usage:   "log file name `FILE`",
			EnvVars: []string{"HR_LOGFILE"},
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   getDefaultFileName(app.ConfigFileName),
			Usage:   "Configuration file `FILE`",
			EnvVars: []string{"HR_CONFIG"},
		},
		&cli.StringFlag{
			Name:    "date-format",
			Value:   "2006/01/02",
			Usage:   "Date format for parsing and printing dates `DATE_FORMAT`",
			EnvVars: []string{"HR_DATE_FORMAT"},
		},
		&cli.IntFlag{
			Name:    "maxdepth",
			Value:   app.DefaultResolverMaxDepth,
			Usage:   "Resolve depth `DEPTH`",
			EnvVars: []string{"HR_MAXDEPTH"},
		},
		&cli.BoolFlag{
			Name:  "no-color",
			Usage: "Disable color output",
		},
	}
	a.Commands = []*cli.Command{
		newRegisterCommand(ol),
		newBalanceCommand(ol),
		newLintCommand(ol),
		newReportCommand(ol),
		newCsvCommand(ol),
		newStatsCommand(ol),
		newSummaryCommand(ol),
		newGenCommand(ol, a),
		newPrintCommand(ol),
	}
	return a
}

func getDefaultFileName(fillePath string) string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.HomeDir + fillePath
}
