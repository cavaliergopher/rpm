package rpm

import "fmt"

// TimeFormat is the time format used by the rpm ecosystem. The time being
// formatted must be in UTC for Format to generate the correct format.
const TimeFormat = "Mon Jan _2 15:04:05 2006"

func errorf(format string, a ...interface{}) error {
	return fmt.Errorf("rpm: "+format, a...)
}
