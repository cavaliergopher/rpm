# rpm
[![Go Reference](https://pkg.go.dev/badge/github.com/cavaliergopher/rpm.svg)](https://pkg.go.dev/github.com/cavaliergopher/rpm) [![Build Status](https://app.travis-ci.com/cavaliergopher/rpm.svg?branch=main)](https://app.travis-ci.com/cavaliergopher/rpm) [![Go Report Card](https://goreportcard.com/badge/github.com/cavaliergopher/rpm)](https://goreportcard.com/report/github.com/cavaliergopher/rpm)

Package rpm providers readers and writers for RPM packages.

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
