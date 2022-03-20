.PHONY: clean test test-release snapshots

unexport HR_DATABASE
unexport HR_LOGFILE
unexport HR_CONFIG
unexport HR_DATE_FORMAT

SHELL=/bin/bash
BINARY=hranoprovod-cli
MDEXEC=mdexec
TARGETS=${BINARY} docs/command-line.md docs/usage.md README.md docs/usage.cast

all: $(WASM) $(BINARY) docs

$(BINARY):
	go build -o $(BINARY)

docs: $(BINARY) docs/command-line.md docs/usage.md README.md

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
	rm -f $(TARGETS)

documentation/usage.cast: $(BINARY) scripts/usage.sh
	asciinema rec --overwrite -c "scripts/usage.sh -n" documentation/usage.cast

cast: $(BINARY) docs/usage.cast

snapshots:
	UPDATE_SNAPSHOTS=1 go test ./...

coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
	# rm coverage.out