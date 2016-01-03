package rpm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// A Lead is the deprecated lead section of an RPM file which is used in legacy
// RPM versions to store package metadata.
type Lead struct {
	VersionMajor    int
	VersionMinor    int
	Name            string
	Type            int
	Architecture    int
	OperatingSystem int
	SignatureType   int
}

// ReadPackageLead reads the deprecated lead section of an RPM file which is
// used in legacy RPM versions to store package metadata.
//
// This function should only be used if you intend to read a package lead in
// isolation.
func ReadPackageLead(r io.Reader) (*Lead, error) {
	// read bytes
	b := make([]byte, 96)
	n, err := r.Read(b)
	if err != nil {
		return nil, fmt.Errorf("Error reading RPM Lead section: %s", err)
	}

	// check length
	if n != 96 {
		return nil, fmt.Errorf("RPM Lead section is incorrect length")
	}

	// check magic number
	if 0 != bytes.Compare(b[:4], []byte{0xED, 0xAB, 0xEE, 0xDB}) {
		return nil, fmt.Errorf("RPM file descriptor is invalid")
	}

	// decode lead
	lead := &Lead{
		VersionMajor:    int(b[5]),
		VersionMinor:    int(b[6]),
		Type:            int(binary.BigEndian.Uint16(b[7:9])),
		Architecture:    int(binary.BigEndian.Uint16(b[9:11])),
		Name:            string(b[10:77]),
		OperatingSystem: int(binary.BigEndian.Uint16(b[76:78])),
		SignatureType:   int(binary.BigEndian.Uint16(b[78:80])),
	}

	// TODO: validate lead value ranges

	return lead, nil
}
