package rpm

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/errors"
	"golang.org/x/crypto/openpgp/packet"
)

var (
	// ErrMD5CheckFailed indicates that an rpm package failed MD5 checksum
	// validation.
	ErrMD5CheckFailed = fmt.Errorf("MD5 checksum validation failed")

	// ErrGPGCheckFailed indicates that an rpm package failed GPG signature
	// validation.
	ErrGPGCheckFailed = fmt.Errorf("GPG signature validation failed")
)

// in order of precedence
var gpgTags = []int{
	1002, // RPMSIGTAG_PGP
	1006, // RPMSIGTAG_PGP5
	1005, // RPMSIGTAG_GPG
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
//      https://github.com/rpm-software-management/rpm/blob/3b1f4b0c6c9407b08620a5756ce422df10f6bd1a/rpmio/rpmpgp.c#L88
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

// GPGCheck validates the integrity of an rpm package file. Public keys in the
// given keyring are used to validate the package signature.
//
// If validation fails, ErrGPGCheckFailed is returned.
func GPGCheck(r io.Reader, keyring openpgp.KeyRing) (string, error) {
	sig, err := readSigHeader(r)
	if err != nil {
		return "", err
	}
	var sigval []byte
	for _, tag := range gpgTags {
		if sigval = sig.GetTag(tag).Bytes(); sigval != nil {
			break
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
	payloadSize := sigheader.GetTag(1000).Int64() // RPMSIGTAG_SIZE
	if payloadSize == 0 {
		return errorf("tag not found: RPMSIGTAG_SIZE")
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
