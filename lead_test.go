package rpm

import (
	"bytes"
	"github.com/cavaliercoder/badio"
	"io/ioutil"
	"os"
	"testing"
)

func TestLeadErrors(t *testing.T) {
	// load package file paths
	files, err := packages(t)
	if err != nil {
		t.Fatalf("Error listing rpm packages: %v", err)
	}

	// load each package
	for _, path := range files {
		// simulate read error
		f, _ := os.Open(path)
		r := badio.NewBreakReader(f, 95)
		_, err := ReadPackageLead(r)
		if err == nil {
			t.Errorf("Expected read error in lead section, got: %v", err)
		}
		f.Close()

		// simulate length error
		f, _ = os.Open(path)
		r = badio.NewTruncateReader(f, 95)
		_, err = ReadPackageLead(r)
		if err != ErrBadLeadLength {
			t.Errorf("Expected bad length error in lead section, got: %v", err)
		}
		f.Close()

		// simulate version error
		f, _ = os.Open(path)
		b, _ := ioutil.ReadAll(f)
		f.Close()

		b[4] = 0x02
		_, err = ReadPackageLead(bytes.NewReader(b))
		if err != ErrUnsupportedVersion {
			t.Errorf("Expected version error in lead section, got: %v", err)
		}

		b[4] = 0x05
		_, err = ReadPackageLead(bytes.NewReader(b))
		if err != ErrUnsupportedVersion {
			t.Errorf("Expected version error in lead section, got: %v", err)
		}
	}

	// simulate magic number error
	b := make([]byte, 96)
	r := bytes.NewReader(b)
	_, err = ReadPackageLead(r)
	if err != ErrNotRPMFile {
		t.Errorf("Expected bad descriptor error in lead section, got: %v", err)
	}

	// test handler in ReadPackageFile
	r = bytes.NewReader(b)
	_, err = ReadPackageFile(r)
	if err != ErrNotRPMFile {
		t.Errorf("Expected lead section error in bad package, got: %v", err)
	}
}
