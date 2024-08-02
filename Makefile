

.PHONY: build
build:
	go build

.PHONY: upgrade
upgrade:
	go get -u -v ./...

.PHONY: up-build
up-build: upgrade build
