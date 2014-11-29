package main

import (
	"code.google.com/p/gcfg"
	"github.com/Hranoprovod/processor"
	"github.com/Hranoprovod/reporter"
	"github.com/codegangsta/cli"
	"os"
	"os/user"
	"time"
)

const (
	configFileName = "/.hranoprovod/config"
)

// Config contains the configuration structure
type Config struct {
	Global struct {
		DbFileName  string
		LogFileName string
		DateFormat  string
	}
	Resolver struct {
		ResolverMaxDepth int
	}
	Processor processor.Options
	Reporter  reporter.Options
	API       struct {
	}
}

// NewConfig returns new config structure.
func NewConfig() *Config {
	config := &Config{}
	config.Reporter = *reporter.NewDefaultOptions()
	config.Processor = *processor.NewDefaultOptions()
	return config
}

// Load loads the settigns from config file/command line params/defauls from given context.
func (config *Config) Load(c *cli.Context) *Config {
	fileName := c.GlobalString("config")
	// First try to load the config file
	if exists(fileName) {
		if err := gcfg.ReadFileInto(config, fileName); err != nil {
			// Config file is not valid
			panic(err)
		}
	}
	config.populateGlobals(c)
	config.populateLocals(c)
	return config
}

// GetDefaultFileName returns the default filename for the config file
func GetDefaultFileName() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.HomeDir + configFileName
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func (config *Config) populateGlobals(c *cli.Context) {
	if c.GlobalIsSet("database") || config.Global.DbFileName == "" {
		config.Global.DbFileName = c.GlobalString("database")
	}

	if c.GlobalIsSet("logfile") || config.Global.LogFileName == "" {
		config.Global.LogFileName = c.GlobalString("logfile")
	}

	if c.GlobalIsSet("date-format") || config.Global.DateFormat == "" {
		config.Global.DateFormat = c.GlobalString("date-format")
	}
}

func (config *Config) populateLocals(c *cli.Context) {
	config.populateResolver(c)
	config.populateProcessor(c)
	config.populateReporter(c)
}

func (config *Config) populateResolver(c *cli.Context) {
	if c.IsSet("maxdepth") || config.Resolver.ResolverMaxDepth == 0 {
		config.Resolver.ResolverMaxDepth = c.Int("maxdepth")
	}
}

func (config *Config) populateProcessor(c *cli.Context) {
	var err error

	if c.IsSet("beginning") {
		config.Processor.BeginningTime, err = time.Parse(config.Global.DateFormat, c.String("beginning"))
		if err != nil {
			panic(err)
		}
		config.Processor.HasBeginning = true
	}

	if c.IsSet("end") {
		config.Processor.EndTime, err = time.Parse(config.Global.DateFormat, c.String("end"))
		if err != nil {
			panic(err)
		}
		config.Processor.HasEnd = true
	}

	config.Processor.Unresolved = c.Bool("unresolved")
	config.Processor.SingleFood = c.String("single-food")
	config.Processor.SingleElement = c.String("single-element")
}

func (config *Config) populateReporter(c *cli.Context) {
	if c.IsSet("csv") {
		config.Reporter.CSV = true
	}
	if c.IsSet("no-color") {
		config.Reporter.Color = false
	}

	if c.IsSet("no-totals") {
		config.Processor.Totals = false
	}
}
