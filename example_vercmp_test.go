package rpm_test

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cavaliercoder/go-rpm"
	"github.com/cavaliercoder/go-rpm/version"
)

// ExampleVersionCompare reads packages in a directory and prints the name of
// the package with the highest version.
func ExampleVersionCompare() {
	// list files in directory
	dir, err := ioutil.ReadDir("testdata")
	if err != nil {
		panic(err)
	}

	// test each package file
	var latest *rpm.PackageFile = nil
	for _, f := range dir {
		if strings.HasSuffix(f.Name(), ".rpm") {
			// read package file
			pkg, err := rpm.OpenPackageFile("testdata/" + f.Name())
			if err != nil {
				panic(err)
			}

			// compare versions
			if 1 == version.Compare(pkg, latest) {
				latest = pkg
			}
		}
	}

	fmt.Printf("Latest package: %v\n", latest)

	// Output: Latest package: centos-release-5-0.0.el5.centos.2.i386
}
