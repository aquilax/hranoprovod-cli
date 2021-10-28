# Installation

## From source

First make sure you have go (golang) installed.

    http://golang.org/

Clone the repository and build/install the tool:

    git clone https://github.com/aquilax/hranoprovod-cli.git
    cd hranoprovod-cli
    go install

## Github Releases

You can download `.deb` and `.rpm` packages from [GitHub releases](https://github.com/aquilax/hranoprovod-cli/releases)

### asdf

[asdf](https://github.com/asdf-vm/asdf) is an extendable version manager for linux and macOS.

hranoprovod can be installed using a plugin as follows:

    asdf plugin add hranoprovod https://github.com/aquilax/asdf-hranoprovod.git
    asdf install hranoprovod latest

## Snapcraft

You can find the latest Snapcraft releases [here](https://snapcraft.io/hranoprovod-cli)

## Docker

You can run hranoprovod-cli from a docker container as well.

### Build the container locally

    git clone https://github.com/aquilax/hranoprovod-cli.git
    cd hranoprovod-cli
    docker build --pull --rm -f "Dockerfile" -t aquilax/hranoprovod-cli:latest .

### Download the latest container from Docker Hub

You can find the latest image from the docker hub page: https://hub.docker.com/r/aquilax/hranoprovod-cli

### Run the container with data files

    docker run --rm -it -v /path/to/data/files/:/data aquilax/hranoprovod-cli:latest -d /data/food.yaml -l /data/log.yaml bal
