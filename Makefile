.PHONY: test run build

test:
	go test -v ./lib

build: test
	go build .

run:
	time go run .
