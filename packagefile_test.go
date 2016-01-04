package rpm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func packages(t *testing.T) ([]string, error) {
	// get a directory full of rpms from RPM_DIR environment variable or
	// failback to ./fixtures
	path := os.Getenv("RPM_DIR")
	if path == "" {
		path = "fixtures"
	}

	// list RPM files
	t.Logf("Loading package files in %s...", path)
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, f := range dir {
		if strings.HasSuffix(f.Name(), ".rpm") {
			files = append(files, filepath.Join("fixtures", f.Name()))
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("No rpm packages found for testing")
	}

	return files, nil
}

func TestReadRPMFile(t *testing.T) {
	// load package file paths
	files, err := packages(t)
	if err != nil {
		t.Fatalf("Error listing rpm packages: %v", err)
	}

	valid := 0
	for _, path := range files {
		// Load package info
		rpm, err := OpenPackageFile(path)
		if err != nil {
			t.Errorf("Error loading RPM file %s: %s", path, err)
			//os.Remove(path)
		} else {
			t.Logf("Loaded package: %v", rpm)
			valid++
		}
	}

	t.Logf("Validated %d RPM files", valid)
}
