package rpm

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
)

var (
	// ErrMD5CheckFailed indicates that an rpm package failed MD5 checksum
	// validation.
	ErrMD5CheckFailed = fmt.Errorf("MD5 checksum validation failed")
)

// GPGSignature is the raw byte representation of a package's signature.
//
// Deprecated: OpenPGP signature verification was removed because its former
// dependency is unmaintained and unsafe. Use RPM-native verification tooling
// when package authenticity must be established.
type GPGSignature []byte

// readSigHeader reads the lead and signature header of a rpm package and stops
// the reader at the beginning of the header header.
func readSigHeader(r io.Reader) (*Header, error) {
	lead, err := readLead(r)
	if err != nil {
		return nil, err
	}
	if lead.SignatureType != 5 { // RPMSIGTYPE_HEADERSIG
		return nil, errorf("unknown signature type: %x", lead.SignatureType)
	}
	sig, err := readHeader(r, true)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

// MD5Check validates the integrity of an rpm package file. The MD5 checksum is
// computed for the package payload and compared with the checksum specified in
// the package header.
//
// If validation fails, ErrMD5CheckFailed is returned.
func MD5Check(r io.Reader) error {
	sigheader, err := readSigHeader(r)
	if err != nil {
		return err
	}
	payloadSize := sigheader.GetTag(270).Int64() // RPMSIGTAG_LONGSIGSIZE
	if payloadSize == 0 {
		payloadSize = sigheader.GetTag(1000).Int64() // RPMSIGTAG_SIGSIZE
		if payloadSize == 0 {
			return fmt.Errorf("tag not found: RPMSIGTAG_SIZE")
		}
	}
	expect := sigheader.GetTag(1004).Bytes() // RPMSIGTAG_MD5
	if expect == nil {
		return errorf("tag not found: RPMSIGTAG_MD5")
	}
	h := md5.New()
	if n, err := io.Copy(h, r); err != nil {
		return err
	} else if n != payloadSize {
		return ErrMD5CheckFailed
	}
	actual := h.Sum(nil)
	if !bytes.Equal(expect, actual) {
		return ErrMD5CheckFailed
	}
	return nil
}
