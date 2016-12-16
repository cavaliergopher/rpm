rpmdump_SOURCES = \
	main.go

all: rpmdump

rpmdump: $(rpmdump_SOURCES)
	go build -x -o rpmdump $(rpmdump_SOURCES)

install:
	go install -x github.com/cavaliercoder/go-rpm/cmd/rpmdump

clean:
	go clean -x -i

.PHONY: all install clean
