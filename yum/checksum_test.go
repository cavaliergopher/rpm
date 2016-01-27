package yum

import (
	"bytes"
	"testing"
)

type ChecksumTest struct {
	Checksum     string
	ChecksumType string
	Value        []byte
}

func TestValidateChecksum(t *testing.T) {
	tests := []ChecksumTest{
		ChecksumTest{"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", "sha256", []byte{}},
		ChecksumTest{"054edec1d0211f624fed0cbca9d4f9400b0e491c43742af2c5b0abebf0c990d8", "sha256", []byte{0x00, 0x01, 0x02, 0x03}},
		ChecksumTest{"1e584b5a9a8387cadf4449efa6a632fd31b307d9d5cdf6cf70ac2bf9d1cb9513", "sha256", []byte{0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA}},
	}

	for i, test := range tests {
		if err := ValidateChecksum(bytes.NewReader(test.Value), test.Checksum, test.ChecksumType); err == ErrChecksumMismatch {
			t.Errorf("Checksum validation failed for test %d", i+1)
		}
	}

	t.Logf("%d checksums validated", len(tests))
}
