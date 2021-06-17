# hranoprovod-cli [![Build Status](https://travis-ci.org/aquilax/hranoprovod-cli.svg?branch=master)](https://travis-ci.org/aquilax/hranoprovod-cli) [![GoDoc](https://godoc.org/github.com/aquilax/hranoprovod-cli?status.svg)](https://godoc.org/github.com/aquilax/hranoprovod-cli) [![Go Report Card](https://goreportcard.com/badge/github.com/aquilax/hranoprovod-cli)](https://goreportcard.com/report/github.com/aquilax/hranoprovod-cli) [![Documentation Status](https://readthedocs.org/projects/hranoprovod/badge/?version=latest)](https://hranoprovod.readthedocs.io/en/latest/?badge=latest) [![hranoprovod-cli](https://snapcraft.io/hranoprovod-cli/badge.svg)](https://snapcraft.io/hranoprovod-cli)

## Description

Hranoprovod is command line tracking tool. It supports nested recipes and custom defined tracking elements, which makes it perfect for tracking calories, nutrition data, exercises and other accumulative data.

[![asciicast](https://asciinema.org/a/257200.svg)](https://asciinema.org/a/257200)

## Installation

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/hranoprovod-cli)

First make sure you have go (golang) installed.

    http://golang.org/

Clone the repository and build/install the tool:

    git clone https://github.com/aquilax/hranoprovod-cli.git
    cd hranoprovod-cli
    go install

## Docker

You can run hranoprovod-cli from Docker too

### Building the image

```sh
docker build --pull --rm -f "Dockerfile" -t aquilax/hranoprovod-cli:latest .
```

### Running a balance report

```sh
docker run --rm -it -v /path/to/data/files/:/data aquilax/hranoprovod-cli:latest -d /data/food.yaml -l /data/log.yaml bal
```

## Help

Running the `hranoprovod-cli` command will show you the command line options

```sh
$ ./hranoprovod-cli --help
NAME:
   hranoprovod-cli - Lifestyle tracker

USAGE:
   hranoprovod-cli [global options] command [command options] [arguments...]

VERSION:
   dev, commit none, built at unknown

COMMANDS:
   register, reg  Shows the log register report
   balance, bal   Shows food balance as tree
   lint           Lints file
   report         Generates various reports
   csv            Generates csv exports
   stats          Provide stats information
   summary        Show summary
   gen            Generate documentation
   help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --begin DATE, -b DATE      Beginning of period DATE
   --end DATE, -e DATE        End of period DATE
   --database FILE, -d FILE   optional database file name FILE (default: "food.yaml") [$HR_DATABASE]
   --logfile FILE, -l FILE    log file name FILE (default: "log.yaml") [$HR_LOGFILE]
   --config FILE, -c FILE     Configuration file FILE (default: "/home/aquilax/.hranoprovod/config") [$HR_CONFIG]
   --date-format DATE_FORMAT  Date format for parsing and printing dates DATE_FORMAT (default: "2006/01/02") [$HR_DATE_FORMAT]
   --maxdepth DEPTH           Resolve depth DEPTH (default: 10) [$HR_MAXDEPTH]
   --no-color                 Disable color output (default: false)
   --help, -h                 show help (default: false)
   --version, -v              print the version (default: false)

```

## Usage

Hranoprovod uses two files with similar format to operate.

### Database file (food.yaml)

Contains all the "recipes" in the following format:

```sh
$ cat examples/food.yaml
# daily nutrition budget
day/nonworking:
  calories: -1200
  fat: -124
  carbohydrate: -50
  protein: -104

bread/rye/100g:
  calories: 259
  fat: 3.3
  carbohydrate: 48
  protein: 9

egg/boiled/100g:
  calories: 155
  fat: 11
  carbohydrate: 1.1
  protein: 13

vegetables/lettuce/romaine/100g:
  calories: 15
  fat: 0.5
  carbohydrate: 1.7
  protein: 0.9

sauce/mayonnaise/100g:
  calories: 680
  fat: 7.5
  carbohydrate: 0.6
  protein: 1

sandwich/egg/lettuce/100g:
  bread/rye/100g: 0.40
  egg/boiled/100g: 0.20
  vegetables/lettuce/romaine/100g: 0.20
  sauce/mayonnaise/100g: 0.20

candy/snickers/bar:
  calories: 280
  fat: 13.6
  carbohydrate: 35.1
  protein: 4.29
```

Hranoprovod is measure agnostic and it's up to the user to use or state the measurements.

### Log file (log.yaml)

The log file contains dated usage of the recipes, defined in the database file.

```sh
$ cat examples/log.yaml
2021/01/24:
  day/nonworking: 1
  coffee/cup: 1
  sandwich/egg/lettuce/100g: 1.20
  candy/snickers/bar: 1

2021/01/25:
  day/nonworking: 1
  coffee/cup: 1
  sandwich/egg/lettuce/100g: 1.50
  coffee/cup: 1

```

Note: it's not mandatory to have the elements in the database file. Elements which are not found will be represented as they are. They can always be added later to the database.

#### Register

Given this example, the result will look like:

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color reg
2021/01/24
	day/nonworking              :      1.00
		            calories   -1200.00
		        carbohydrate     -50.00
		                 fat    -124.00
		             protein    -104.00
	coffee/cup                  :      1.00
		          coffee/cup       1.00
	sandwich/egg/lettuce/100g   :      1.20
		            calories     328.32
		        carbohydrate      23.86
		                 fat       6.14
		             protein       7.90
	candy/snickers/bar          :      1.00
		            calories     280.00
		        carbohydrate      35.10
		                 fat      13.60
		             protein       4.29
	-- TOTAL  ----------------------------------------------------
		            calories     608.32   -1200.00 =   -591.68
		        carbohydrate      58.96     -50.00 =      8.96
		          coffee/cup       1.00       0.00 =      1.00
		                 fat      19.74    -124.00 =   -104.26
		             protein      12.19    -104.00 =    -91.81
2021/01/25
	day/nonworking              :      1.00
		            calories   -1200.00
		        carbohydrate     -50.00
		                 fat    -124.00
		             protein    -104.00
	coffee/cup                  :      2.00
		          coffee/cup       2.00
	sandwich/egg/lettuce/100g   :      1.50
		            calories     410.40
		        carbohydrate      29.82
		                 fat       7.68
		             protein       9.87
	-- TOTAL  ----------------------------------------------------
		            calories     410.40   -1200.00 =   -789.60
		        carbohydrate      29.82     -50.00 =    -20.18
		          coffee/cup       2.00       0.00 =      2.00
		                 fat       7.68    -124.00 =   -116.32
		             protein       9.87    -104.00 =    -94.13

```

#### Balance tree

You can also generate balance tree for single nutrition value:

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal -s calories
    280.00 | candy
    280.00 |   snickers
    280.00 |     bar
  -2400.00 | day
  -2400.00 |   nonworking
    738.72 | sandwich
    738.72 |   egg
    738.72 |     lettuce
    738.72 |       100g
-----------|
  -1381.28 | calories

```

Same result in slightly more compact format:

```sh
$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal -s calories -c
    280.00 | candy/snickers/bar
  -2400.00 | day/nonworking
    738.72 | sandwich/egg/lettuce/100g
-----------|
  -1381.28 | calories

```
