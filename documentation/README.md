# hranoprovod-cli [![GoDoc](https://godoc.org/github.com/aquilax/hranoprovod-cli?status.svg)](https://godoc.org/github.com/aquilax/hranoprovod-cli) [![Go Report Card](https://goreportcard.com/badge/github.com/aquilax/hranoprovod-cli)](https://goreportcard.com/report/github.com/aquilax/hranoprovod-cli) [![Documentation Status](https://readthedocs.org/projects/hranoprovod/badge/?version=latest)](https://hranoprovod.readthedocs.io/en/latest/?badge=latest) [![hranoprovod-cli](https://snapcraft.io/hranoprovod-cli/badge.svg)](https://snapcraft.io/hranoprovod-cli)

## Description

Hranoprovod is command line diet tracking tool. It supports nested recipes and is agnostic to the metrics you choose to track, which makes it perfect for tracking calories, nutrition data, exercises and any other measurement that is accumulated on daily basis.

[![asciicast](https://asciinema.org/a/257200.svg)](https://asciinema.org/a/257200)

## Installation

### asdf

[asdf](https://github.com/asdf-vm/asdf) is an extendable version manager for linux and macOS.

hranoprovod can be installed using a plugin as follows:

    asdf plugin add hranoprovod https://github.com/aquilax/asdf-hranoprovod.git
    asdf install hranoprovod latest

### from code

First make sure you have go (golang) installed.

    https://go.dev/

Clone the repository and build/install the tool:

    git clone https://github.com/aquilax/hranoprovod-cli.git
    cd hranoprovod-cli/cmd/hranoprovod-cli/
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

`$ ./hranoprovod-cli --help`

## Usage

Hranoprovod uses two files with very similar format to operate.

### Database file (food.yaml)

Contains all the "recipes" in the following format:

`$ cat examples/food.yaml`

Hranoprovod is measure agnostic and it's up to the user to use or state the measurements.

### Log journal file (log.yaml)

The log file contains dated usage of the recipes, defined in the database file.

`$ cat examples/log.yaml`

Note: it's not mandatory to have the elements in the database file. Elements which are not found will be represented as they are. They can always be added later to the database.

#### Register

Given this example, the result will look like:

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color reg`

#### Balance tree

You can also generate balance tree for single nutrition value:

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal -s calories`

Same result in slightly more compact format:

`$ ./hranoprovod-cli -d examples/food.yaml -l examples/log.yaml --no-color bal -s calories -c`
