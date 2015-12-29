package rpm

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"time"
)

// A Package is an RPM package.
type Package struct {
	Lead    Lead
	Headers Headers
}

// A Lead is the deprecated lead section of an RPM file which is used in legacy
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

// A Header stores metadata about a rpm package.
type Header struct {
	Version    int
	IndexCount int
	Length     int64
	Indexes    IndexEntries
}

// Headers is an array of Header structs.
type Headers []Header

// OpenPackage reads a rpm package from the file systems and returns a pointer
// to it.
func OpenPackage(path string) (*Package, error) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error opening RPM file: %s", err)
	}
	defer f.Close()

	return ReadPackage(f)
}

// ReadPackage reads a rpm package from a stream and returns a pointer to it.
func ReadPackage(r io.Reader) (*Package, error) {
	p := &Package{}

	// read the deprecated "lead"
	lead := make([]byte, 96)
	n, err := r.Read(lead)
	if err != nil {
		return nil, fmt.Errorf("Error reading RPM Lead section: %s", err)
	}

	if n != 96 {
		return nil, fmt.Errorf("RPM Lead section is incorrect length")
	}

	// check magic number
	if 0 != bytes.Compare(lead[:4], []byte{0xED, 0xAB, 0xEE, 0xDB}) {
		return nil, fmt.Errorf("RPM file descriptor is invalid")
	}

	// translate lead
	p.Lead.VersionMajor = int(lead[5])
	p.Lead.VersionMinor = int(lead[6])
	p.Lead.Type = (int(lead[7]) << 8) + int(lead[8])
	p.Lead.Architecture = (int(lead[9]) << 8) + int(lead[10])
	p.Lead.Name = string(lead[10:77])
	p.Lead.OperatingSystem = (int(lead[76]) << 8) + int(lead[77])
	p.Lead.SignatureType = (int(lead[78]) << 8) + int(lead[79])

	// TODO: validate lead value ranges

	// parse headers
	p.Headers = make(Headers, 0)

	// TODO: find last header without using hard limit of 2
	for i := 1; i < 3; i++ {
		// read the "header structure header"
		header := make([]byte, 16)
		n, err = r.Read(header)
		if err != nil {
			return nil, fmt.Errorf("Error reading RPM structure header for header %d: %v", i, err)
		}

		if n != 16 {
			return nil, fmt.Errorf("Error reading RPM structure header for header %d: only %d bytes returned", i, n)
		}

		// check magic number
		if 0 != bytes.Compare(header[:3], []byte{0x8E, 0xAD, 0xE8}) {
			return nil, fmt.Errorf("RPM header %d is invalid", i)
		}

		// translate header
		h := Header{}
		h.Version = int(header[3])
		h.IndexCount = (int(header[8]) << 24) + (int(header[9]) << 16) + (int(header[10]) << 8) + int(header[11])
		h.Length = (int64(header[12]) << 24) + (int64(header[13]) << 16) + (int64(header[14]) << 8) + int64(header[15])
		h.Indexes = make(IndexEntries, h.IndexCount)

		// read indexes
		indexLength := 16 * h.IndexCount
		indexes := make([]byte, indexLength)
		n, err = r.Read(indexes)
		if err != nil {
			return nil, fmt.Errorf("Error reading index entries for header %d: %v", i, err)
		}

		if n != indexLength {
			return nil, fmt.Errorf("Error reading index entries for header %d: only %d bytes returned", i, n)
		}

		for x := 0; x < h.IndexCount; x++ {
			o := 16 * x
			index := IndexEntry{}

			index.Tag = (int64(indexes[o]) << 24) + (int64(indexes[o+1]) << 16) + (int64(indexes[o+2]) << 8) + int64(indexes[o+3])
			index.Type = (int64(indexes[o+4]) << 24) + (int64(indexes[o+5]) << 16) + (int64(indexes[o+6]) << 8) + int64(indexes[o+7])
			index.Offset = (int64(indexes[o+8]) << 24) + (int64(indexes[o+9]) << 16) + (int64(indexes[o+10]) << 8) + int64(indexes[o+11])
			index.ItemCount = (int64(indexes[o+12]) << 24) + (int64(indexes[o+13]) << 16) + (int64(indexes[o+14]) << 8) + int64(indexes[o+15])
			h.Indexes[x] = index
		}

		// read the "store"
		store := make([]byte, h.Length)
		n, err = r.Read(store)
		if err != nil {
			return nil, fmt.Errorf("Error reading store for header %d: %v", i, err)
		}

		if int64(n) != h.Length {
			return nil, fmt.Errorf("Error reading store for header %d: only %d bytes returned", i, n)
		}

		for x := 0; x < h.IndexCount; x++ {
			index := h.Indexes[x]

			switch index.Type {
			case IndexDataTypeChar:
				index.Value = uint8(store[index.Offset])
				break

			case IndexDataTypeInt8:
				index.Value = int8(store[index.Offset])
				break

			case IndexDataTypeInt16:
				index.Value = (int16(store[index.Offset]) << 8) + int16(store[index.Offset+1])
				break

			case IndexDataTypeInt32:
				index.Value = (int32(store[index.Offset]) << 24) + (int32(store[index.Offset+1]) << 16) + (int32(store[index.Offset+2]) << 8) + int32(store[index.Offset+3])
				break

			case IndexDataTypeInt64:
				index.Value = (int64(store[index.Offset]) << 56) + (int64(store[index.Offset+1]) << 48) + (int64(store[index.Offset+2]) << 40) + (int64(store[index.Offset+3]) << 32) + (int64(store[index.Offset+4]) << 24) + (int64(store[index.Offset+5]) << 16) + (int64(store[index.Offset+6]) << 8) + int64(store[index.Offset+7])
				break

			case IndexDataTypeBinary:
				b := make([]byte, index.ItemCount)
				copy(b, store[index.Offset:index.Offset+index.ItemCount])

				index.Value = b

				break

			case IndexDataTypeString, IndexDataTypeStringArray, IndexDataTypeI8NString:
				vals := make([]string, index.ItemCount)

				o := index.Offset
				for s := 0; int64(s) < index.ItemCount; s++ {
					// calculate string length
					var j int64
					for j = 0; store[int64(j)+o] != 0; j++ {
					}

					vals[s] = string(store[o : o+j])
					o += j + 1

				}

				index.Value = vals

				break
			}

			// save in array
			h.Indexes[x] = index
		}

		// add header
		p.Headers = append(p.Headers, h)

		// calculate location of next header by padding to a multiple of 8
		o := 8 - int(math.Mod(float64(h.Length), 8))

		// seek to next header
		if o > 0 {
			pad := make([]byte, o)
			n, err = r.Read(pad)
			if err != nil {
				return nil, fmt.Errorf("Error seeking beyond header padding of %d bytes: %v", o, err)
			}

			if n != o {
				return nil, fmt.Errorf("Error seeking beyond header padding of %d bytes: only %d bytes returned", o, n)
			}
		}
	}

	return p, nil
}

