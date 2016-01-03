package rpm

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
)

var (
	// ErrMD5ValidationFailed indicates that a RPM package failed checksum
	// validation.
	ErrMD5ValidationFailed = fmt.Errorf("Package checksum validation failed")
)

// MD5Check validates the integrity of a RPM package file read from the given
// io.Reader. An MD5 checksum is computed for the package payload and compared
// with the checksum value specified in the package header.
//
// If validation succeeds, nil is returned. If validation fails,
// ErrMD5ValidationFailed is returned.
//
// This function is an expensive operation which reads the entire package file.
func MD5Check(r io.Reader) error {
	_, _, err := md5check(r)
	return err
}

// md5check reads a RPM package file from the given io.Reader and returns the
// signature header and computed md5 sum for the package payload.
//
// If the package fails checksum validation, a GPGCheckError is returned.
func md5check(r io.Reader) (*Header, []byte, error) {
	// read package lead
	if _, err := ReadPackageLead(r); err != nil {
		return nil, nil, err
	}

	// read signature header
	sigheader, err := ReadPackageHeader(r)
	if err != nil {
		return nil, nil, err
	}

	// get expected payload size
	payloadSize := sigheader.Indexes.IntByTag(1000)
	if payloadSize == 0 {
		return nil, nil, fmt.Errorf("No payload size specified")
	}

	// get expected payload md5 sum
	sigmd5 := sigheader.Indexes.BytesByTag(1004)
	if sigmd5 == nil {
		return nil, nil, fmt.Errorf("No payload md5 sum specified")
	}

	// compute payload sum
	h := md5.New()
	if n, err := io.Copy(h, r); err != nil {
		return nil, nil, fmt.Errorf("Error reading payload: %v", err)
	} else if n != payloadSize {
		return nil, nil, ErrMD5ValidationFailed
	}

	// compare sums
	payloadmd5 := h.Sum(nil)
	if !bytes.Equal(payloadmd5, sigmd5) {
		return nil, nil, ErrMD5ValidationFailed
	}

	return sigheader, h.Sum(nil), nil
}
