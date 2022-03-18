package cmd

import (
	"fmt"
	"os/user"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/aquilax/hranoprovod-cli/v2/parser"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func GetApp() *cli.App {
	u := NewCmdUtils()

	a := &cli.App{
		Name:        app.Name,
		Usage:       app.Usage,
		Description: app.Description,
		Version:     fmt.Sprintf("%v, commit %v, built at %v", version, commit, date),
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
			Value:   parser.DefaultDateFormat,
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
		&cli.BoolFlag{
			Name:  "no-database",
			Usage: "Disables loading the database (even if database filename is set)",
			Value: false,
		},
	}
	a.Commands = []*cli.Command{
		newRegisterCommand(u),
		newBalanceCommand(u),
		newLintCommand(u, app.Lint),
		newReportCommand(u),
		newCSVCommand(u),
		newStatsCommand(u),
		newSummaryCommand(u),
		newGenCommand(u, a),
		newPrintCommand(u, app.Print),
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
