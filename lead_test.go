package rpm

import (
	"bytes"
	"io"
	"testing"

	"github.com/cavaliercoder/badio"
)

func TestLeadErrors(t *testing.T) {
	// load package file paths
	files := getTestFiles()

	// load each package
	for _, fb := range files {
		// local copy the file so the following corruptions don't break
		// other tests
		b := make([]byte, len(fb))
		copy(b, fb)

		// simulate read error
		f := bytes.NewReader(b)
		r := badio.NewBreakReader(f, 95)
		_, err := ReadPackageLead(r)
		if err == nil {
			t.Errorf("Expected read error in lead section, got: %v", err)
		}

		// simulate length error
		f.Seek(0, 0) // io.SeekStart since ~1.7
		r = badio.NewTruncateReader(f, 95)
		_, err = ReadPackageLead(r)
		if err != io.ErrUnexpectedEOF {
			t.Errorf("Expected bad length error in lead section, got: %v", err)
		}

		// simulate version error
		f.Seek(0, 0) // io.SeekStart since ~1.7
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
	_, err := ReadPackageLead(r)
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
