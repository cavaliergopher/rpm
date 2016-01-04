package rpm_test

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"os"
)

// ExampleGPGCheck loads a public key from /etc/pki/rpm-gpg/RPM-GPG-KEY-MY-KEY
// and uses it to validate the signature in a local rpm package named
// 'my-package.rpm'.
func ExampleGPGCheck() {
	// open gpgkey file
	f_gpgkey, err := os.Open("/etc/pki/rpm-gpg/RPM-GPG-KEY-MY-KEY")
	if err != nil {
		panic(err)
	}

	defer f_gpgkey.Close()

	// decode gpgkey file
	p, err := armor.Decode(f_gpgkey)
	if err != nil {
		panic(err)
	}

	// read public key from decoded gpgkey file
	keyring, err := openpgp.ReadKeyRing(p.Body)
	if err != nil {
		panic(err)
	}

	// open a rpm package for reading
	f_rpm, err := os.Open("my-package.rpm")
	if err != nil {
		panic(err)
	}

	defer f_rpm.Close()

	// validate gpg signature
	if signer, err := rpm.GPGCheck(f_rpm, keyring); err == nil {
		fmt.Printf("Package signed by '%s'\n", signer)
	} else if err == rpm.ErrGPGValidationFailed {
		fmt.Printf("Package failed GPG signature validation\n")
	} else {
		panic(err)
	}
}
