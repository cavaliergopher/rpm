# go-rpm [![GoDoc](https://godoc.org/github.com/cavaliercoder/go-rpm?status.svg)](https://godoc.org/github.com/cavaliercoder/go-rpm)

A native implementation of the RPM file specification in Go.

	$ go get github.com/cavaliercoder/go-rpm


The go-rpm package aims to enable cross-platform tooling for yum/dnf/rpm
written in Go (E.g. [y10k](https://github.com/cavaliercoder/y10k)).

Initial goals include like-for-like implementation of existing rpm ecosystem
features such as:

* Reading of modern and legacy rpm package file formats
* Reading, creating and updating modern and legacy yum repository metadata


## License

Copyright (c) 2015 Ryan Armstrong 

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.