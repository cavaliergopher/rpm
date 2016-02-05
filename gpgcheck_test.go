package rpm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMD5Check(t *testing.T) {
	// load package file paths
	files, err := packages(t)
	if err != nil {
		t.Fatalf("Error listing rpm packages: %v", err)
	}

	valid := 0
	for _, path := range files {
		// MD5 Check
		f, _ := os.Open(path)
		defer f.Close()

		if err := MD5Check(f); err != nil {
			t.Errorf("Validation error for %s: %v", f.Name(), err)
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
	files, err := packages(t)
	if err != nil {
		t.Fatalf("Error listing rpm packages: %v", err)
	}

	// check each package
	valid := 0
	for _, path := range files {
		// GPG Check
		f, _ := os.Open(path)
		defer f.Close()

		if signer, err := GPGCheck(f, keyring); err != nil {
			t.Errorf("Validation error for %s: %v", filepath.Base(path), err)
		} else {
			t.Logf("%s signed by '%v'", filepath.Base(path), signer)
			valid++
		}
	}

	t.Logf("Validated GPG signature for %d packages", valid)
}
