package rpm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadRPMFile(t *testing.T) {
	// get a directory full of rpms from RPM_DIR environment variable or
	// failback to ./fixtures
	dir := os.Getenv("RPM_DIR")
	if dir == "" {
		dir = "fixtures"
	}

	// list RPM files
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf(err.Error())
	}

	valid := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".rpm") {
			path := filepath.Join(dir, f.Name())

			// MD5 Check
			f, _ := os.Open(path)
			defer f.Close()

			if err := MD5Check(f); err != nil {
				t.Errorf("Validation error for %s: %v", f.Name(), err)
				//os.Remove(path)
			}
			f.Close()

			// Load package info
			rpm, err := OpenPackageFile(path)
			if err != nil {
				t.Errorf("Error loading RPM file %s: %s", f.Name(), err)
				//os.Remove(path)
			} else {
				t.Logf("Loaded package: %v", rpm)
				valid++
			}
		}
	}

	if valid == 0 {
		t.Errorf("No RPM files found for testing with in %s", dir)
	} else {
		t.Logf("Validated %d RPM files", valid)
	}
}
