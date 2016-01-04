package rpm_test

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm"
	"os"
)

func ExampleMD5Check() {
	// open a rpm package for reading
	f, err := os.Open("my-package.rpm")
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
}
