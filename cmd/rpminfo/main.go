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
