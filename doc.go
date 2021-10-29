/*
Package rpm implements the rpm package file format.

Methods for comparing package versions are provided in pkg/rpmver.

See: http://ftp.rpm.org/max-rpm/s1-rpm-file-format-rpm-file-format.html

	package main

	import (
		"fmt"
		"log"

		"github.com/cavaliergopher/rpm"
	)

	func main() {
		pkg, err := rpm.Open("golang-1.17.2-1.el7.x86_64.rpm")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Package:", pkg)
		fmt.Println("Summary:", pkg.Summary())

		// Output:
		// Package: golang-1.17.2-1.el7.x86_64
		// Summary: The Go Programming Language
	}

*/
package rpm
