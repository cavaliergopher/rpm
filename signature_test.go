package rpm

import (
	"bytes"
	"fmt"
	"log"
	"os"
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

// ExampleMD5Check validates a local RPM package using the MD5 checksum value
// specified in the package header.
func ExampleMD5Check() {
	f, err := os.Open("testdata/centos-release-7-2.1511.el7.centos.2.10.x86_64.rpm")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Print(err)
		}
	}()

	if err := MD5Check(f); err == nil {
		fmt.Printf("Package passed checksum validation\n")
	} else if err == ErrMD5CheckFailed {
		fmt.Printf("Package failed checksum validation\n")
	} else {
		log.Fatal(err)
	}

	// Output: Package passed checksum validation
}
