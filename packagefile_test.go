package rpm

import (
	"bytes"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var testFiles map[string][]byte

func getTestFiles() map[string][]byte {
	if testFiles != nil {
		return testFiles
	}

	// get a directory full of rpms from RPM_DIR environment variable or
	// failback to ./testdata
	path := os.Getenv("RPM_DIR")
	if path == "" {
		path = "testdata"
	}

	// list RPM files
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	files := make([]string, 0)
	for _, f := range dir {
		if strings.HasSuffix(f.Name(), ".rpm") {
			files = append(files, filepath.Join(path, f.Name()))
		}
	}

	if len(files) == 0 {
		panic("No rpm packages found for testing")
	}

	testFiles = make(map[string][]byte, len(files))
	for _, filename := range files {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		testFiles[filename] = b
	}

	return testFiles
}

func TestReadRPMFile(t *testing.T) {
	// load package file paths
	files := getTestFiles()

	valid := 0
	for path, b := range files {
		// Load package info
		rpm, err := ReadPackageFile(bytes.NewReader(b))
		if err != nil {
			t.Errorf("Error loading RPM file %s: %s", path, err)
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

func TestPackageFiles(t *testing.T) {
	names := []string{
		"/etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7",
		"/etc/yum.repos.d/epel-testing.repo",
		"/etc/yum.repos.d/epel.repo",
		"/usr/lib/rpm/macros.d/macros.epel",
		"/usr/lib/systemd/system-preset/90-epel.preset",
		"/usr/share/doc/epel-release-7",
		"/usr/share/doc/epel-release-7/GPL",
	}
	modes := []int64{0644, 0644, 0644, 0644, 0644, 0755, 0644}
	sizes := []int64{1662, 1056, 957, 41, 2813, 4096, 18385}
	owners := []string{"root", "root", "root", "root", "root", "root", "root"}
	groups := []string{"root", "root", "root", "root", "root", "root", "root"}
	modtimes := []time.Time{
		time.Unix(1416932629, 0),
		time.Unix(1416932629, 0),
		time.Unix(1416932629, 0),
		time.Unix(1416932629, 0),
		time.Unix(1416932629, 0),
		time.Unix(1416932778, 0),
		time.Unix(1416932629, 0),
	}

	path := "./testdata/epel-release-7-5.noarch.rpm"

	p, err := OpenPackageFile(path)
	if err != nil {
		t.Fatalf("Error opening %s: %v", path, err)
	}

	files := p.Files()
	if len(files) != len(names) {
		t.Fatalf("expected %v files in RPM package but got %v", len(names), len(files))
	}

	for i, fi := range files {
		name := fi.Name()
		if name != names[i] {
			t.Errorf("expected file %v with name %v but got %v", i, names[i], name)
			continue
		}

		if mode := int64(fi.Mode().Perm()); mode != modes[i] {
			t.Errorf("expected mode %v but got %v for %v", modes[i], mode, name)
		}

		if size := fi.Size(); size != sizes[i] {
			t.Errorf("expected size %v but got %v for %v", sizes[i], size, name)
		}

		if owner := fi.Owner(); owner != owners[i] {
			t.Errorf("expected owner %v but got %v for %v", owners[i], owner, name)
		}

		if group := fi.Group(); group != groups[i] {
			t.Errorf("expected group %v but got %v for %v", groups[i], group, name)
		}

		if modtime := fi.ModTime(); modtime != modtimes[i] {
			t.Errorf("expected modtime %v but got %v for %v", modtimes[i], modtime.Unix(), name)
		}
	}
}

func TestByteTags(t *testing.T) {
	tests := []struct {
		Path            string
		GPGSignatureCRC uint32
	}{
		{
			Path:            "testdata/centos-release-6-0.el6.centos.5.i686.rpm",
			GPGSignatureCRC: 1788312322,
		},
		{
			Path:            "testdata/centos-release-6-0.el6.centos.5.x86_64.rpm",
			GPGSignatureCRC: 3194808352,
		},
		{
			Path:            "testdata/centos-release-7-2.1511.el7.centos.2.10.x86_64.rpm",
			GPGSignatureCRC: 3466078337,
		},
		{
			Path:            "testdata/epel-release-7-5.noarch.rpm",
			GPGSignatureCRC: 2817187108,
		},
	}
	for _, test := range tests {
		p, err := OpenPackageFile(test.Path)
		if err != nil {
			t.Errorf("error opening %v: %v", test.Path, err)
			continue
		}

		if crc := crc32.ChecksumIEEE(p.GPGSignature()); crc != test.GPGSignatureCRC {
			t.Errorf("expected GPG Signature CRC %v, got %v for %v", test.GPGSignatureCRC, crc, test.Path)
		}
	}
}

func BenchmarkPackageOpens(b *testing.B) {
	files := getTestFiles()
	// parse packages from byte arrays b.N times
	var V interface{}
	for n := 0; n < b.N; n++ {
		for _, b := range files {
			p, err := ReadPackageFile(bytes.NewReader(b))
			if err != nil {
				panic(err)
			}

			V = p
		}

		X = V
	}
}
