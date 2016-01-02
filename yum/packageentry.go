package yum

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm"
	"time"
)

// PackageEntry is a RPM package as defined in a yum repository database.
type PackageEntry struct {
	db *PrimaryDatabase

	key           int
	architecture  string
	archive_size  int64
	checksum      string
	checksum_type string
	epoch         int
	install_size  int64
	locationhref  string
	name          string
	package_size  int64
	release       string
	version       string
	time_build    int64
}

// PackageEntries is a slice of PackageEntry structs.
type PackageEntries []PackageEntry

// String reassembles package metadata to form a standard rpm package name;
// including the package name, version, release and architecture.
func (c PackageEntry) String() string {
	return fmt.Sprintf("%s-%s-%s.%s", c.Name(), c.Version(), c.Release(), c.Architecture())
}

// Key is the unique identifier of a package within a primary_db
func (c *PackageEntry) Key() int {
	return c.key
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

func (c *PackageEntry) Epoch() int {
	return c.epoch
}

func (c *PackageEntry) BuildTime() time.Time {
	return time.Unix(c.time_build, 0)
}

func (c *PackageEntry) Requires() rpm.Dependencies {
	if deps, err := c.db.DependenciesByPackage(c.key, "requires"); err != nil {
		return nil
	} else {
		return deps
	}
}

func (c *PackageEntry) Provides() rpm.Dependencies {
	if deps, err := c.db.DependenciesByPackage(c.key, "provides"); err != nil {
		return nil
	} else {
		return deps
	}
}

func (c *PackageEntry) Conflicts() rpm.Dependencies {
	if deps, err := c.db.DependenciesByPackage(c.key, "conflicts"); err != nil {
		return nil
	} else {
		return deps
	}
}

func (c *PackageEntry) Obsoletes() rpm.Dependencies {
	if deps, err := c.db.DependenciesByPackage(c.key, "obsoletes"); err != nil {
		return nil
	} else {
		return deps
	}
}

func (c *PackageEntry) Files() []string {
	if files, err := c.db.FilesByPackage(c.key); err != nil {
		return nil
	} else {
		return files
	}
}
