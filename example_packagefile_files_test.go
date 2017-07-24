package rpm_test

import (
	"fmt"

	"github.com/cavaliercoder/go-rpm"
)

// Lists all the files in a RPM package.
func ExamplePackageFile_Files() {
	// open a package file
	pkg, err := rpm.OpenPackageFile("./testdata/epel-release-7-5.noarch.rpm")
	if err != nil {
		panic(err)
	}

	// list each file
	files := pkg.Files()
	fmt.Printf("total %v\n", len(files))
	for _, fi := range files {
		fmt.Printf("%v %v %v %5v %v %v\n",
			fi.Mode().Perm(),
			fi.Owner(),
			fi.Group(),
			fi.Size(),
			fi.ModTime().UTC().Format("Jan 02 15:04"),
			fi.Name())
	}

	// Output:
	// total 7
	// -rw-r--r-- root root  1662 Nov 25 16:23 /etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7
	// -rw-r--r-- root root  1056 Nov 25 16:23 /etc/yum.repos.d/epel-testing.repo
	// -rw-r--r-- root root   957 Nov 25 16:23 /etc/yum.repos.d/epel.repo
	// -rw-r--r-- root root    41 Nov 25 16:23 /usr/lib/rpm/macros.d/macros.epel
	// -rw-r--r-- root root  2813 Nov 25 16:23 /usr/lib/systemd/system-preset/90-epel.preset
	// -rwxr-xr-x root root  4096 Nov 25 16:26 /usr/share/doc/epel-release-7
	// -rw-r--r-- root root 18385 Nov 25 16:23 /usr/share/doc/epel-release-7/GPL
}