// String reassembles package metadata to form a standard rpm package name;
// including the package name, version, release and architecture.
func (c *Package) String() string {
	return fmt.Sprintf("%s-%s-%s.%s", c.Name(), c.Version(), c.Release(), c.Architecture())
}

func (c *Package) Name() string {
	return c.Headers[1].Indexes.GetString(1000)
}

func (c *Package) Version() string {
	return c.Headers[1].Indexes.GetString(1001)
}

func (c *Package) Release() string {
	return c.Headers[1].Indexes.GetString(1002)
}

func (c *Package) Epoch() time.Time {
	return c.Headers[1].Indexes.GetTime(1003)
}

func (c *Package) Summary() []string {
	return c.Headers[1].Indexes.GetStrings(1004)
}

func (c *Package) Description() []string {
	return c.Headers[1].Indexes.GetStrings(1005)
}

func (c *Package) BuildTime() time.Time {
	return c.Headers[1].Indexes.GetTime(1006)
}

func (c *Package) BuildHost() string {
	return c.Headers[1].Indexes.GetString(1007)
}

func (c *Package) InstallTime() time.Time {
	return c.Headers[1].Indexes.GetTime(1008)
}

func (c *Package) Size() int64 {
	return c.Headers[1].Indexes.GetInt(1009)
}

func (c *Package) Distribution() string {
	return c.Headers[1].Indexes.GetString(1010)
}

func (c *Package) Vendor() string {
	return c.Headers[1].Indexes.GetString(1011)
}

func (c *Package) GIFImage() []byte {
	return c.Headers[1].Indexes.GetBytes(1012)
}

func (c *Package) XPMImage() []byte {
	return c.Headers[1].Indexes.GetBytes(1013)
}

func (c *Package) License() string {
	return c.Headers[1].Indexes.GetString(1014)
}

func (c *Package) Packager() string {
	return c.Headers[1].Indexes.GetString(1015)
}

func (c *Package) Groups() []string {
	return c.Headers[1].Indexes.GetStrings(1016)
}

func (c *Package) ChangeLog() []string {
	return c.Headers[1].Indexes.GetStrings(1017)
}

func (c *Package) Source() []string {
	return c.Headers[1].Indexes.GetStrings(1018)
}

func (c *Package) Patch() []string {
	return c.Headers[1].Indexes.GetStrings(1019)
}

func (c *Package) URL() string {
	return c.Headers[1].Indexes.GetString(1020)
}

func (c *Package) OperatingSystem() string {
	return c.Headers[1].Indexes.GetString(1021)
}

func (c *Package) Architecture() string {
	return c.Headers[1].Indexes.GetString(1022)
}

func (c *Package) PreInstallScript() string {
	return c.Headers[1].Indexes.GetString(1023)
}

func (c *Package) PostInstallScript() string {
	return c.Headers[1].Indexes.GetString(1024)
}

func (c *Package) PreUninstallScript() string {
	return c.Headers[1].Indexes.GetString(1025)
}

func (c *Package) PostUninstallScript() string {
	return c.Headers[1].Indexes.GetString(1026)
}

func (c *Package) OldFilenames() []string {
	return c.Headers[1].Indexes.GetStrings(1027)
}

func (c *Package) Icon() []byte {
	return c.Headers[1].Indexes.GetBytes(1043)
}

func (c *Package) SourceRPM() string {
	return c.Headers[1].Indexes.GetString(1044)
}

func (c *Package) Provides() []string {
	return c.Headers[1].Indexes.GetStrings(1047)
}

func (c *Package) Requires() []string {
	return c.Headers[1].Indexes.GetStrings(1049)
}

func (c *Package) RPMVersion() string {
	return c.Headers[1].Indexes.GetString(1064)
}

func (c *Package) Platform() string {
	return c.Headers[1].Indexes.GetString(1132)
}
