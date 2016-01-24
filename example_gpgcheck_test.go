package rpm_test

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm"
	"os"
)

// ExampleGPGCheck reads a public GPG key and uses it to validate the signature
// of a local rpm package.
func ExampleGPGCheck() {
	// read public key from gpgkey file
	keyring, err := rpm.KeyRingFromFile("testdata/RPM-GPG-KEY-CentOS-7")
	if err != nil {
		panic(err)
	}

	// open a rpm package for reading
	f, err := os.Open("testdata/centos-release-7-2.1511.el7.centos.2.10.x86_64.rpm")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	// validate gpg signature
	if signer, err := rpm.GPGCheck(f, keyring); err == nil {
		fmt.Printf("Package signed by '%s'\n", signer)
	} else if err == rpm.ErrGPGValidationFailed {
		fmt.Printf("Package failed GPG signature validation\n")
	} else {
		panic(err)
	}

	// Output: Package signed by 'CentOS-7 Key (CentOS 7 Official Signing Key) <security@centos.org>'
}
