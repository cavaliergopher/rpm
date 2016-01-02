package yum

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm"
	"time"
)

// PackageEntry is a RPM package as defined in a package repository database.
type PackageEntry struct {
	key           int64
	architecture  string
	archive_size  int64
	checksum      string
	checksum_type string
	epoch         int64
	install_size  int64
	locationhref  string
	name          string
	package_size  int64
	release       string
	version       string
	time_build    int64
}

type PackageEntries []PackageEntry

// String reassembles package metadata to form a standard rpm package name;
// including the package name, version, release and architecture.
func (c PackageEntry) String() string {
	return fmt.Sprintf("%s-%s-%s.%s", c.Name(), c.Version(), c.Release(), c.Architecture())
}

// LocationHref is the location of the package, relative to the parent
// repository.
func (c *PackageEntry) LocationHref() string {
	return c.locationhref
}

func (c *PackageEntry) Checksum() string {
	return c.checksum
}

func (c *PackageEntry) ChecksumType() string {
	return c.checksum_type
}

func (c *PackageEntry) PackageSize() int64 {
	return c.package_size
}

func (c *PackageEntry) InstallSize() int64 {
	return c.install_size
}

func (c *PackageEntry) ArchiveSize() int64 {
	return c.archive_size
}

func (c *PackageEntry) Name() string {
	return c.name
}

func (c *PackageEntry) Version() string {
	return c.version
}

func (c *PackageEntry) Release() string {
	return c.release
}

func (c *PackageEntry) Architecture() string {
	return c.architecture
}

func (c *PackageEntry) Epoch() int64 {
	return c.epoch
}

func (c *PackageEntry) BuildTime() time.Time {
	return time.Unix(c.time_build, 0)
}

func (c *PackageEntry) Requires() rpm.Dependencies {
	return nil
}

func (c *PackageEntry) Provides() rpm.Dependencies {
	return nil
}

func (c *PackageEntry) Conflicts() rpm.Dependencies {
	return nil
}

func (c *PackageEntry) Obsoletes() rpm.Dependencies {
	return nil
}
