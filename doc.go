/*
Package rpm implements the rpm package file format.

	package main

	import (
		"fmt"
		"log"

		"github.com/cavaliergopher/rpm"
	)

	func main() {
		pkg, err := rpm.Open("golang-1.17.2-1.el7.x86_64.rpm")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Package:", pkg)
		fmt.Println("Summary:", pkg.Summary())

		// Output:
		// Package: golang-1.17.2-1.el7.x86_64
		// Summary: The Go Programming Language
	}

For more information about the rpm file format, see:

http://ftp.rpm.org/max-rpm/s1-rpm-file-format-rpm-file-format.html

Packages are composed of two headers: the Signature header and the "Header"
header. Each contains key-value pairs called tags. Tags map an integer key to a
value whose data type will be one of the TagType types. Tag values can be
decoded with the appropriate Tag method for the data type.

Many known tags are available as Package methods. For example, RPMTAG_NAME and
RPMTAG_BUILDTIME are available as Package.Name and Package.BuildTime
respectively.

	fmt.Println(pkg.Name(), pkg.BuildTime())

Tags can be retrieved and decoded from the Signature or Header headers directly
using Header.GetTag and their tag identifier.

	const (
		RPMTagName      = 1000
		RPMTagBuidlTime = 1006
	)

	fmt.Println(
		pkg.Header.GetTag(RPMTagName).String()),
		time.Unix(pkg.Header.GetTag(RPMTagBuildTime).Int64(), 0),
	)

Header.GetTag and all Tag methods will return a zero value if the header or the
tag do not exist, or if the tag has a different data type.

You may enumerate all tags in a header with Header.Tags:

	for id, tag := range pkg.Header.Tags {
		fmt.Println(id, tag.Type, tag.Value)
	}

Comparing versions

In the rpm ecosystem, package versions are compared using EVR; epoch, version,
release. Versions may be compared using the Compare function.

	if rpm.Compare(pkgA, pkgB) == 1 {
		fmt.Println("A is more recent than B")
	}

Packages may be be sorted using the PackageSlice type which implements
sort.Interface. Packages are sorted lexically by name ascending and then by
version descending. Version is evaluated first by epoch, then by version string,
then by release.

	sort.Sort(PackageSlice(pkgs))

The Sort function is provided for your convenience.

	rpm.Sort(pkgs)

Checksum validation

Packages may be validated using MD5Check or GPGCheck. See the example for each
function.

Extracting files

The payload of an rpm package is typically archived in cpio format and
compressed with xz. To decompress and unarchive an rpm payload, the reader that
read the rpm package headers will be positioned at the beginning of the payload
and can be reused with the appropriate Go packages for the rpm payload format.

You can check the archive format with Package.PayloadFormat and the compression
algorithm with Package.PayloadCompression.

For the cpio archive format, the following package is recommended:

https://github.com/cavaliergopher/cpio

For xz compression, the following package is recommended:

https://github.com/ulikunitz/xz

See README.md for a working example of extracting files from a cpio/xz rpm
package using these packages.

Example programs

See cmd/rpmdump and cmd/rpminfo for example programs that emulate tools from the
rpm ecosystem.
*/
package rpm
