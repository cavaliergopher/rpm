package rpm

import (
	"time"
)

// PackageVersion is an interface which holds version information for a single
// package version.
type PackageVersion interface {
	Name() string
	Version() string
	Release() string
	Epoch() int
}

// Package is an interface which represents an RPM package and its supported
// tags.
type Package interface {
	PackageVersion

	Architecture() string

	Path() string
	FileTime() time.Time
	FileSize() uint64

	Size() uint64
	ArchiveSize() uint64

	BuildTime() time.Time

	HeaderStart() uint64
	HeaderEnd() uint64

	Summary() string
	Description() string
	URL() string

	License() string
	Vendor() string
	Groups() []string
	BuildHost() string
	SourceRPM() string
	Packager() string

	Requires() Dependencies
	Conflicts() Dependencies
	Obsoletes() Dependencies
	Provides() Dependencies

	Checksum() (string, error)
	ChecksumType() string
}

// Packages is a slice of Package interfaces.
type Packages []Package
