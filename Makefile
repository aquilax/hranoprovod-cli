.PHONY: clean test

unexport HR_DATABASE
unexport HR_LOGFILE
unexport HR_CONFIG
unexport HR_DATE_FORMAT

SHELL=/bin/bash
BINARY=hranoprovod-cli
MDEXEC=mdexec
TARGETS=${BINARY} docs/command-line.md docs/usage.md README.md

all: hranoprovod-cli documentation

$(BINARY):
	go build -o $(BINARY)

documentation: docs/command-line.md docs/usage.md README.md

docs/command-line.md:
	./$(BINARY) gen markdown > docs/command-line.md

docs/usage.md:
	$(MDEXEC) documentation/usage.md > docs/usage.md

README.md:
	$(MDEXEC) documentation/README.md > README.md

test:
	go test -v ./...

test-release:
	goreleaser --snapshot --skip-publish --rm-dist

clean:
	rm $(TARGETS)