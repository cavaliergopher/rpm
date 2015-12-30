package yum

import (
	"fmt"
)

// PackageEntry is a RPM package as defined in a package repository database.
type PackageEntry struct {
	row []interface{}
}

type PackageEntries []PackageEntry

func NewPackageEntry(row []interface{}) (*PackageEntry, error) {
	return &PackageEntry{
		row: row,
	}, nil
}

// String reassembles package metadata to form a standard rpm package name;
// including the package name, version, release and architecture.
func (c PackageEntry) String() string {
	return fmt.Sprintf("%s-%s-%s.%s", c.Name(), c.Version(), c.Release(), c.Architecture())
}

// LocationHref is the location of the package, relative to the parent
// repository.
func (c *PackageEntry) LocationHref() string {
	return string(c.row[23].([]byte))
}

func (c *PackageEntry) Checksum() string {
	return string(c.row[1].([]byte))
}

func (c *PackageEntry) ChecksumType() string {
	return string(c.row[25].([]byte))
}

func (c *PackageEntry) Name() string {
	return string(c.row[2].([]byte))
}

func (c *PackageEntry) Version() string {
	return string(c.row[4].([]byte))
}

func (c *PackageEntry) Release() string {
	return string(c.row[6].([]byte))
}

func (c *PackageEntry) Architecture() string {
	return string(c.row[3].([]byte))
}
