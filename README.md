# rpm
[![Go Reference](https://pkg.go.dev/badge/github.com/cavaliergopher/rpm.svg)](https://pkg.go.dev/github.com/cavaliergopher/rpm) [![Build Status](https://app.travis-ci.com/cavaliergopher/rpm.svg?branch=main)](https://app.travis-ci.com/cavaliergopher/rpm) [![Go Report Card](https://goreportcard.com/badge/github.com/cavaliergopher/rpm)](https://goreportcard.com/report/github.com/cavaliergopher/rpm)

Package rpm implements the rpm package file format.

	$ go get github.com/cavaliergopher/rpm

This package also includes two example tools, `rpmdump` and `rpminfo`.

```
$ rpminfo golang-1.6.3-2.el7.x86_64.rpm
Name        : golang
Version     : 1.6.3
Release     : 2.el7
Architecture: x86_64
Group       : Unspecified
Size        : 11809071
License     : BSD and Public Domain
Signature   : RSA/SHA256, Sun Nov 20 18:01:16 2016, Key ID 24c6a8a7f4a80eb5
Source RPM  : golang-1.6.3-2.el7.src.rpm
Build Date  : Tue Nov 15 12:20:30 2016
Build Host  : c1bm.rdu2.centos.org
Packager    : CentOS BuildSystem <http://bugs.centos.org>
Vendor      : CentOS
URL         : http://golang.org/
Summary     : The Go Programming Language
Description :
The Go Programming Language.
```

## Extracting rpm packages

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

// The following working example demonstrates how to extract files from an rpm
// package. In this example, only the cpio format and xz compression are
// supported.
//
// Implementations should consider additional formats and compressions
// algorithms, as well as support for extracting irregular file types and
// configuring permissions, uids and guids, etc..
func ExtractRPM(name string) {
	// Open a package file for reading
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// read the package headers
	pkg, err := rpm.Read(f)
	if err != nil {
		log.Fatal(err)
	}

	// check the compression algorithm of the payload
	if compression := pkg.PayloadCompression(); compression != "xz" {
		log.Fatalf("Unsupported compression: %s", compression)
	}

	// decompress the payload
	xzReader, err := xz.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}

	// check the archive format of the payload
	if format := pkg.PayloadFormat(); format != "cpio" {
		log.Fatalf("Unsupported payload format: %s", format)
	}

	// Unarchive each file in the payload
	cpioReader := cpio.NewReader(xzReader)
	for {
		// move to the next file in the archive
		hdr, err := cpioReader.Next()
		if err == io.EOF {
			break // no more files
		}
		if err != nil {
			log.Fatal(err)
		}

		// skip directories and other irregular file types in this example
		if !hdr.Mode.IsRegular() {
			continue
		}

		// create the target directory
		if dirName := filepath.Dir(hdr.Name); dirName != "" {
			if err := os.MkdirAll(dirName, 0o755); err != nil {
				log.Fatal(err)
			}
		}

		// create and write the file
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
