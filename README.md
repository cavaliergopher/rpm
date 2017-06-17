# go-rpm [![Build Status](https://travis-ci.org/cavaliercoder/go-rpm.svg?branch=master)](https://travis-ci.org/cavaliercoder/go-rpm) [![GoDoc](https://godoc.org/github.com/cavaliercoder/go-rpm?status.svg)](https://godoc.org/github.com/cavaliercoder/go-rpm)

A native implementation of the RPM file specification in Go.

	$ go get github.com/cavaliercoder/go-rpm


The go-rpm package aims to enable cross-platform tooling for yum/dnf/rpm
written in Go (E.g. [y10k](https://github.com/cavaliercoder/y10k)).

Initial goals include like-for-like implementation of existing rpm ecosystem
features such as:

* Reading of modern and legacy rpm package file formats
* Reading, creating and updating modern and legacy yum repository metadata
* Reading of the rpm database

```go
package main

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm"
)

func main() {
	p, err := rpm.OpenPackageFile("golang-1.6.3-2.el7.rpm")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded package: %v - %s\n", p, p.Summary())

	// Output: golang-0:1.6.3-2.el7.x86_64 - The Go Programming Language
}
```

## Tools

This package also includes two tools `rpmdump` and `rpminfo`.

The code for both tools demonstrates some use-cases of this package. They are
both also useful for interrogating RPM packages.

## License
Copyright (c) 2017 Ryan Armstrong. All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software without
   specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
