PACKAGE = github.com/cavaliercoder/go-rpm

all: check install

check:
	go test -v -cover $(PACKAGE)/...

install:
	go install -x $(PACKAGE)/...

clean: clean-fuzz
	go clean -x -i $(PACKAGE)/...

rpm-fuzz.zip: *.go
	go-fuzz-build $(PACKAGE)

fuzz: rpm-fuzz.zip
	go-fuzz -bin=./rpm-fuzz.zip -workdir=.fuzz/

clean-fuzz:
	rm -rf rpm-fuzz.zip .fuzz/crashers/* .fuzz/suppressions/*

get-deps:
	go get github.com/cavaliercoder/badio
	go get github.com/dvyukov/go-fuzz/go-fuzz
	go get github.com/dvyukov/go-fuzz/go-fuzz-build
	go get golang.org/x/crypto/openpgp

.PHONY: all check install clean fuzz clean-fuzz get-deps
