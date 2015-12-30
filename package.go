package rpm

// Package is an interface which represents an RPM package from any source.
type Package interface {
	Name() string
	Version() string
	Release() string
	Epoch() int64
	Architecture() string
}
