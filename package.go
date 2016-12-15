package rpm

// PackageVersion is an interface which holds version information for a single
// package version.
type PackageVersion interface {
	Name() string
	Version() string
	Release() string
	Epoch() int
}
