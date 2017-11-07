.PHONY: install install-test

default: install
	go build

install:
	go get -v ./...

install-tests:
	go get -t -v ./...
	go get github.com/golang/lint/golint
	go get github.com/AlekSi/gocoverutil

lint: install-test
	golint ./...

test: install-test
	gocoverutil -coverprofile=cover.out test -v -covermode=count ./...
	go tool cover -html=cover.out -o coverage.html
