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

// TestSHA1Check validates the SHA1 header checksum for each test package
// that carries the RPMTAG_SHA1HEADER (269) tag.
func TestSHA1Check(t *testing.T) {
	files := getTestFiles()

	valid := 0
	for filename, b := range files {
		sig, err := readSigHeader(bytes.NewReader(b))
		if err != nil {
			t.Errorf("Validation error for %s: %v", filename, err)
			continue
		}
		if sig.GetTag(269) == nil {
			continue // skip packages without SHA1HEADER
		}
		if err := SHA1Check(bytes.NewReader(b)); err != nil {
			t.Errorf("Validation error for %s: %v", filename, err)
		} else {
			valid++
		}
	}

	t.Logf("Validated SHA1 header checksum for %d packages", valid)
}

// TestSHA256Check validates the SHA256 header checksum for each test package
// that carries the RPMTAG_SHA256HEADER (273) tag.
func TestSHA256Check(t *testing.T) {
	files := getTestFiles()

	valid := 0
	for filename, b := range files {
		sig, err := readSigHeader(bytes.NewReader(b))
		if err != nil {
			t.Errorf("Validation error for %s: %v", filename, err)
			continue
		}
		if sig.GetTag(273) == nil {
			continue // skip packages without SHA256HEADER
		}
		if err := SHA256Check(bytes.NewReader(b)); err != nil {
			t.Errorf("Validation error for %s: %v", filename, err)
		} else {
			valid++
		}
	}

	t.Logf("Validated SHA256 header checksum for %d packages", valid)
}

// TestSHA3_256Check validates the SHA3_256 header checksum for each test
// package that carries the RPMTAG_SHA3_256HEADER (279) tag.
func TestSHA3_256Check(t *testing.T) {
	files := getTestFiles()

	valid := 0
	for filename, b := range files {
		sig, err := readSigHeader(bytes.NewReader(b))
		if err != nil {
			t.Errorf("Validation error for %s: %v", filename, err)
			continue
		}
		if sig.GetTag(279) == nil {
			continue // skip packages without SHA3_256HEADER
		}
		if err := SHA3_256Check(bytes.NewReader(b)); err != nil {
			t.Errorf("Validation error for %s: %v", filename, err)
		} else {
			valid++
		}
	}

	t.Logf("Validated SHA3_256 header checksum for %d packages", valid)
}

// TestPayloadSHA256Check validates the SHA256 payload digest for each test
// package that carries the RPMTAG_PAYLOADSHA256 (5092) tag in its main header.
func TestPayloadSHA256Check(t *testing.T) {
	files := getTestFiles()

	valid := 0
	for filename, b := range files {
		r := bytes.NewReader(b)
		if _, err := readSigHeader(r); err != nil {
			t.Errorf("Validation error for %s: %v", filename, err)
			continue
		}
		hdr, err := readHeader(r, false)
		if err != nil {
			t.Errorf("Validation error for %s: %v", filename, err)
			continue
		}
		if hdr.GetTag(5092) == nil {
			continue // skip packages without PAYLOADSHA256
		}
		if err := PayloadSHA256Check(bytes.NewReader(b)); err != nil {
			t.Errorf("Validation error for %s: %v", filename, err)
		} else {
			valid++
		}
	}

	t.Logf("Validated SHA256 payload digest for %d packages", valid)
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
