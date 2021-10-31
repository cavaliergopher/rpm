/*
rpminfo displays package information, akin to rpm --info.

	usage: rpminfo [package ...]

Example:

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
*/
package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/cavaliergopher/rpm"
)

var tmpl = template.Must(template.New("rpminfo").
	Funcs(template.FuncMap{
		"join": func(a []string) string {
			return strings.Join(a, ", ")
		},
		"strftime": func(t time.Time) string {
			return t.Format(rpm.TimeFormat)
		},
	}).
	Parse(`Name        : {{ .Name }}
Version     : {{ .Version }}
Release     : {{ .Release }}
Architecture: {{ .Architecture }}
Group       : {{ .Groups | join }}
Size        : {{ .Size }}
License     : {{ .License }}
Signature   : {{ .GPGSignature }}
Source RPM  : {{ .SourceRPM }}
Build Date  : {{ strftime .BuildTime }}
Build Host  : {{ .BuildHost }}
Packager    : {{ .Packager }}
Vendor      : {{ .Vendor }}
URL         : {{ .URL }}
Summary     : {{ .Summary }}
Description :
{{ .Description }}
`))

func main() {
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		os.Exit(usage(1))
	}
	for i, name := range os.Args[1:] {
		if i > 0 {
			fmt.Printf("\n")
		}
		p, err := rpm.Open(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", name, err)
			continue
		}
		if err := tmpl.Execute(os.Stdout, p); err != nil {
			fmt.Fprintf(os.Stderr, "error formatting %s: %v\n", name, err)
			continue
		}
	}
}

func usage(exitCode int) int {
	w := os.Stdout
	if exitCode != 0 {
		w = os.Stderr
	}
	fmt.Fprintf(w, "usage: %v [path ...]\n", os.Args[0])
	return exitCode
}
