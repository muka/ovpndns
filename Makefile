
# The name of the executable (default is current directory name)
TARGET := `basename ${PWD}`
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 1.0.0
BUILD := `git rev-parse HEAD`

ARCH ?= amd64
GOARCH ?= ${ARCH}
GOARM ?= 7

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-s -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall run docker/build docker/push

all: build

$(TARGET): $(SRC)
	CGO_ENABLED=0 ARCH=${ARCH} GOARCH=${GOARCH} GOARM=${GOARM} go build $(LDFLAGS) -o ./build/${TARGET}-$(ARCH)

build: $(TARGET)
	@true

clean:
	@rm -f build/$(TARGET)

install:
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})

run: install
	@$(TARGET)

docker/build:
	@docker build . -t raptorbox/$(TARGET)-$(ARCH)

docker/push: docker/build
	@docker push raptorbox/$(TARGET)-$(ARCH)
