package main

import (
	"fmt"
	"os/user"
	"time"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/balance"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/csv"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/gen"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/lint"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/options"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/print"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/register"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/report"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/stats"
	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/summary"
	"github.com/aquilax/hranoprovod-cli/lib/parser/v3"
	"github.com/aquilax/hranoprovod-cli/lib/resolver/v3"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// GetApp returns a cli app
func GetApp() *cli.App {
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
			Name:    "end",
			Aliases: []string{"e"},
			Usage:   "End of period `DATE`",
		},
		&cli.StringFlag{
			Name:  "today",
			Usage: "Overwrite today's date `DATE`",
		},
		&cli.StringFlag{
			Name:    "database",
			Aliases: []string{"d"},
			Value:   options.DefaultDbFilename,
			Usage:   "optional database file name `FILE`",
			EnvVars: []string{"HR_DATABASE"},
		},
		&cli.StringFlag{
			Name:    "logfile",
			Aliases: []string{"l"},
			Value:   options.DefaultLogFilename,
			Usage:   "log file name `FILE`",
			EnvVars: []string{"HR_LOGFILE"},
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   getDefaultFileName(options.ConfigFileName),
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
			Value:   resolver.DefaultMaxDepth,
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
		register.Command(),
		balance.Command(),
		lint.Command(),
		report.Command(),
		csv.Command(),
		stats.Command(),
		summary.Command(),
		print.Command(),
		gen.Command(a),
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
