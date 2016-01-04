package rpm_test

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"os"
)

// ExampleGPGCheck reads a public GPG key and uses it to validate the signature
// of a local rpm package.
func ExampleGPGCheck() {
	// open gpgkey file (typically in /etc/pki/rpm-gpg)
	f_gpgkey, err := os.Open("fixtures/RPM-GPG-KEY-CentOS-7")
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
	f_rpm, err := os.Open("fixtures/centos-release-7-2.1511.el7.centos.2.10.x86_64.rpm")
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

	// Output: Package signed by 'CentOS-7 Key (CentOS 7 Official Signing Key) <security@centos.org>'
}
