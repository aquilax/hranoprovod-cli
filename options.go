package main

import (
	"errors"
	"os"
	"os/user"
	"time"

	client "github.com/aquilax/hranoprovod-cli/api-client"
	"github.com/aquilax/hranoprovod-cli/parser"
	"github.com/aquilax/hranoprovod-cli/reporter"
	"github.com/urfave/cli/v2"
	gcfg "gopkg.in/gcfg.v1"
)

const (
	optionsFileName = "/.hranoprovod/config"
)

// Options contains the options structure
type Options struct {
	Global struct {
		DbFileName  string
		LogFileName string
		DateFormat  string
	}
	Resolver struct {
		ResolverMaxDepth int
	}
	Parser   parser.Options
	Reporter reporter.Options
	API      client.Options
}

// NewOptions returns new options structure.
func NewOptions() *Options {
	o := &Options{}
	o.Reporter = *reporter.NewDefaultOptions()
	o.Reporter.Color = true
	o.Parser = *parser.NewDefaultOptions()
	o.API = *client.NewDefaultOptions()
	return o
}

// Load loads the settings from config file/command line params/defaults from given context.
func (o *Options) Load(c *cli.Context) error {
	fileName := c.String("config")
	// First try to load the o file
	exists, err := fileExists(fileName)
	if err != nil {
		return err
	}
	// Non existing file passed
	if !exists && c.IsSet("config") {
		return errors.New("File " + fileName + "not found")
	}
	if exists {
		if err := gcfg.ReadFileInto(o, fileName); err != nil {
			return err
		}
	}
	o.populateGlobals(c)
	o.populateLocals(c)
	return nil
}

// GetDefaultFileName returns the default filename for the config file
func GetDefaultFileName() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.HomeDir + optionsFileName
}

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

func (o *Options) populateGlobals(c *cli.Context) {
	if c.IsSet("database") || o.Global.DbFileName == "" {
		o.Global.DbFileName = c.String("database")
	}

	if c.IsSet("logfile") || o.Global.LogFileName == "" {
		o.Global.LogFileName = c.String("logfile")
	}

	if c.IsSet("date-format") || o.Global.DateFormat == "" {
		o.Global.DateFormat = c.String("date-format")
	}
}

func (o *Options) populateLocals(c *cli.Context) {
	o.populateResolver(c)
	o.populateReporter(c)
}

func (o *Options) populateResolver(c *cli.Context) {
	if c.IsSet("maxdepth") || o.Resolver.ResolverMaxDepth == 0 {
		o.Resolver.ResolverMaxDepth = c.Int("maxdepth")
	}
}

func mustGetTime(format string, date string) time.Time {
	if date == "today" {
		return time.Now().Local()
	}
	if date == "yesterday" {
		return time.Now().AddDate(0, 0, -1)
	}
	if date == "last7" {
		return time.Now().AddDate(0, 0, -7)
	}
	if date == "last30" {
		return time.Now().AddDate(0, 0, -30)
	}
	customTime, err := time.Parse(format, date)
	if err != nil {
		panic(err)
	}
	return customTime
}

func (o *Options) populateReporter(c *cli.Context) {
	for i := len(c.Lineage()) - 1; i >= 0; i-- {
		if c.Lineage()[i].IsSet("csv") {
			o.Reporter.CSV = true
		}
		if c.Lineage()[i].IsSet("no-color") {
			o.Reporter.Color = false
		}

		if c.Lineage()[i].IsSet("collapse-last") {
			o.Reporter.CollapseLast = true
		}

		if c.Lineage()[i].IsSet("collapse") {
			o.Reporter.Collapse = true
		}

		if c.Lineage()[i].IsSet("no-totals") {
			o.Reporter.Totals = false
		}

		if c.Lineage()[i].IsSet("totals-only") {
			o.Reporter.TotalsOnly = true
		}

		if c.Lineage()[i].IsSet("shorten") {
			o.Reporter.ShortenStrings = true
		}

		if c.Lineage()[i].IsSet("begin") {
			o.Reporter.BeginningTime = mustGetTime(o.Global.DateFormat, c.Lineage()[i].String("begin"))
			o.Reporter.HasBeginning = true
		}
		if c.Lineage()[i].IsSet("end") {
			o.Reporter.EndTime = mustGetTime(o.Global.DateFormat, c.Lineage()[i].String("end"))
			o.Reporter.HasEnd = true
		}
	}
	o.Reporter.Unresolved = c.Bool("unresolved")
	o.Reporter.SingleFood = c.String("single-food")
	o.Reporter.ElementGroupByFood = c.Bool("group-food")
	o.Reporter.SingleElement = c.String("single-element")
}
