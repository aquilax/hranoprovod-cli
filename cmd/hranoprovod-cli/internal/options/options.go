package options

import (
	"bufio"
	"errors"
	"io"
	"os"
	"time"

	"github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3/internal/reporter"
	"github.com/aquilax/hranoprovod-cli/v3/filter"
	"github.com/aquilax/hranoprovod-cli/v3/parser"
	"github.com/aquilax/hranoprovod-cli/v3/resolver"
	"github.com/tj/go-naturaldate"
	"github.com/urfave/cli/v2"
	gcfg "gopkg.in/gcfg.v1"
)

const (
	ConfigFileName     = "/.hranoprovod/config"
	DefaultDbFilename  = "food.yaml"
	DefaultLogFilename = "log.yaml"
)

type GlobalConfig struct {
	Now         time.Time
	DbFileName  string
	LogFileName string
	DateFormat  string
}

func NewDefaultGlobalConfig() GlobalConfig {
	return GlobalConfig{
		time.Now().Local(),
		DefaultDbFilename,
		DefaultLogFilename,
		parser.DefaultDateFormat,
	}
}

// Options contains the options structure
type Options struct {
	GlobalConfig   GlobalConfig    `gcfg:"Global"`
	ResolverConfig resolver.Config `gcfg:"Resolver"`
	ParserConfig   parser.Config
	ReporterConfig reporter.Config
	FilterConfig   filter.Config
}

// New returns new options structure.
func New() *Options {
	return &Options{
		NewDefaultGlobalConfig(),
		resolver.NewDefaultConfig(),
		parser.NewDefaultConfig(),
		reporter.NewDefaultConfig(),
		filter.NewDefaultConfig(),
	}
}

// Load loads the settings from config file/command line params/defaults from given context.
func (o *Options) Load(c *cli.Context, useConfigFile bool) error {
	if useConfigFile {
		fileName := c.String("config")
		// First try to load the o file
		exists, err := fileExists(fileName)
		if err != nil {
			return err
		}
		// Non existing file passed
		if !exists && c.IsSet("config") {
			return errors.New("File " + fileName + " not found")
		}
		if exists {
			if f, err := os.Open(fileName); err != nil {
				return err
			} else {
				r := bufio.NewReader(f)
				if err := loadFromConfigFile(o, r); err != nil {
					return err
				}
			}
		}
	}
	if err := o.populateGlobals(c); err != nil {
		return err
	}
	o.populateResolver(c)
	o.populateReporter(c)
	if err := o.populateFilter(c); err != nil {
		return err
	}
	return validateOptions(c, o)
}

func loadFromConfigFile(o *Options, r io.Reader) error {
	return gcfg.ReadInto(o, r)
}

func validateOptions(c *cli.Context, o *Options) error {
	// TODO: validate options
	return nil
}

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

func (o *Options) populateGlobals(c *cli.Context) error {
	if !c.IsSet("no-database") && (c.IsSet("database") || o.GlobalConfig.DbFileName == "") {
		o.GlobalConfig.DbFileName = c.String("database")
	}

	if c.IsSet("logfile") || o.GlobalConfig.LogFileName == "" {
		o.GlobalConfig.LogFileName = c.String("logfile")
	}

	if c.IsSet("date-format") || o.GlobalConfig.DateFormat == "" {
		o.GlobalConfig.DateFormat = c.String("date-format")
	}
	if c.IsSet("today") {
		now, err := time.Parse(o.GlobalConfig.DateFormat, c.String("today"))
		if err != nil {
			return err
		}
		o.GlobalConfig.Now = now
	}
	return nil
}

func (o *Options) populateResolver(c *cli.Context) {
	if c.IsSet("maxdepth") || o.ResolverConfig.MaxDepth == 0 {
		o.ResolverConfig.MaxDepth = c.Int("maxdepth")
	}
}

func GetTimeFromString(now time.Time, format string, date string) (time.Time, error) {
	if date == "today" {
		return now.Local(), nil
	}
	if date == "yesterday" {
		return now.AddDate(0, 0, -1), nil
	}
	if date == "last7" {
		return now.AddDate(0, 0, -7), nil
	}
	if date == "last30" {
		return now.AddDate(0, 0, -30), nil
	}
	var err error
	var customTime time.Time
	customTime, err = time.Parse(format, date)
	if err == nil {
		return customTime, nil
	}
	return naturaldate.Parse(date, time.Now())
}

func (o *Options) populateFilter(c *cli.Context) error {
	for i := len(c.Lineage()) - 1; i >= 0; i-- {
		if c.Lineage()[i].IsSet("begin") {
			time, err := GetTimeFromString(o.GlobalConfig.Now, o.GlobalConfig.DateFormat, c.Lineage()[i].String("begin"))
			if err != nil {
				return err
			}
			o.FilterConfig.BeginningTime = &time
		}
		if c.Lineage()[i].IsSet("end") {
			time, err := GetTimeFromString(o.GlobalConfig.Now, o.GlobalConfig.DateFormat, c.Lineage()[i].String("end"))
			if err != nil {
				return err
			}
			o.FilterConfig.EndTime = &time
		}
	}
	return nil
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
	}
	o.ReporterConfig.SingleFood = c.String("single-food")
	o.ReporterConfig.ElementGroupByFood = c.Bool("group-food")
	o.ReporterConfig.SingleElement = c.String("single-element")
}
