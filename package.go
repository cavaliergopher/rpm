package rpm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"
)

// A Package is an rpm package file.
type Package struct {
	Lead      Lead
	Signature Header
	Header    Header
}

var _ Version = &Package{}

// Read reads an rpm package from r.
//
// When this function returns, the reader will be positioned at the start of the
// package payload. Use Package.PayloadFormat and Package.PayloadCompression to
// determine how to decompress and unarchive the payload.
func Read(r io.Reader) (*Package, error) {
	lead, err := readLead(r)
	if err != nil {
		return nil, err
	}
	sig, err := readHeader(r, true)
	if err != nil {
		return nil, err
	}
	hdr, err := readHeader(r, false)
	if err != nil {
		return nil, err
	}
	return &Package{
		Lead:      *lead,
		Signature: *sig,
		Header:    *hdr,
	}, nil
}

// Open opens an rpm package from the file system.
//
// Once the package headers are read, the underlying reader is closed and cannot
// be used to read the package payload. To read the package payload, open the
// package with os.Open and read the headers with Read. You may then use the
// same reader to read the payload.
func Open(name string) (*Package, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(bufio.NewReader(f))
}

// dependencies translates the given tag values into a slice of package
// relationships such as provides, conflicts, obsoletes and requires.
func (c *Package) dependencies(nevrsTagID, flagsTagID, namesTagID, versionsTagID int) []Dependency {
	// TODO: Implement NEVRS tags
	// TODO: error handling
	flgs := c.Header.GetTag(flagsTagID).Int64Slice()
	names := c.Header.GetTag(namesTagID).StringSlice()
	vers := c.Header.GetTag(versionsTagID).StringSlice()
	deps := make([]Dependency, len(names))
	for i := 0; i < len(names); i++ {
		deps[i] = &dependency{
			flags:   int(flgs[i]),
			name:    names[i],
			version: vers[i],
		}
	}
	return deps
}

// String returns the package identifier in the form
// '[name]-[version]-[release].[architecture]'.
func (c *Package) String() string {
	return fmt.Sprintf("%s-%s-%s.%s", c.Name(), c.Version(), c.Release(), c.Architecture())
}

func (c *Package) GPGSignature() GPGSignature {
	return c.Signature.GetTag(1002).Bytes()
}

// For tag definitions, see:
// https://github.com/rpm-software-management/rpm/blob/master/lib/rpmtag.h#L61

func (c *Package) Name() string {
	return c.Header.GetTag(1000).String()
}

func (c *Package) Version() string {
	return c.Header.GetTag(1001).String()
}

func (c *Package) Release() string {
	return c.Header.GetTag(1002).String()
}

func (c *Package) Epoch() int {
	return int(c.Header.GetTag(1003).Int64())
}

func (c *Package) Requires() []Dependency {
	return c.dependencies(5041, 1048, 1049, 1050)
}

func (c *Package) Provides() []Dependency {
	return c.dependencies(5042, 1112, 1047, 1113)
}

func (c *Package) Conflicts() []Dependency {
	return c.dependencies(5044, 1053, 1054, 1055)
}

func (c *Package) Obsoletes() []Dependency {
	return c.dependencies(5043, 1114, 1090, 1115)
}

func (c *Package) Suggests() []Dependency {
	return c.dependencies(5059, 5051, 5049, 5050)
}

func (c *Package) Enhances() []Dependency {
	return c.dependencies(5061, 5057, 5055, 5056)
}

func (c *Package) Recommends() []Dependency {
	return c.dependencies(5058, 5048, 5046, 5047)
}

func (c *Package) Supplements() []Dependency {
	return c.dependencies(5060, 5051, 5052, 5053)
}

// Files returns file information for each file that is installed by this RPM
// package.
func (c *Package) Files() []FileInfo {
	ixs := c.Header.GetTag(1116).Int64Slice()
	names := c.Header.GetTag(1117).StringSlice()
	dirs := c.Header.GetTag(1118).StringSlice()
	modes := c.Header.GetTag(1030).Int64Slice()
	sizes := c.Header.GetTag(1028).Int64Slice()
	times := c.Header.GetTag(1034).Int64Slice()
	flags := c.Header.GetTag(1037).Int64Slice()
	owners := c.Header.GetTag(1039).StringSlice()
	groups := c.Header.GetTag(1040).StringSlice()
	digests := c.Header.GetTag(1035).StringSlice()
	linknames := c.Header.GetTag(1036).StringSlice()
	a := make([]FileInfo, len(names))
	for i := 0; i < len(names); i++ {
		a[i] = FileInfo{
			name:     dirs[ixs[i]] + names[i],
			mode:     fileModeFromInt64(modes[i]),
			size:     sizes[i],
			modTime:  time.Unix(times[i], 0),
			flags:    flags[i],
			owner:    owners[i],
			group:    groups[i],
			digest:   digests[i],
			linkname: linknames[i],
		}
	}
	return a
}

