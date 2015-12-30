all: build

build:
	go build -x

test:
	go test -v

get-deps:
	go get github.com/mattn/go-sqlite3
