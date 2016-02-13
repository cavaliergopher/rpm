package rpm

import (
	"fmt"
	"testing"
)

type DepTest struct {
	dep Dependency
	str string
}

func TestDependencies(t *testing.T) {
	tests := []DepTest{
		DepTest{NewDependency(DepFlagAny, "test", 0, "", ""), "test"},
		DepTest{NewDependency(DepFlagAny, "test", 0, "1", ""), "test 1"},
		DepTest{NewDependency(DepFlagAny, "test", 0, "1", "2"), "test 1.2"},
		DepTest{NewDependency(DepFlagAny, "test", 1, "2", "3"), "test 2.3"},
		DepTest{NewDependency(DepFlagLesser, "test", 0, "1", ""), "test < 1"},
		DepTest{NewDependency(DepFlagLesser, "test", 0, "1", "2"), "test < 1.2"},
		DepTest{NewDependency(DepFlagLesser, "test", 1, "2", "3"), "test < 2.3"},
		DepTest{NewDependency(DepFlagLesserOrEqual, "test", 0, "1", ""), "test <= 1"},
		DepTest{NewDependency(DepFlagLesserOrEqual, "test", 0, "1", "2"), "test <= 1.2"},
		DepTest{NewDependency(DepFlagLesserOrEqual, "test", 1, "2", "3"), "test <= 2.3"},
		DepTest{NewDependency(DepFlagGreaterOrEqual, "test", 0, "1", ""), "test >= 1"},
		DepTest{NewDependency(DepFlagGreaterOrEqual, "test", 0, "1", "2"), "test >= 1.2"},
		DepTest{NewDependency(DepFlagGreaterOrEqual, "test", 1, "2", "3"), "test >= 2.3"},
		DepTest{NewDependency(DepFlagLesser, "test", 0, "1", ""), "test < 1"},
		DepTest{NewDependency(DepFlagLesser, "test", 0, "1", "2"), "test < 1.2"},
		DepTest{NewDependency(DepFlagLesser, "test", 1, "2", "3"), "test < 2.3"},
	}

	for i, test := range tests {
		if str := fmt.Sprintf("%v", test.dep); str != test.str {
			t.Errorf("Expected '%s' for test %d, got: '%s'", test.str, i+1, str)
		}
	}
}
