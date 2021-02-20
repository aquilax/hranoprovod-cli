.PHONY: clean

SHELL=/bin/bash
BINARY=hranoprovod-cli
MDEXEC=mdexec
TARGETS=${BINARY} docs/command-line.md docs/usage.md

all: hranoprovod-cli documentation

$(BINARY):
	go build -o $(BINARY)

documentation: docs/command-line.md docs/usage.md

docs/command-line.md:
	./$(BINARY) gen markdown > docs/command-line.md

docs/usage.md:
	$(MDEXEC) documentation/usage.md > docs/usage.md

clean:
	rm $(TARGETS)