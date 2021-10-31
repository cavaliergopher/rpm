# rpm
[![Go Reference](https://pkg.go.dev/badge/github.com/cavaliergopher/rpm.svg)](https://pkg.go.dev/github.com/cavaliergopher/rpm) [![Build Status](https://app.travis-ci.com/cavaliergopher/rpm.svg?branch=main)](https://app.travis-ci.com/cavaliergopher/rpm) [![Go Report Card](https://goreportcard.com/badge/github.com/cavaliergopher/rpm)](https://goreportcard.com/report/github.com/cavaliergopher/rpm)

Package rpm implements the rpm package file format.

	$ go get github.com/cavaliergopher/rpm

See the [package documentation](https://pkg.go.dev/github.com/cavaliergopher/rpm)
or the examples programs in `cmd/` to get started.

## Extracting rpm packages

The following working example demonstrates how to extract files from an rpm
package. In this example, only the cpio format and xz compression are supported
which will cover most cases.

Implementations should consider additional formats and compressions algorithms,
as well as support for extracting irregular file types and configuring
permissions, uids and guids, etc.

```go
package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/cavaliergopher/cpio"
	"github.com/cavaliergopher/rpm"
	"github.com/ulikunitz/xz"
)

func ExtractRPM(name string) {
	// Open a package file for reading
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Read the package headers
	pkg, err := rpm.Read(f)
	if err != nil {
		log.Fatal(err)
	}

	// Check the compression algorithm of the payload
	if compression := pkg.PayloadCompression(); compression != "xz" {
		log.Fatalf("Unsupported compression: %s", compression)
	}

	// Attach a reader to decompress the payload
	xzReader, err := xz.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}

	// Check the archive format of the payload
	if format := pkg.PayloadFormat(); format != "cpio" {
		log.Fatalf("Unsupported payload format: %s", format)
	}

	// Attach a reader to unarchive each file in the payload
	cpioReader := cpio.NewReader(xzReader)
	for {
		// Move to the next file in the archive
		hdr, err := cpioReader.Next()
		if err == io.EOF {
			break // no more files
		}
		if err != nil {
			log.Fatal(err)
		}

		// Skip directories and other irregular file types in this example
		if !hdr.Mode.IsRegular() {
			continue
		}

		// Create the target directory
		if dirName := filepath.Dir(hdr.Name); dirName != "" {
			if err := os.MkdirAll(dirName, 0o755); err != nil {
				log.Fatal(err)
			}
		}

		// Create and write the file
		outFile, err := os.Create(hdr.Name)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(outFile, cpioReader); err != nil {
			outFile.Close()
			log.Fatal(err)
		}
		outFile.Close()
	}
}
```
