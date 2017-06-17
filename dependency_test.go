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
		{NewDependency(DepFlagAny, "test", 0, "", ""), "test"},
		{NewDependency(DepFlagAny, "test", 0, "1", ""), "test 1"},
		{NewDependency(DepFlagAny, "test", 0, "1", "2"), "test 1.2"},
		{NewDependency(DepFlagAny, "test", 1, "2", "3"), "test 2.3"},
		{NewDependency(DepFlagLesser, "test", 0, "1", ""), "test < 1"},
		{NewDependency(DepFlagLesser, "test", 0, "1", "2"), "test < 1.2"},
		{NewDependency(DepFlagLesser, "test", 1, "2", "3"), "test < 2.3"},
		{NewDependency(DepFlagLesserOrEqual, "test", 0, "1", ""), "test <= 1"},
		{NewDependency(DepFlagLesserOrEqual, "test", 0, "1", "2"), "test <= 1.2"},
		{NewDependency(DepFlagLesserOrEqual, "test", 1, "2", "3"), "test <= 2.3"},
		{NewDependency(DepFlagGreaterOrEqual, "test", 0, "1", ""), "test >= 1"},
		{NewDependency(DepFlagGreaterOrEqual, "test", 0, "1", "2"), "test >= 1.2"},
		{NewDependency(DepFlagGreaterOrEqual, "test", 1, "2", "3"), "test >= 2.3"},
		{NewDependency(DepFlagLesser, "test", 0, "1", ""), "test < 1"},
		{NewDependency(DepFlagLesser, "test", 0, "1", "2"), "test < 1.2"},
		{NewDependency(DepFlagLesser, "test", 1, "2", "3"), "test < 2.3"},
	}

	for i, test := range tests {
		if str := fmt.Sprintf("%v", test.dep); str != test.str {
			t.Errorf("Expected '%s' for test %d, got: '%s'", test.str, i+1, str)
		}
	}
}