// fileModeFromInt64 converts the 16 bit value returned from a typical
// unix/linux stat call to the bitmask that go uses to produce an os
// neutral representation.  It is incorrect to just cast the 16 bit
// value directly to a os.FileMode.  The result of stat is 4 bits to
// specify the type of the object, this is a value in the range 0 to
// 15, rather than a bitfield, 3 bits to note suid, sgid and sticky,
// and 3 sets of 3 bits for rwx permissions for user, group and other.
// An os.FileMode has the same 9 bits for permissions, but rather than
// using an enum for the type it has individual bits.  As a concrete
// example, a block device has the 1<<26 bit set (os.ModeDevice) in
// the os.FileMode, but has type 0x6000 (syscall.S_IFBLK). A regular
// file is represented in os.FileMode by not having any of the bits in
// os.ModeType set (i.e. is not a directory, is not a symlink, is not
// a named pipe...) whilst a regular file has value syscall.S_IFREG
// (0x8000) in the mode field from stat.
func fileModeFromInt64(mode int64) os.FileMode {
	fm := os.FileMode(mode & 0777)
	switch mode & syscall.S_IFMT {
	case syscall.S_IFBLK:
		fm |= os.ModeDevice
	case syscall.S_IFCHR:
		fm |= os.ModeDevice | os.ModeCharDevice
	case syscall.S_IFDIR:
		fm |= os.ModeDir
	case syscall.S_IFIFO:
		fm |= os.ModeNamedPipe
	case syscall.S_IFLNK:
		fm |= os.ModeSymlink
	case syscall.S_IFREG:
		// nothing to do
	case syscall.S_IFSOCK:
		fm |= os.ModeSocket
	}
	if mode&syscall.S_ISGID != 0 {
		fm |= os.ModeSetgid
	}
	if mode&syscall.S_ISUID != 0 {
		fm |= os.ModeSetuid
	}
	if mode&syscall.S_ISVTX != 0 {
		fm |= os.ModeSticky
	}
	return fm
}

func (c *Package) Summary() string {
	return strings.Join(c.Header.GetTag(1004).StringSlice(), "\n")
}

func (c *Package) Description() string {
	return strings.Join(c.Header.GetTag(1005).StringSlice(), "\n")
}

func (c *Package) BuildTime() time.Time {
	return time.Unix(c.Header.GetTag(1006).Int64(), 0)
}

func (c *Package) BuildHost() string {
	return c.Header.GetTag(1007).String()
}

func (c *Package) InstallTime() time.Time {
	return time.Unix(c.Header.GetTag(1008).Int64(), 0)
}

// Size specifies the disk space consumed by installation of the package.
func (c *Package) Size() uint64 {
	return uint64(c.Header.GetTag(1009).Int64())
}

// ArchiveSize specifies the size of the archived payload of the package in
// bytes.
func (c *Package) ArchiveSize() uint64 {
	if i := uint64(c.Signature.GetTag(1007).Int64()); i > 0 {
		return i
	}

	return uint64(c.Header.GetTag(1046).Int64())
}

func (c *Package) Distribution() string {
	return c.Header.GetTag(1010).String()
}

func (c *Package) Vendor() string {
	return c.Header.GetTag(1011).String()
}

func (c *Package) GIFImage() []byte {
	return c.Header.GetTag(1012).Bytes()
}

func (c *Package) XPMImage() []byte {
	return c.Header.GetTag(1013).Bytes()
}

func (c *Package) License() string {
	return c.Header.GetTag(1014).String()
}

func (c *Package) Packager() string {
	return c.Header.GetTag(1015).String()
}

func (c *Package) Groups() []string {
	return c.Header.GetTag(1016).StringSlice()
}

func (c *Package) ChangeLog() []string {
	return c.Header.GetTag(1017).StringSlice()
}

func (c *Package) Source() []string {
	return c.Header.GetTag(1018).StringSlice()
}

func (c *Package) Patch() []string {
	return c.Header.GetTag(1019).StringSlice()
}

func (c *Package) URL() string {
	return c.Header.GetTag(1020).String()
}

func (c *Package) OperatingSystem() string {
	return c.Header.GetTag(1021).String()
}

func (c *Package) Architecture() string {
	return c.Header.GetTag(1022).String()
}

func (c *Package) PreInstallScript() string {
	return c.Header.GetTag(1023).String()
}

func (c *Package) PostInstallScript() string {
	return c.Header.GetTag(1024).String()
}

func (c *Package) PreUninstallScript() string {
	return c.Header.GetTag(1025).String()
}

func (c *Package) PostUninstallScript() string {
	return c.Header.GetTag(1026).String()
}

func (c *Package) OldFilenames() []string {
	return c.Header.GetTag(1027).StringSlice()
}

func (c *Package) Icon() []byte {
	return c.Header.GetTag(1043).Bytes()
}

func (c *Package) SourceRPM() string {
	return c.Header.GetTag(1044).String()
}

func (c *Package) RPMVersion() string {
	return c.Header.GetTag(1064).String()
}

func (c *Package) Platform() string {
	return c.Header.GetTag(1132).String()
}

// PayloadFormat returns the name of the format used for the package payload.
// Typically cpio.
func (c *Package) PayloadFormat() string {
	return c.Header.GetTag(1124).String()
}

// PayloadCompression returns the name of the compression used for the package
// payload. Typically xz.
func (c *Package) PayloadCompression() string {
	return c.Header.GetTag(1125).String()
}

// Sort sorts a slice of packages lexically by name ascending and then by
// version descending. Version is evaluated first by epoch, then by version
// string, then by release.
func Sort(x []*Package) { sort.Sort(PackageSlice(x)) }

// PackageSlice implements sort.Interface for a slice of packages. Packages are
// sorted lexically by name ascending and then by version descending. Version is
// evaluated first by epoch, then by version string, then by release.
type PackageSlice []*Package

// Sort is a convenience method: x.Sort() calls sort.Sort(x).
func (x PackageSlice) Sort() { sort.Sort(x) }

func (x PackageSlice) Len() int { return len(x) }

func (x PackageSlice) Less(i, j int) bool {
	a, b := x[i].Name(), x[j].Name()
	if a == b {
		return Compare(x[i], x[j]) == 1
	}
	return a < b
}

func (x PackageSlice) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

var _ sort.Interface = PackageSlice{}
