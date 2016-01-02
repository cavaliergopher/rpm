package rpm

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

	Requires() Dependencies
	Conflicts() Dependencies
	Obsoletes() Dependencies
	Provides() Dependencies
}

// Packages is a slice of Package interfaces.
type Packages []Package
