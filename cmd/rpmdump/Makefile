rpmdump_SOURCES = \
	main.go

all: rpmdump

rpmdump: $(rpmdump_SOURCES)
	go build -x -o rpmdump $(rpmdump_SOURCES)

clean:
	go clean -x

.PHONY: all clean
