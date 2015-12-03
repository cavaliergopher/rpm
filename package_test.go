package rpm

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadRPMFile(t *testing.T) {
	dir := "./rpms"

	// list RPM files
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf(err.Error())
	}

	valid := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".rpm") {
			path := filepath.Join(dir, f.Name())
			_, err := OpenPackage(path)
			if err != nil {
				t.Errorf("Error loading RPM file %s: %s", f.Name(), err)
			} else {
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
