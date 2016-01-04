package rpm

import (
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// keyring loads all gpgkey files in the fixtures directory and returned a
// keyring containing all found entities.
func keyring() (openpgp.KeyRing, error) {
	// read fixtures directory
	dir, err := ioutil.ReadDir("fixtures")
	if err != nil {
		return nil, fmt.Errorf("Error reading fixtures directory: %v", err)
	}

	keyring := make(openpgp.EntityList, 0)
	for _, fi := range dir {
		if strings.HasPrefix(fi.Name(), "RPM-GPG-KEY-") {
			// open gpgkey file
			f, err := os.Open(filepath.Join("fixtures", fi.Name()))
			if err != nil {
				return nil, fmt.Errorf("Error reading %s: %v", fi.Name(), err)
			}
			defer f.Close()

			// decode gpgkey file
			p, err := armor.Decode(f)
			if err != nil {
				return nil, fmt.Errorf("Error decoding %s: %v", fi.Name(), err)
			}

			// extract keys
			el, err := openpgp.ReadKeyRing(p.Body)
			if err != nil {
				return nil, fmt.Errorf("Error reading keyring in %s: %v", fi.Name(), err)
			}

			// append to keyring
			for _, e := range el {
				keyring = append(keyring, e)
			}
		}
	}

	return keyring, nil
}

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
			//os.Remove(path)
		} else {
			valid++
		}
	}

	t.Logf("Validated MD5 checksum for %d packages", valid)
}

func TestGPGCheck(t *testing.T) {
	// load package file paths
	files, err := packages(t)
	if err != nil {
		t.Fatalf("Error listing rpm packages: %v", err)
	}

	// load keys
	keyring, err := keyring()
	if err != nil {
		t.Fatalf(err.Error())
	}

	// check each package
	valid := 0
	for _, path := range files {
		// GPG Check
		f, _ := os.Open(path)
		defer f.Close()

		if signer, err := GPGCheck(f, keyring); err != nil {
			t.Errorf("Validation error for %s: %v", filepath.Base(path), err)
			//os.Remove(path)
		} else {
			t.Logf("%s signed by '%v'", filepath.Base(path), signer)
			valid++
		}
	}

	t.Logf("Validated GPG signature for %d packages", valid)

}
