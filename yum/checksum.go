package yum

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// ErrChecksumMismatch indicates that the checksum value of two items does not
// match.
var ErrChecksumMismatch = fmt.Errorf("Checksum mismatch")

type RepoDatabaseChecksum struct {
	Type string `xml:"type,attr"`
	Hash string `xml:",chardata"`
}

// Check creates a checksum of the given io.Reader content and compares it the
// the expected checksum value. If the checksums match, nil is returned. If the
// checksums do not match, ErrChecksumMismatch is returned. If any other error
// occurs, the error is returned.
func (c *RepoDatabaseChecksum) Check(r io.Reader) error {
	// get checksum value based by type
	actual := ""
	switch c.Type {
	case "sha256":
		s := sha256.New()
		if _, err := io.Copy(s, r); err != nil {
			return err
		}

		actual = hex.EncodeToString(s.Sum(nil))

	default:
		return fmt.Errorf("Unsupported checksum type: %s", c.Type)
	}

	// check against expected value
	if c.Hash != actual {
		return ErrChecksumMismatch
	}

	return nil
}
