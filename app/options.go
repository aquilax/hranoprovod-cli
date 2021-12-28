package app

import (
	"errors"
	"os"
	"time"

	"github.com/aquilax/hranoprovod-cli/v2/parser"
	"github.com/aquilax/hranoprovod-cli/v2/reporter"
	"github.com/aquilax/hranoprovod-cli/v2/resolver"
	"github.com/tj/go-naturaldate"
	"github.com/urfave/cli/v2"
	gcfg "gopkg.in/gcfg.v1"
)

type GlobalConfig struct {
	DbFileName  string
	LogFileName string
	DateFormat  string
}

type FilterConfig struct {
	BeginningTime *time.Time
	EndTime       *time.Time
}

// Options contains the options structure
type Options struct {
	GlobalConfig   GlobalConfig
	ResolverConfig resolver.Config `gcfg:"Resolver"`
	ParserConfig   parser.Config   `gcfg:"Parser"`
	ReporterConfig reporter.Config `gcfg:"Reporter"`
	FilterConfig   FilterConfig
}

// NewOptions returns new options structure.
func NewOptions() *Options {
	o := &Options{}
	o.ReporterConfig = reporter.NewDefaultConfig()
	o.ReporterConfig.Color = true
	o.ParserConfig = parser.NewDefaultConfig()
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

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

func (o *Options) populateGlobals(c *cli.Context) {
	if c.IsSet("database") || o.GlobalConfig.DbFileName == "" {
		o.GlobalConfig.DbFileName = c.String("database")
	}

	if c.IsSet("logfile") || o.GlobalConfig.LogFileName == "" {
		o.GlobalConfig.LogFileName = c.String("logfile")
	}

	if c.IsSet("date-format") || o.GlobalConfig.DateFormat == "" {
		o.GlobalConfig.DateFormat = c.String("date-format")
	}
}

func (o *Options) populateLocals(c *cli.Context) {
	o.populateResolver(c)
	o.populateReporter(c)
}

func (o *Options) populateResolver(c *cli.Context) {
	if c.IsSet("maxdepth") || o.ResolverConfig.MaxDepth == 0 {
		o.ResolverConfig.MaxDepth = c.Int("maxdepth")
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
	var err error
	var customTime time.Time
	customTime, err = time.Parse(format, date)
	if err == nil {
		return customTime
	}
	customTime, err = naturaldate.Parse(date, time.Now())
	if err == nil {
		return customTime
	}
	panic(err)
}

func (o *Options) populateReporter(c *cli.Context) {
	for i := len(c.Lineage()) - 1; i >= 0; i-- {
		if c.Lineage()[i].IsSet("csv") {
			o.ReporterConfig.CSV = true
		}
		if c.Lineage()[i].IsSet("no-color") {
			o.ReporterConfig.Color = false
		}

		if c.Lineage()[i].IsSet("collapse-last") {
			o.ReporterConfig.CollapseLast = true
		}

		if c.Lineage()[i].IsSet("collapse") {
			o.ReporterConfig.Collapse = true
		}

		if c.Lineage()[i].IsSet("no-totals") {
			o.ReporterConfig.Totals = false
		}

		if c.Lineage()[i].IsSet("totals-only") {
			o.ReporterConfig.TotalsOnly = true
		}

		if c.Lineage()[i].IsSet("shorten") {
			o.ReporterConfig.ShortenStrings = true
		}

		if c.Lineage()[i].IsSet("use-old-reg-reporter") {
			o.ReporterConfig.UseOldRegReporter = true
		}

		if c.Lineage()[i].IsSet("internal-template-name") {
			o.ReporterConfig.InternalTemplateName = c.Lineage()[i].String("internal-template-name")
		}

		if c.Lineage()[i].IsSet("begin") {
			time := mustGetTime(o.GlobalConfig.DateFormat, c.Lineage()[i].String("begin"))
			o.FilterConfig.BeginningTime = &time
		}
		if c.Lineage()[i].IsSet("end") {
			time := mustGetTime(o.GlobalConfig.DateFormat, c.Lineage()[i].String("end"))
			o.FilterConfig.EndTime = &time
		}
	}
	o.ReporterConfig.Unresolved = c.Bool("unresolved")
	o.ReporterConfig.SingleFood = c.String("single-food")
	o.ReporterConfig.ElementGroupByFood = c.Bool("group-food")
	o.ReporterConfig.SingleElement = c.String("single-element")
}
