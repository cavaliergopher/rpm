PACKAGE = github.com/cavaliercoder/go-rpm

rpm_la_SOURCES = \
	dependency.go \
	doc.go \
	dump.go \
	gpgcheck.go \
	header.go \
	index.go \
	keyring.go \
	lead.go \
	packagefile.go \
	tags.go \
	version.go

all: build

build: $(rpm_la_SOURCES)
	go build -x $(PACKAGE)

test:
	go test -v -cover

get-deps:
	go get github.com/cavaliercoder/badio
	go get github.com/dvyukov/go-fuzz/go-fuzz
	go get github.com/dvyukov/go-fuzz/go-fuzz-build
	go get golang.org/x/crypto/openpgp

rpm-fuzz.zip: *.go
	go-fuzz-build github.com/cavaliercoder/go-rpm

fuzz: rpm-fuzz.zip
	go-fuzz -bin=./rpm-fuzz.zip -workdir=.fuzz/

clean-fuzz:
	rm -rf rpm-fuzz.zip .fuzz/crashers/* .fuzz/suppressions/*

clean: clean-fuzz

.PHONY: all build test get-deps fuzz clean-fuzz clean
