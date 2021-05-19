all: usage 

usage:
	@echo "Default is showing usage"

build-linux:
	@GITHASH = $(shell git rev-parse --short HEAD)
	$(GITHASH)

