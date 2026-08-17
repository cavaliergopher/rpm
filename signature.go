package rpm

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/errors"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/sha3"
)

var (
	// ErrMD5CheckFailed indicates that an rpm package failed MD5 checksum
	// validation.
	ErrMD5CheckFailed = fmt.Errorf("MD5 checksum validation failed")

	// ErrGPGCheckFailed indicates that an rpm package failed GPG signature
	// validation.
	ErrGPGCheckFailed = fmt.Errorf("GPG signature validation failed")

	// ErrSHA1CheckFailed indicates that an rpm package failed SHA1 header
	// checksum validation.
	ErrSHA1CheckFailed = fmt.Errorf("SHA1 header checksum validation failed")

	// ErrSHA256CheckFailed indicates that an rpm package failed SHA256 header
	// checksum validation.
	ErrSHA256CheckFailed = fmt.Errorf("SHA256 header checksum validation failed")

	// ErrSHA3_256CheckFailed indicates that an rpm package failed SHA3-256
	// header checksum validation.
	ErrSHA3_256CheckFailed = fmt.Errorf("SHA3_256 header checksum validation failed")

	// ErrPayloadSHA256CheckFailed indicates that an rpm package failed SHA256
	// payload digest validation.
	ErrPayloadSHA256CheckFailed = fmt.Errorf("SHA256 payload digest validation failed")
)

// in order of precedence
var gpgTags = []int{
	1002, // RPMSIGTAG_PGP
	1006, // RPMSIGTAG_PGP5
	1005, // RPMSIGTAG_GPG
}

// in order of precedence, signatures covering the header region only
var headerGpgTags = []int{
	267, // RPMTAG_DSAHEADER
	268, // RPMTAG_RSAHEADER
}

// see: https://github.com/rpm-software-management/rpm/blob/3b1f4b0c6c9407b08620a5756ce422df10f6bd1a/rpmio/rpmpgp.c#L51
var gpgPubkeyTbl = map[packet.PublicKeyAlgorithm]string{
	packet.PubKeyAlgoRSA:            "RSA",
	packet.PubKeyAlgoRSASignOnly:    "RSA(Sign-Only)",
	packet.PubKeyAlgoRSAEncryptOnly: "RSA(Encrypt-Only)",
	packet.PubKeyAlgoElGamal:        "Elgamal",
	packet.PubKeyAlgoDSA:            "DSA",
	packet.PubKeyAlgoECDH:           "Elliptic Curve",
	packet.PubKeyAlgoECDSA:          "ECDSA",
}

// Map Go hashes to rpm info name
// See: https://golang.org/src/crypto/crypto.go?s=#L23
//
//	https://github.com/rpm-software-management/rpm/blob/3b1f4b0c6c9407b08620a5756ce422df10f6bd1a/rpmio/rpmpgp.c#L88
var gpgHashTbl = []string{
	"Unknown hash algorithm",
	"MD4",
	"MD5",
	"SHA1",
	"SHA224",
	"SHA256",
	"SHA384",
	"SHA512",
	"MD5SHA1",
	"RIPEMD160",
	"SHA3_224",
	"SHA3_256",
	"SHA3_384",
	"SHA3_512",
	"SHA512_224",
	"SHA512_256",
}

// GPGSignature is the raw byte representation of a package's signature.
type GPGSignature []byte

func (b GPGSignature) String() string {
	pkt, err := packet.Read(bytes.NewReader(b))
	if err != nil {
		return ""
	}
	switch sig := pkt.(type) {
	case *packet.SignatureV3:
		algo, ok := gpgPubkeyTbl[sig.PubKeyAlgo]
		if !ok {
			algo = "Unknown public key algorithm"
		}
		hasher := gpgHashTbl[0]
		if int(sig.Hash) < len(gpgHashTbl) {
			hasher = gpgHashTbl[sig.Hash]
		}
		ctime := sig.CreationTime.UTC().Format(TimeFormat)
		return fmt.Sprintf("%v/%v, %v, Key ID %x", algo, hasher, ctime, sig.IssuerKeyId)
	}
	return ""
}

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

