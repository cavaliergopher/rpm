package rpm

import (
	"fmt"
	"strings"
)

// TimeFormat is the time format used by the rpm ecosystem. The time being
// formatted must be in UTC for Format to generate the correct format.
const TimeFormat = "Mon Jan _2 15:04:05 2006"

func errorf(format string, a ...interface{}) error {
	return fmt.Errorf("rpm: "+format, a...)
}

func parseVersion(v string) (epoch int, version, release string) {
	if i := strings.IndexByte(v, ':'); i >= 0 {
		epoch, v = parseInt(v[:i]), v[i+1:]
	}

	if i := strings.IndexByte(v, '-'); i >= 0 {
		return epoch, v[:i], v[i+1:]
	}

	return epoch, v, ""
}

func parseInt(s string) int {
	var n int
	for _, dec := range s {
		if dec < '0' || dec > '9' {
			return 0
		}
		n = n*10 + (int(dec) - '0')
	}
	return n
}
