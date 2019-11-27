## hranoprovod-cli [![Build Status](https://travis-ci.org/aquilax/hranoprovod-cli.svg?branch=master)](https://travis-ci.org/aquilax/hranoprovod-cli) [![GoDoc](https://godoc.org/github.com/aquilax/hranoprovod-cli?status.svg)](https://godoc.org/github.com/aquilax/hranoprovod-cli) [![Go Report Card](https://goreportcard.com/badge/github.com/aquilax/hranoprovod-cli)](https://goreportcard.com/report/github.com/aquilax/hranoprovod-cli)

## Description

Hranoprovod is command line tracking tool. It supports nested recipies and custom defined tracking elements, which makes it perfect for tracking calories, nutrition data, exercises and other accumulative data.

[![asciicast](https://asciinema.org/a/257200.svg)](https://asciinema.org/a/257200)

## Installation

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/hranoprovod-cli)

First make sure you have go (golang) installed.

    http://golang.org/

Clone the repository and build/install the tool:

    git clone https://github.com/aquilax/hranoprovod-cli.git
    cd hranoprovod-cli
    go install

## Help

Running the `hranoprovod-cli` command will show you the command line options

```
$ hranoprovod-cli
NAME:
   hranoprovod-cli - Lifestyle tracker

USAGE:
   hranoprovod-cli [global options] command [command options] [arguments...]

VERSION:
   2.2.0

AUTHOR:
   aquilax <aquilax@gmail.com>

COMMANDS:
     register, reg  Shows the log register report
     balance, bal   Shows food balance as tree
     add            Adds new item to the log
     api            Service API commands
     lint           Lints file
     report         Generates various reports
     csv            Generates csv exports
     stats          Provide summary information
     help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --begin value, -b value     Beginning of period
   --end value, -e value       End of period
   --database value, -d value  database file name (default: "food.yaml") [$HR_DATABASE]
   --logfile value, -l value   log file name (default: "log.yaml") [$HR_LOGFILE]
   --config value, -c value    Configuration file (default: "/home/aquilax/.hranoprovod/config") [$HR_CONFIG]
   --date-format value         Date format for parsing and printing dates (default: "2006/01/02") [$HR_DATE_FORMAT]
   --maxdepth value            Resolve depth (default: 10) [$HR_MAXDEPTH]
   --help, -h                  show help
   --version, -v               print the version
```

## Usage

Hranoprovod uses two files with similar format to operate.

### Database file (food.yaml)

Contains all the "recipes" in the following format:

```
fish/tuna/canned/100g:
  calories: 184
  fat: 6
  carbohydrate: 0
  protein: 0

bread/white/100g:
  calories: 265
  fat: 3.2
  carbohydrate: 49
  protein: 9
```

Let's say you love tuna sandwiches then you can combine these two ingredients into one:

```
sandwich/tuna/100g:
  fish/tuna/canned/100g: .6
  bread/white/100g: .4

sandwich/tuna/pc:
  sandwich/tuna/100g: 1.5
```

This means that the sandwich is composed of 60% tuna and 40% bread and a sandwich weights arount 150g.

Hranoprovod is measure agnostic and it's up to the user to use or state the measurements.

### Log file (log.yaml)

The log file contains dated usage of the recipes, defined in the database file.

```
2014/12/17:
  tea/cup: 1
  sandwich/tuna/pc: 2
  calories: 300
  biking/km: 10
```

Note: it's not mandatory to have the elements in the database file. Elements which are not found will be represented as they are. They can always be added later to the database.

#### Register

Given this example, the result will look like:

```
$hranoprovod-cli -d food.yaml -l log.yaml reg
2014/12/17
	tea/cup                     :      1.00
		             tea/cup       1.00
	sandwich/tuna/pc            :      2.00
		            calories     649.20
		        carbohydrate      58.80
		                 fat      14.64
		             protein      10.80
	calories                    :    300.00
		            calories     300.00
	biking/km                   :     10.00
		           biking/km      10.00
	-- TOTAL  ----------------------------------------------------
		           biking/km      10.00       0.00 =     10.00
		            calories     949.20       0.00 =    949.20
		        carbohydrate      58.80       0.00 =     58.80
		                 fat      14.64       0.00 =     14.64
		             protein      10.80       0.00 =     10.80
		             tea/cup       1.00       0.00 =      1.00
```

#### Balance tree

You can also generate balance tree for single nutrition value:

```
$ hranoprovod-cli bal -b yesterday -s calories
    329.82 | butter
    329.82 |   cow milk
    329.82 |     100g
     44.20 | cream
     44.20 |   heavy
     44.20 |     36%
     44.20 |       100g
  -1632.00 | day
  -1632.00 |   nonworking
      2.40 | drinks
      2.40 |   coffee
      2.40 |     cup
    305.61 | eggs
    305.61 |   fried
    305.61 |     pc
      8.94 | garlic
      8.94 |   100g
    100.80 | olives
    100.80 |   brown
    100.80 |     100g
      7.20 | rucola
      7.20 |   100g
     54.90 | vegetables
     54.90 |   spinach
     54.90 |     frozen
     54.90 |       100g
    148.40 | vegokorv
    148.40 |   pc
-----------|
   -629.73 | calories
```

Same result in slightly more compact format:
```
$ hranoprovod-cli bal -b yesterday -s calories -c
    329.82 | butter
    329.82 |   cow milk/100g
     44.20 | cream
     44.20 |   heavy
     44.20 |     36%/100g
  -1632.00 | day/nonworking
      2.40 | drinks
      2.40 |   coffee/cup
    305.61 | eggs
    305.61 |   fried/pc
      8.94 | garlic/100g
    100.80 | olives
    100.80 |   brown/100g
      7.20 | rucola/100g
     54.90 | vegetables
     54.90 |   spinach
     54.90 |     frozen/100g
    148.40 | vegokorv/pc
-----------|
   -629.73 | calories
```
