.PHONY: all clean test test-release snapshots

unexport HR_DATABASE
unexport HR_LOGFILE
unexport HR_CONFIG
unexport HR_DATE_FORMAT

SHELL=/bin/bash
TARGET=hranoprovod-cli
MDEXEC=mdexec
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
DOC_TARGETS=${TARGET} docs/command-line.md docs/usage.md README.md docs/usage.cast

.DEFAULT_GOAL: $(TARGET)

all: $(TARGET) docs

$(TARGET): $(SRC)
	go build -o $(TARGET) hranoprovod/main.go

docs: $(TARGET) docs/command-line.md docs/usage.md README.md

docs/command-line.md: $(TARGET)
	./$(TARGET) gen markdown > $@

docs/usage.md: $(TARGET) documentation/usage.md
	$(MDEXEC) documentation/usage.md > $@

README.md: $(TARGET) documentation/README.md
	$(MDEXEC) documentation/README.md > $@

test:
	go test -v ./...

test-release:
	goreleaser --snapshot --skip-publish --rm-dist

clean:
	rm -f $(DOC_TARGETS)

documentation/usage.cast: $(TARGET) scripts/usage.sh
	asciinema rec --overwrite -c "scripts/usage.sh -n" documentation/usage.cast

cast: $(TARGET) documentation/usage.cast

snapshots:
	UPDATE_SNAPSHOTS=1 go test ./...

coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
	rm coverage.out
