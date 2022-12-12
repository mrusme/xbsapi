.PHONY: swagger ent build install-deps install-dep-swag
VERSION := $(shell git describe --tags 2> /dev/null || git rev-parse --short HEAD)

all: ent build swagger

swagger:
	swag init -g api/v1/v1.go

ent:
	go generate ./ent

build:
	go build -ldflags "-X github.com/mrusme/xbsapi/lib.VERSION=$(VERSION)"

install-deps: install-deps-go install-dep-ent install-dep-swag

install-deps-go:
	go get

install-dep-ent:
	go install entgo.io/ent/cmd/ent@latest

install-dep-swag:
	go install github.com/swaggo/swag/cmd/swag@latest


