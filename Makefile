.PHONY: clean test test-release

unexport HR_DATABASE
unexport HR_LOGFILE
unexport HR_CONFIG
unexport HR_DATE_FORMAT

SHELL=/bin/bash
BINARY=hranoprovod-cli
MDEXEC=mdexec
TARGETS=${BINARY} docs/command-line.md docs/usage.md README.md

all: $(WASM) $(BINARY) documentation

$(BINARY):
	go build -o $(BINARY)

documentation: $(BINARY) docs/command-line.md docs/usage.md README.md

docs/command-line.md:
	./$(BINARY) gen markdown > $@

docs/usage.md:
	$(MDEXEC) documentation/usage.md > $@

README.md:
	$(MDEXEC) documentation/README.md > $@

test:
	go test -v ./...

test-release:
	goreleaser --snapshot --skip-publish --rm-dist

clean:
	rm $(TARGETS)