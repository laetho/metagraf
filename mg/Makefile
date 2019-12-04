# Usage:
# make
# make clean
# make build

.PHONY = build

GIT_COMMIT = $(git rev-list -1 HEAD) 

build:
	@echo "Building Linux version..."
	echo ${GIT_COMMIT}
	go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" 

build-osx:

build-win:

clean:
	@echo "Cleaning up binaries..."
	rm mg
	rm mg.exe

