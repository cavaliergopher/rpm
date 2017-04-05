package rpm

import (
	"bytes"
	"io/ioutil"
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
	keyring, err := KeyRingFromFiles(keyfiles)
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
