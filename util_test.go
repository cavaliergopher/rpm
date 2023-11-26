package rpm

import "testing"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		in      string
		epoch   int
		version string
		release string
	}{
		{"", 0, "", ""},
		{"1.0", 0, "1.0", ""},
		{"1:1.0", 1, "1.0", ""},
		{"1:1.0-test", 1, "1.0", "test"},
		{"1.0-test", 0, "1.0", "test"},
		{":1.0-", 0, "1.0", ""}, // Ensure malformed version doesn't panic.
	}

	for _, test := range tests {
		epoch, ver, rel := parseVersion(test.in)
		if epoch != test.epoch {
			t.Errorf("Expected epoch %d for %q; got %d", test.epoch, test.in, epoch)
		}
		if ver != test.version {
			t.Errorf("Expected version %s for %q; got %s", test.version, test.in, ver)
		}
		if rel != test.release {
			t.Errorf("Expected release %s for %q; got %s", test.release, test.in, rel)
		}
	}
}
