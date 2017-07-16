package rpm

import (
	"fmt"
	"testing"
)

type DepTest struct {
	dep Dependency
	str string
}

func TestDependencyStrings(t *testing.T) {
	tests := []DepTest{
		{&dependency{DepFlagAny, "test", 0, "", ""}, "test"},
		{&dependency{DepFlagAny, "test", 0, "1", ""}, "test 1"},
		{&dependency{DepFlagAny, "test", 0, "1", "2"}, "test 1.2"},
		{&dependency{DepFlagAny, "test", 1, "2", "3"}, "test 2.3"},
		{&dependency{DepFlagLesser, "test", 0, "1", ""}, "test < 1"},
		{&dependency{DepFlagLesser, "test", 0, "1", "2"}, "test < 1.2"},
		{&dependency{DepFlagLesser, "test", 1, "2", "3"}, "test < 2.3"},
		{&dependency{DepFlagLesserOrEqual, "test", 0, "1", ""}, "test <= 1"},
		{&dependency{DepFlagLesserOrEqual, "test", 0, "1", "2"}, "test <= 1.2"},
		{&dependency{DepFlagLesserOrEqual, "test", 1, "2", "3"}, "test <= 2.3"},
		{&dependency{DepFlagGreaterOrEqual, "test", 0, "1", ""}, "test >= 1"},
		{&dependency{DepFlagGreaterOrEqual, "test", 0, "1", "2"}, "test >= 1.2"},
		{&dependency{DepFlagGreaterOrEqual, "test", 1, "2", "3"}, "test >= 2.3"},
		{&dependency{DepFlagLesser, "test", 0, "1", ""}, "test < 1"},
		{&dependency{DepFlagLesser, "test", 0, "1", "2"}, "test < 1.2"},
		{&dependency{DepFlagLesser, "test", 1, "2", "3"}, "test < 2.3"},
	}

	for i, test := range tests {
		if str := fmt.Sprintf("%v", test.dep); str != test.str {
			t.Errorf("Expected '%s' for test %d, got: '%s'", test.str, i+1, str)
		}
	}
}
