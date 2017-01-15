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
	// failback to ./testdata
	path := os.Getenv("RPM_DIR")
	if path == "" {
		path = "testdata"
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
			files = append(files, filepath.Join(path, f.Name()))
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

func TestReadRPMDirectory(t *testing.T) {
	expected := 10
	packages, err := OpenPackageFiles("./testdata")
	if err != nil {
		t.Fatalf("Error reading RPMs in directory: %v", err)
	}

	// count packages
	if len(packages) != expected {
		t.Errorf("Expected %d packages in directory; got %d", expected, len(packages))
	}
}

func TestChecksum(t *testing.T) {
	path := "./testdata/epel-release-7-5.noarch.rpm"
	expected := "d6f332ed157de1d42058ec785b392a1cc4b5836c27830af8fbf083cce29ef0ab"

	p, err := OpenPackageFile(path)
	if err != nil {
		t.Fatalf("Error opening %s: %v", path, err)
	}

	sum, err := p.Checksum()
	if err != nil {
		t.Errorf("Error validating checksum for %s: %v", path, err)
	} else {
		if sum != expected {
			t.Errorf("Expected sum %s for %s; got %s", expected, path, sum)
		}
	}
}

func TestFileModes(t *testing.T) {
	expectedModes := map[string]os.FileMode{
		"/etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7":           0644,
		"/etc/yum.repos.d/epel-testing.repo":            0644,
		"/etc/yum.repos.d/epel.repo":                    0644,
		"/usr/lib/rpm/macros.d/macros.epel":             0644,
		"/usr/lib/systemd/system-preset/90-epel.preset": 0644,
		"/usr/share/doc/epel-release-7":                 0755,
		"/usr/share/doc/epel-release-7/GPL":             0644,
	}

	path := "./testdata/epel-release-7-5.noarch.rpm"

	p, err := OpenPackageFile(path)
	if err != nil {
		t.Fatalf("Error opening %s: %v", path, err)
	}

	names := p.Files()
	modes := p.Modes()

	if len(names) != len(modes) {
		t.Fatal("Mismatched slice lenghts for Files (len=%v) and Modes (len=%v)", len(names), len(modes))
	}

	for i, name := range names {
		mode := modes[i].Perm()

		m, found := expectedModes[name]
		if !found {
			t.Errorf("unexpected file found in RPM: %v", name)
			continue
		}

		if m != mode {
			t.Errorf("expected %v but got %v for %v", m, mode, name)
			continue
		}
	}
}
