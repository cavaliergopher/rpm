/*
Package rpm provides an implementation the rpm package file format.

Methods for comparing package versions are provided in pkg/rpmver.

See: http://ftp.rpm.org/max-rpm/s1-rpm-file-format-rpm-file-format.html

	package main

	import (
		"fmt"
		"log"

		"github.com/cavaliergopher/rpm"
	)

	func main() {
		p, err := rpm.Open("my-package.rpm")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Loaded package: %v", p)
	}

*/
package rpm
