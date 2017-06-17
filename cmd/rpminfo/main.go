package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/cavaliercoder/go-rpm"
)

// Signature see:
// https://github.com/rpm-software-management/rpm/blob/b74096a751293ed770b3cfb0a8793117fc45e7f2/lib/formats.c#L371

const defaultQueryFormat = `Name        : {{ .Name }}
Version     : {{ .Version }}
Release     : {{ .Release }}
Architecture: {{ .Architecture }}
Group       : {{ .Groups | join }}
Size        : {{ .Size }}
License     : {{ .License }}
Signature   : {{ .GPGSignature }}
Source RPM  : {{ .SourceRPM }}
Build Date  : {{ .BuildTime | timestamp }}
Build Host  : {{ .BuildHost }}
Packager    : {{ .Packager }}
Vendor      : {{ .Vendor }}
URL         : {{ .URL }}
Summary     : {{ .Summary }}
Description :
{{ .Description }}
`

func main() {
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		os.Exit(usage(1))
	}

	qf, err := queryformat(defaultQueryFormat)
	dieOn(err)

	for i, path := range os.Args[1:] {
		if i > 0 {
			fmt.Printf("\n")
		}

		p, err := rpm.OpenPackageFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
			continue
		}

		if err := qf.Execute(os.Stdout, p); err != nil {
			fmt.Fprintf(os.Stderr, "error formatting %s: %v\n", path, err)
			continue
		}
	}
}

func queryformat(tmpl string) (*template.Template, error) {
	return template.New("queryformat").
		Funcs(template.FuncMap{
			"join": func(a []string) string {
				return strings.Join(a, ", ")
			},
			"timestamp": func(t time.Time) rpm.Time {
				return rpm.Time(t)
			},
		}).
		Parse(tmpl)
}

func usage(exitCode int) int {
	w := os.Stdout
	if exitCode != 0 {
		w = os.Stderr
	}

	fmt.Fprintf(w, "usage: %v [path ...]\n", os.Args[0])
	return exitCode
}

func die(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func dieOn(err error) {
	if err != nil {
		die("%v\n", err)
	}
}
