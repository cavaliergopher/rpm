package rpm

import (
	"bytes"
	"encoding/binary"
	"io"
)

// ErrNotRPMFile indicates that the file is not an rpm package.
var ErrNotRPMFile = errorf("invalid file descriptor")

// Lead is the deprecated lead section of an rpm file which is used in legacy
// rpm versions to store package metadata.
type Lead struct {
	VersionMajor    int
	VersionMinor    int
	Name            string
	Type            int
	Architecture    int
	OperatingSystem int
	SignatureType   int
}

type leadBytes [96]byte

func (c leadBytes) Magic() []byte        { return c[:4] }
func (c leadBytes) VersionMajor() int    { return int(c[4]) }
func (c leadBytes) VersionMinor() int    { return int(c[5]) }
func (c leadBytes) Type() int            { return int(binary.BigEndian.Uint16(c[6:8])) }
func (c leadBytes) Architecture() int    { return int(binary.BigEndian.Uint16(c[8:10])) }
func (c leadBytes) Name() string         { return string(c[10:76]) }
func (c leadBytes) OperatingSystem() int { return int(binary.BigEndian.Uint16(c[76:78])) }
func (c leadBytes) SignatureType() int   { return int(binary.BigEndian.Uint16(c[78:80])) }

// readLead reads the deprecated lead section of an rpm package which is used in
// legacy rpm versions to store package metadata.
func readLead(r io.Reader) (*Lead, error) {
	var lead leadBytes
	_, err := io.ReadFull(r, lead[:])
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(lead.Magic(), []byte{0xED, 0xAB, 0xEE, 0xDB}) {
		return nil, ErrNotRPMFile
	}
	if lead.VersionMajor() < 3 || lead.VersionMajor() > 4 {
		return nil, errorf("unsupported rpm version: %d", lead.VersionMajor())
	}
	// TODO: validate lead value ranges
	return &Lead{
		VersionMajor:    lead.VersionMajor(),
		VersionMinor:    lead.VersionMinor(),
		Type:            lead.Type(),
		Architecture:    lead.Architecture(),
		Name:            lead.Name(),
		OperatingSystem: lead.OperatingSystem(),
		SignatureType:   lead.SignatureType(),
	}, nil
}
