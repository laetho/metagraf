# Usage:
# make
# make clean
# make build

.PHONY = build

GIT_COMMIT = $(git rev-list -1 HEAD) 

build:
	@echo "Building Linux version..."
	go build
	sudo cp mg /usr/local/bin/mg
#	echo ${GIT_COMMIT}
#	go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" 

build-osx:
	@echo "Building Mac OSX version..."
	GOOS=darwin ARCH=amd64 go build

build-win:

clean:
	@echo "Cleaning up binaries..."
	rm mg
	rm mg.exe