// readHeaderRegion reads the raw header region of a rpm package from r and
// returns it. The header region covers the header data structure, the index
// entries and all string and binary data of the header.
func readHeaderRegion(r io.Reader) ([]byte, error) {
	// decode the header structure header
	var hdrBytes rpmHeader
	if _, err := io.ReadFull(r, hdrBytes[:]); err != nil {
		return nil, err
	}
	if hdrBytes.Size() > r_MaxHeaderSize {
		return nil, errorf(
			"header size exceeds the maximum of %d: %d",
			r_MaxHeaderSize,
			hdrBytes.Size(),
		)
	}
	if hdrBytes.IndexCount()*len(hdrBytes) > r_MaxHeaderSize {
		return nil, errorf(
			"header index size exceeds the maximum of %d: %d",
			r_MaxHeaderSize,
			hdrBytes.IndexCount()*len(hdrBytes),
		)
	}

	// read the raw header region bytes
	buf := make([]byte, len(hdrBytes)+hdrBytes.Size()+hdrBytes.IndexCount()*len(hdrBytes))
	copy(buf, hdrBytes[:])
	if _, err := io.ReadFull(r, buf[len(hdrBytes):]); err != nil {
		return nil, err
	}
	return buf, nil
}

// GPGCheck validates the integrity of an rpm package file. Public keys in the
// given keyring are used to validate the package signature.
//
// If validation fails, ErrGPGCheckFailed is returned.
func GPGCheck(r io.Reader, keyring openpgp.KeyRing) (string, error) {
	sig, err := readSigHeader(r)
	if err != nil {
		return "", err
	}

	// Check for v6 OPENPGP signature (tag 278) first.
	// Tag 278 contains base64-encoded OpenPGP signature(s) covering the header region.
	if tag278 := sig.GetTag(278); tag278 != nil {
		sigStrs := tag278.StringSlice()
		if len(sigStrs) > 0 {
			// Read the raw header region which the signature covers
			headerBytes, err := readHeaderRegion(r)
			if err != nil {
				return "", err
			}
			// Try each signature in the array
			for _, sigStr := range sigStrs {
				rawSig, err := base64.StdEncoding.DecodeString(sigStr)
				if err != nil {
					continue
				}
				signer, err := openpgp.CheckDetachedSignature(
					keyring,
					bytes.NewReader(headerBytes),
					bytes.NewReader(rawSig),
				)
				if err == errors.ErrUnknownIssuer {
					continue
				} else if err != nil {
					continue
				}
				for id := range signer.Identities {
					return id, nil
				}
			}
		}
		return "", ErrGPGCheckFailed
	}

	// search for a signature covering the header and the payload
	var sigval []byte
	for _, tag := range gpgTags {
		if sigval = sig.GetTag(tag).Bytes(); sigval != nil {
			break
		}
	}
	// search for a signature covering only the header
	if sigval == nil {
		for _, tag := range headerGpgTags {
			if sigval = sig.GetTag(tag).Bytes(); sigval != nil {
				break
			}
		}
		if sigval != nil {
			// read the raw header region which the signature covers
			headerBytes, err := readHeaderRegion(r)
			if err != nil {
				return "", err
			}
			r = bytes.NewReader(headerBytes)
		}
	}
	if sigval == nil {
		return "", errorf("package signature not found")
	}
	signer, err := openpgp.CheckDetachedSignature(keyring, r, bytes.NewReader(sigval))
	if err == errors.ErrUnknownIssuer {
		return "", ErrGPGCheckFailed
	} else if err != nil {
		return "", err
	}
	for id := range signer.Identities {
		return id, nil
	}
	return "", errorf("no identity found in public key")
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
	// v6 packages use SHA3_256 (tag 279) for header integrity, not MD5.
	// If tag 279 is present, skip MD5 validation.
	if sigheader.GetTag(279) != nil {
		return nil
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

// SHA1Check validates the integrity of an rpm package file. The SHA1 checksum
// is computed for the package header and compared with the checksum specified
// in the package signature.
//
// If validation fails, ErrSHA1CheckFailed is returned.
func SHA1Check(r io.Reader) error {
	sig, err := readSigHeader(r)
	if err != nil {
		return err
	}
	expected := sig.GetTag(269).String() // RPMTAG_SHA1HEADER
	if expected == "" {
		return errorf("tag not found: RPMTAG_SHA1HEADER")
	}
	headerBytes, err := readHeaderRegion(r)
	if err != nil {
		return err
	}
	h := sha1.New()
	h.Write(headerBytes)
	actual := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(actual, expected) {
		return ErrSHA1CheckFailed
	}
	return nil
}

// SHA256Check validates the integrity of an rpm package file. The SHA256
// checksum is computed for the package header and compared with the checksum
// specified in the package signature.
//
// If validation fails, ErrSHA256CheckFailed is returned.
func SHA256Check(r io.Reader) error {
	sig, err := readSigHeader(r)
	if err != nil {
		return err
	}
	expected := sig.GetTag(273).String() // RPMTAG_SHA256HEADER
	if expected == "" {
		return errorf("tag not found: RPMTAG_SHA256HEADER")
	}
	headerBytes, err := readHeaderRegion(r)
	if err != nil {
		return err
	}
	h := sha256.New()
	h.Write(headerBytes)
	actual := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(actual, expected) {
		return ErrSHA256CheckFailed
	}
	return nil
}

// SHA3_256Check validates the integrity of an rpm package file. The SHA3-256
// checksum is computed for the package header and compared with the checksum
// specified in the package signature.
//
// If validation fails, ErrSHA3_256CheckFailed is returned.
func SHA3_256Check(r io.Reader) error {
	sig, err := readSigHeader(r)
	if err != nil {
		return err
	}
	expected := sig.GetTag(279).String() // RPMTAG_SHA3_256HEADER
	if expected == "" {
		return errorf("tag not found: RPMTAG_SHA3_256HEADER")
	}
	headerBytes, err := readHeaderRegion(r)
	if err != nil {
		return err
	}
	h := sha3.New256()
	h.Write(headerBytes)
	actual := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(actual, expected) {
		return ErrSHA3_256CheckFailed
	}
	return nil
}

// PayloadSHA256Check validates the integrity of an rpm package file. The
// SHA256 digest of the compressed payload is computed and compared with the
// digest specified in the package main header.
//
// If validation fails, ErrPayloadSHA256CheckFailed is returned.
func PayloadSHA256Check(r io.Reader) error {
	if _, err := readSigHeader(r); err != nil {
		return err
	}
	hdr, err := readHeader(r, false)
	if err != nil {
		return err
	}
	expected := hdr.GetTag(5092).String() // RPMTAG_PAYLOADSHA256
	if expected == "" {
		if hdr.GetTag(5097) != nil {
			return nil // ALT-only: tolerated, requires decompression (out of scope)
		}
		return errorf("tag not found: RPMTAG_PAYLOADSHA256")
	}
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return err
	}
	actual := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(actual, expected) {
		return ErrPayloadSHA256CheckFailed
	}
	return nil
}

// ReadKeyRing reads a openpgp.KeyRing from the given io.Reader which may then
// be used to validate GPG keys in rpm packages.
func ReadKeyRing(r io.Reader) (openpgp.KeyRing, error) {
	// decode gpgkey file
	p, err := armor.Decode(r)
	if err != nil {
		return nil, err
	}

	// extract keys
	return openpgp.ReadKeyRing(p.Body)
}

func openKeyRing(name string) (openpgp.KeyRing, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadKeyRing(f)
}

// KeyRingFromFiles reads a openpgp.KeyRing from the given file paths which may
// then be used to validate GPG keys in rpm packages.
//
// This function might typically be used to read all keys in /etc/pki/rpm-gpg.
func OpenKeyRing(name ...string) (openpgp.KeyRing, error) {
	entityList := make(openpgp.EntityList, 0)
	for _, path := range name {
		keyring, err := openKeyRing(path)
		if err != nil {
			return nil, err
		}
		entityList = append(entityList, keyring.(openpgp.EntityList)...)
	}
	return entityList, nil
}
