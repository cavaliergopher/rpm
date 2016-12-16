package main

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(usage(1))
	}

	fmt.Printf("---\n")
	for i, path := range os.Args[1:] {
		if i > 0 {
			fmt.Printf("\n")
		}

		raw(path)
	}
}

func raw(path string) {
	fmt.Printf("- path: %v\n", path)
	p, err := rpm.OpenPackageFile(path)
	if err != nil {
		die("%v: %v\n", path, err)
	}

	fmt.Printf("  headers:\n")
	for i, h := range p.Headers {
		if i > 0 {
			fmt.Printf("\n")
		}

		fmt.Printf("  - index: %v\n", i)
		fmt.Printf("    version: %v\n", h.Version)
		fmt.Printf("    start: %v\n", h.Start)
		fmt.Printf("    end: %v\n", h.End)
		fmt.Printf("    length: %v\n", h.Length)
		fmt.Printf("    indexes:\n")
		for j, ix := range h.Indexes {
			if j > 0 {
				fmt.Printf("\n")
			}

			fmt.Printf("    - index: %v\n", j)
			fmt.Printf("      tag: %v\n", ix.Tag)
			fmt.Printf("      type: %v\n", ix.Type)
			fmt.Printf("      offset: %v\n", ix.Offset)

			switch ix.Value.(type) {
			case []string:
				ss := ix.Value.([]string)
				if len(ss) == 1 && strings.Index(ss[0], "\n") == -1 {
					fmt.Printf("      value: [ \"%v\" ]\n", ss[0])
				} else {
					fmt.Printf("      value:\n")
					for _, s := range ss {
						if strings.Index(s, "\n") == -1 {
							fmt.Printf("      - \"%v\"\n", s)
						} else {
							fmt.Printf("      - |\n")
							lines := strings.Split(s, "\n")
							for _, line := range lines {
								fmt.Printf("        %v\n", line)
							}
						}
					}
				}

			default:
				fmt.Printf("      value: %v\n", ix.Value)
			}
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

func die(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func dieOn(err error) {
	if err != nil {
		die("%v\n", err)
	}
}
