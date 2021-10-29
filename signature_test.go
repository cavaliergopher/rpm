package rpm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMD5Check(t *testing.T) {
	files := getTestFiles()

	valid := 0
	for filename, b := range files {
		if err := MD5Check(bytes.NewReader(b)); err != nil {
			t.Errorf("Validation error for %s: %v", filename, err)
		} else {
			valid++
		}
	}

	t.Logf("Validated MD5 checksum for %d packages", valid)
}

func TestGPGCheck(t *testing.T) {
	// read testdata directory
	dir, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatalf(err.Error())
	}

	// filter for gpgkey files
	keyfiles := make([]string, 0)
	for _, fi := range dir {
		if strings.HasPrefix(fi.Name(), "RPM-GPG-KEY-") {
			keyfiles = append(keyfiles, filepath.Join("testdata", fi.Name()))
		}
	}

	// build keyring
	keyring, err := OpenKeyRing(keyfiles...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// load package file paths
	files := getTestFiles()

	// check each package
	valid := 0
	for filename, b := range files {
		if signer, err := GPGCheck(bytes.NewReader(b), keyring); err != nil {
			t.Errorf("Validation error for %s: %v", filepath.Base(filename), err)
		} else {
			t.Logf("%s signed by '%v'", filepath.Base(filename), signer)
			valid++
		}
	}

	t.Logf("Validated GPG signature for %d packages", valid)
}

// ExampleGPGCheck reads a public GPG key and uses it to validate the signature
// of a local rpm package.
func ExampleGPGCheck() {
	// read public key from gpgkey file
	keyring, err := OpenKeyRing("testdata/RPM-GPG-KEY-CentOS-7")
	if err != nil {
		log.Fatal(err)
	}

	// open a rpm package for reading
	f, err := os.Open("testdata/centos-release-7-2.1511.el7.centos.2.10.x86_64.rpm")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// validate gpg signature
	if signer, err := GPGCheck(f, keyring); err == nil {
		fmt.Printf("Package signed by '%s'\n", signer)
	} else if err == ErrGPGCheckFailed {
		fmt.Printf("Package failed GPG signature validation\n")
	} else {
		log.Fatal(err)
	}

	// Output: Package signed by 'CentOS-7 Key (CentOS 7 Official Signing Key) <security@centos.org>'
}

// ExampleMD5Check validates a local rpm package named using the MD5 checksum
// value specified in the package header.
func ExampleMD5Check() {
	// open a rpm package for reading
	f, err := os.Open("testdata/centos-release-7-2.1511.el7.centos.2.10.x86_64.rpm")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// validate md5 checksum
	if err := MD5Check(f); err == nil {
		fmt.Printf("Package passed checksum validation\n")
	} else if err == ErrMD5CheckFailed {
		fmt.Printf("Package failed checksum validation\n")
	} else {
		log.Fatal(err)
	}

	// Output: Package passed checksum validation
}
