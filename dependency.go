package rpm

import (
	"fmt"
)

// RPM dependency flags.
// See: https://github.com/rpm-software-management/rpm/blob/master/lib/rpmds.h#L25
const (
	DepFlagAny            = 0
	DepFlagLess           = (1 << 1)
	DepFlagGreater        = (1 << 2)
	DepFlagEqual          = (1 << 3)
	DepFlagLesserOrEqual  = (DepFlagEqual | DepFlagLess)
	DepFlagGreaterOrEqual = (DepFlagEqual | DepFlagGreater)
)

// Dependency is an interface which represents a relationship between two
// packages. It might indicate that one package requires, conflicts with,
// obsoletes or provides another package.
type Dependency interface {
	PackageVersion

	// One of EQ, LT, LE, GE or GT
	Flags() int64
}

// private basic implementation or a package dependency.
type dependency struct {
	flags   int64
	name    string
	epoch   int64
	version string
	release string
}

// Dependencies are a slice of Dependency interfaces.
type Dependencies []Dependency

func NewDependency(flgs int64, name string, epoch int64, version string, release string) Dependency {
	return &dependency{
		flags:   flgs,
		name:    name,
		epoch:   epoch,
		version: version,
		release: release,
	}
}

func (c *dependency) String() string {
	s := c.name

	switch {
	case DepFlagLesserOrEqual == (c.flags & DepFlagLesserOrEqual):
		s = fmt.Sprintf("%s <=", s)

	case DepFlagGreaterOrEqual == (c.flags & DepFlagGreaterOrEqual):
		s = fmt.Sprintf("%s >=", s)

	case DepFlagEqual == (c.flags & DepFlagEqual):
		s = fmt.Sprintf("%s =", s)
	}

	if c.version != "" {
		s = fmt.Sprintf("%s %s", s, c.version)
	}

	if c.release != "" {
		s = fmt.Sprintf("%s.%s", s, c.release)
	}

	return s
}

func (c *dependency) Flags() int64 {
	return c.flags
}

func (c *dependency) Name() string {
	return c.name
}

func (c *dependency) Epoch() int64 {
	return c.epoch
}

func (c *dependency) Version() string {
	return c.version
}

func (c *dependency) Release() string {
	return c.release
}
