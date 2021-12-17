package rpm

import (
	"encoding/json"
	"os"
	"testing"
)

type VerTest struct {
	A      string `json:"a"`
	B      string `json:"b"`
	Expect int    `json:"expect"`
}

type TestPkg struct {
	E int
	V string
	R string
}

func (c *TestPkg) Name() string    { return "test" }
func (c *TestPkg) Epoch() int      { return c.E }
func (c *TestPkg) Version() string { return c.V }
func (c *TestPkg) Release() string { return c.R }

func sign(r int) string {
	if r < 0 {
		return "<"
	} else if r > 0 {
		return ">"
	}
	return "=="
}

func TestCompare(t *testing.T) {
	tests := make([]*VerTest, 0)
	f, err := os.Open("./testdata/vercmp.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&tests); err != nil {
		t.Fatal(err)
	}
	if len(tests) == 0 {
		t.Fatal("version tests are empty")
	}
	for _, test := range tests {
		// compare 'version'
		a := &TestPkg{0, test.A, ""}
		b := &TestPkg{0, test.B, ""}
		if r := Compare(a, b); r != test.Expect {
			t.Errorf(
				"Expected %s %s %s; got %s %s %s",
				test.A,
				sign(test.Expect),
				test.B,
				test.A,
				sign(r),
				test.B,
			)
		}

		// compare 'release'
		a = &TestPkg{0, "", test.A}
		b = &TestPkg{0, "", test.B}
		if r := Compare(a, b); r != test.Expect {
			t.Errorf(
				"Expected %s %s %s; got %s %s %s",
				test.A,
				sign(test.Expect),
				test.B,
				test.A,
				sign(r),
				test.B,
			)
		}
	}
	if r := Compare(nil, nil); r != 0 {
		t.Errorf("Expected <nil> == <nil>; got <nil> %s <nil>", sign(r))
	}
	if r := Compare(nil, &TestPkg{1, "", ""}); r != -1 {
		t.Errorf("Expected <nil> < <nil>; got <nil> %s <nil>", sign(r))
	}
	if r := Compare(&TestPkg{1, "", ""}, nil); r != 1 {
		t.Errorf("Expected <nil> > <nil>; got <nil> %s <nil>", sign(r))
	}
}
