all: build

build:
	go build -x

test:
	go test -v

get-deps:
	go get -u github.com/mattn/go-sqlite3

