package cmd

import (
	"fmt"
	"os/user"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/app"
	"github.com/aquilax/hranoprovod-cli/v2/parser"
	"github.com/urfave/cli/v2"
)

const (
	configFileName = "/.hranoprovod/config"

	defaultDbFilename       = "food.yaml"
	defaultLogFilename      = "log.yaml"
	defaultResolverMaxDepth = 10
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// GetApp returns a cli app
func GetApp() *cli.App {
	u := newCmdUtils()

	a := &cli.App{
		Name:        "hranoprovod-cli",
		Usage:       "Diet tracker for the command line",
		Description: "A command line tool to keep log of diet and exercise in text files",
		Authors:     []*cli.Author{{Name: "aquilax", Email: "aquilax@gmail.com"}},
		Version:     fmt.Sprintf("%v, commit %v, built at %v", version, commit, date),
		Compiled:    time.Now(),
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
			Value:   defaultDbFilename,
			Usage:   "optional database file name `FILE`",
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
			Value:   getDefaultFileName(configFileName),
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
			Value:   defaultResolverMaxDepth,
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
		newRegisterCommand(u, app.Register),
		newBalanceCommand(u, app.Balance),
		newLintCommand(u, app.Lint),
		newReportCommand(u),
		newCSVCommand(u),
		newStatsCommand(u, app.Stats),
		newSummaryCommand(u, app.Summary),
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
