package rpm_test

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm"
	"os"
)

// ExampleMD5Check validates a local rpm package named using the MD5 checksum
// value specified in the package header.
func ExampleMD5Check() {
	// open a rpm package for reading
	f, err := os.Open("testdata/centos-release-7-2.1511.el7.centos.2.10.x86_64.rpm")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	// validate md5 checksum
	if err := rpm.MD5Check(f); err == nil {
		fmt.Printf("Package passed checksum validation\n")
	} else if err == rpm.ErrMD5ValidationFailed {
		fmt.Printf("Package failed checksum validation\n")
	} else {
		panic(err)
	}

	// Output: Package passed checksum validation
}
