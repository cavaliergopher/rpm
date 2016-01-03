package rpm

import (
	"fmt"
	"io"
	"os"
	"time"
)

// A PackageFile is an RPM package definition loaded directly from the pacakge
// file itself.
type PackageFile struct {
	Lead    Lead
	Headers Headers
}

// OpenPackageFile reads a rpm package from the file systems and returns a pointer
// to it.
func OpenPackageFile(path string) (*PackageFile, error) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error opening RPM file: %s", err)
	}
	defer f.Close()

	return ReadPackageFile(f)
}

// ReadPackageFile reads a rpm package file from a stream and returns a pointer
// to it.
func ReadPackageFile(r io.Reader) (*PackageFile, error) {
	// See: http://www.rpm.org/max-rpm/s1-rpm-file-format-rpm-file-format.html
	p := &PackageFile{}

	// read the deprecated "lead"
	lead, err := ReadPackageLead(r)
	if err != nil {
		return nil, err
	}

	p.Lead = *lead

	// parse headers
	p.Headers = make(Headers, 2)

	// read signature and header headers
	for i := 0; i < 2; i++ {
		h, err := ReadPackageHeader(r)
		if err != nil {
			return nil, err
		}

		// add header
		p.Headers[i] = *h
	}

	return p, nil
}

// dependencies translates the given tag values into a slice of package
// relationships such as provides, conflicts, obsoletes and requires.
func (c *PackageFile) dependencies(nevrsTagId, flagsTagId, namesTagId, versionsTagId int) Dependencies {
	// TODO: Implement NEVRS tags

	flgs := c.Headers[1].Indexes.IntsByTag(flagsTagId)
	names := c.Headers[1].Indexes.StringsByTag(namesTagId)
	vers := c.Headers[1].Indexes.StringsByTag(versionsTagId)

	deps := make(Dependencies, len(names))
	for i := 0; i < len(names); i++ {
		deps[i] = NewDependency(int(flgs[i]), names[i], 0, vers[i], "")
	}

	return deps
}

// String reassembles package metadata to form a standard rpm package name;
// including the package name, version, release and architecture.
func (c *PackageFile) String() string {
	return fmt.Sprintf("%s-%s-%s.%s", c.Name(), c.Version(), c.Release(), c.Architecture())
}

// For tag definitions, see:
// https://github.com/rpm-software-management/rpm/blob/master/lib/rpmtag.h

func (c *PackageFile) Name() string {
	return c.Headers[1].Indexes.StringByTag(1000)
}

func (c *PackageFile) Version() string {
	return c.Headers[1].Indexes.StringByTag(1001)
}

func (c *PackageFile) Release() string {
	return c.Headers[1].Indexes.StringByTag(1002)
}

func (c *PackageFile) Epoch() int64 {
	return c.Headers[1].Indexes.IntByTag(1003)
}

func (c *PackageFile) Requires() Dependencies {
	return c.dependencies(5041, 1048, 1049, 1050)
}

func (c *PackageFile) Provides() Dependencies {
	return c.dependencies(5042, 1112, 1047, 1113)
}

func (c *PackageFile) Conflicts() Dependencies {
	return c.dependencies(5044, 1053, 1054, 1055)
}

func (c *PackageFile) Obsoletes() Dependencies {
	return c.dependencies(5043, 1114, 1090, 1115)
}

func (c *PackageFile) Files() []string {
	ixs := c.Headers[1].Indexes.IntsByTag(1116)
	names := c.Headers[1].Indexes.StringsByTag(1117)
	dirs := c.Headers[1].Indexes.StringsByTag(1118)

	files := make([]string, len(names))
	for i := 0; i < len(names); i++ {
		files[i] = dirs[ixs[i]] + names[i]
	}

	return files
}

func (c *PackageFile) Summary() []string {
	return c.Headers[1].Indexes.StringsByTag(1004)
}

func (c *PackageFile) Description() []string {
	return c.Headers[1].Indexes.StringsByTag(1005)
}

func (c *PackageFile) BuildTime() time.Time {
	return c.Headers[1].Indexes.TimeByTag(1006)
}

func (c *PackageFile) BuildHost() string {
	return c.Headers[1].Indexes.StringByTag(1007)
}

func (c *PackageFile) InstallTime() time.Time {
	return c.Headers[1].Indexes.TimeByTag(1008)
}

func (c *PackageFile) Size() int64 {
	return c.Headers[1].Indexes.IntByTag(1009)
}

func (c *PackageFile) Distribution() string {
	return c.Headers[1].Indexes.StringByTag(1010)
}

func (c *PackageFile) Vendor() string {
	return c.Headers[1].Indexes.StringByTag(1011)
}

func (c *PackageFile) GIFImage() []byte {
	return c.Headers[1].Indexes.BytesByTag(1012)
}

func (c *PackageFile) XPMImage() []byte {
	return c.Headers[1].Indexes.BytesByTag(1013)
}

func (c *PackageFile) License() string {
	return c.Headers[1].Indexes.StringByTag(1014)
}

func (c *PackageFile) PackageFiler() string {
	return c.Headers[1].Indexes.StringByTag(1015)
}

func (c *PackageFile) Groups() []string {
	return c.Headers[1].Indexes.StringsByTag(1016)
}

func (c *PackageFile) ChangeLog() []string {
	return c.Headers[1].Indexes.StringsByTag(1017)
}

func (c *PackageFile) Source() []string {
	return c.Headers[1].Indexes.StringsByTag(1018)
}

func (c *PackageFile) Patch() []string {
	return c.Headers[1].Indexes.StringsByTag(1019)
}

func (c *PackageFile) URL() string {
	return c.Headers[1].Indexes.StringByTag(1020)
}

func (c *PackageFile) OperatingSystem() string {
	return c.Headers[1].Indexes.StringByTag(1021)
}

func (c *PackageFile) Architecture() string {
	return c.Headers[1].Indexes.StringByTag(1022)
}

func (c *PackageFile) PreInstallScript() string {
	return c.Headers[1].Indexes.StringByTag(1023)
}

func (c *PackageFile) PostInstallScript() string {
	return c.Headers[1].Indexes.StringByTag(1024)
}

func (c *PackageFile) PreUninstallScript() string {
	return c.Headers[1].Indexes.StringByTag(1025)
}

func (c *PackageFile) PostUninstallScript() string {
	return c.Headers[1].Indexes.StringByTag(1026)
}

func (c *PackageFile) OldFilenames() []string {
	return c.Headers[1].Indexes.StringsByTag(1027)
}

func (c *PackageFile) Icon() []byte {
	return c.Headers[1].Indexes.BytesByTag(1043)
}

func (c *PackageFile) SourceRPM() string {
	return c.Headers[1].Indexes.StringByTag(1044)
}

func (c *PackageFile) RPMVersion() string {
	return c.Headers[1].Indexes.StringByTag(1064)
}

func (c *PackageFile) Platform() string {
	return c.Headers[1].Indexes.StringByTag(1132)
}
