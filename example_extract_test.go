package rpm

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cavaliercoder/go-cpio"
)

// ExampleExtractArchive demonstrates how to use go-rpm with the compress/gzip
// library and go-cpio to extract all files from a package archive.
func ExampleExtractArchive() {
	// create destination directory
	path, err := ioutil.TempDir("", "rpm-example-")
	if err != nil {
		panic(err)
	}

	// open a rpm package for reading
	r, err := os.Open("testdata/centos-release-4-0.1.x86_64.rpm")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	// read package info
	_, err = ReadPackageFile(r)
	if err != nil {
		panic(err)
	}

	// deflate payload with gzip (newer packages need xz)
	gzr, err := gzip.NewReader(r)
	if err != nil {
		panic(err)
	}
	defer gzr.Close()

	// read each file in payload with cpio
	cpr := cpio.NewReader(gzr)
	for {
		fi, err := cpr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		dst := filepath.Join(path, fi.Name)

		// create parent directories
		dir := dst
		if !fi.IsDir() {
			dir = filepath.Dir(dst)
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}

		// copy file contents
		if !fi.IsDir() {
			fmt.Printf("extracting %s...\n", fi.Name[1:])
			f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			_, err = io.Copy(f, cpr)
			if err != nil {
				panic(err)
			}
			f.Close()
		}

		// clean up
		os.RemoveAll(path)
	}

	// Output:
	// extracting /etc/issue...
	// extracting /etc/issue.net...
	// extracting /etc/redhat-release...
	// extracting /usr/share/doc/centos-release-4/EULA...
	// extracting /usr/share/doc/centos-release-4/GPL...
	// extracting /usr/share/doc/centos-release-4/RELEASE-NOTES-en...
	// extracting /usr/share/doc/centos-release-4/RELEASE-NOTES-en.html...
	// extracting /usr/share/doc/centos-release-4/RPM-GPG-KEY...
	// extracting /usr/share/doc/centos-release-4/RPM-GPG-KEY-centos4...
	// extracting /usr/share/doc/centos-release-4/autorun-template...
	// extracting /usr/share/doc/centos-release-4/centosdocs-man.css...
	// extracting /usr/share/eula/eula.en_US...
	// extracting /usr/share/firstboot/modules/eula.py...
	// extracting /usr/share/firstboot/modules/eula.pyc...
	// extracting /var/lib/supportinfo...
}
