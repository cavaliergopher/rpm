package rpm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"time"
)

// A PackageFile is an RPM package definition loaded directly from the pacakge
// file itself.
type PackageFile struct {
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
	Length     int
	Indexes    IndexEntries
}

// Headers is an array of Header structs.
type Headers []Header

// OpenPackageFile reads a rpm package from the file systems and returns a pointer
// to it.
func OpenPackageFile(path string) (*PackageFile, error) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error opening RPM file: %s", err)
	}
	defer f.Close()

	return ReadPackageFile(f)
}

// ReadPackageFile reads a rpm package file from a stream and returns a pointer
// to it.
func ReadPackageFile(r io.Reader) (*PackageFile, error) {
	p := &PackageFile{}

	// See: http://www.rpm.org/max-rpm/s1-rpm-file-format-rpm-file-format.html

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
	p.Lead.Type = int(binary.BigEndian.Uint16(lead[7:9]))
	p.Lead.Architecture = int(binary.BigEndian.Uint16(lead[9:11]))
	p.Lead.Name = string(lead[10:77])
	p.Lead.OperatingSystem = int(binary.BigEndian.Uint16(lead[76:78]))
	p.Lead.SignatureType = int(binary.BigEndian.Uint16(lead[78:80]))

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
		h.IndexCount = int(binary.BigEndian.Uint32(header[8:12]))
		h.Length = int(binary.BigEndian.Uint32(header[12:16]))
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

			index.Tag = int(binary.BigEndian.Uint32(indexes[o : o+4]))
			index.Type = int(binary.BigEndian.Uint32(indexes[o+4 : o+8]))
			index.Offset = int(binary.BigEndian.Uint32(indexes[o+8 : o+12]))
			index.ItemCount = int(binary.BigEndian.Uint32(indexes[o+12 : o+16]))
			h.Indexes[x] = index
		}

		// read the "store"
		store := make([]byte, h.Length)
		n, err = r.Read(store)
		if err != nil {
			return nil, fmt.Errorf("Error reading store for header %d: %v", i, err)
		}

		if n != h.Length {
			return nil, fmt.Errorf("Error reading store for header %d: only %d bytes returned", i, n)
		}

		// parse the value of each index from the store
		for x := 0; x < h.IndexCount; x++ {
			index := h.Indexes[x]
			o := index.Offset

			switch index.Type {
			case IndexDataTypeChar:
				vals := make([]uint8, index.ItemCount)
				for v := 0; v < index.ItemCount; v++ {
					vals[v] = uint8(store[o])
					o += 1
				}

				index.Value = vals

			case IndexDataTypeInt8:
				vals := make([]int8, index.ItemCount)
				for v := 0; v < index.ItemCount; v++ {
					vals[v] = int8(store[o])
					o += 1
				}

				index.Value = vals

			case IndexDataTypeInt16:
				vals := make([]int16, index.ItemCount)
				for v := 0; v < index.ItemCount; v++ {
					vals[v] = int16(binary.BigEndian.Uint16(store[o : o+2]))
					o += 2
				}

				index.Value = vals

			case IndexDataTypeInt32:
				vals := make([]int32, index.ItemCount)
				for v := 0; v < index.ItemCount; v++ {
					vals[v] = int32(binary.BigEndian.Uint32(store[o : o+4]))
					o += 4
				}

				index.Value = vals

			case IndexDataTypeInt64:
				vals := make([]int64, index.ItemCount)
				for v := 0; v < index.ItemCount; v++ {
					vals[v] = int64(binary.BigEndian.Uint64(store[o : o+8]))
					o += 8
				}

				index.Value = vals

			case IndexDataTypeBinary:
				b := make([]byte, index.ItemCount)
				copy(b, store[o:o+index.ItemCount])

				index.Value = b

			case IndexDataTypeString, IndexDataTypeStringArray, IndexDataTypeI8NString:
				vals := make([]string, index.ItemCount)

				for s := 0; s < int(index.ItemCount); s++ {
					// calculate string length
					var j int
					for j = 0; store[j+o] != 0; j++ {
					}

					vals[s] = string(store[o : o+j])
					o += j + 1
				}

				index.Value = vals
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

// dependencies translates the given tag values into a slice of package
// relationships such as provides, conflicts, obsoletes and requires.
func (c *PackageFile) dependencies(nevrsTagId, flagsTagId, namesTagId, versionsTagId int) Dependencies {
	// TODO: Implement NEVRS tags

	flgs := c.Headers[1].Indexes.IntsByTag(flagsTagId)
	names := c.Headers[1].Indexes.StringsByTag(namesTagId)
	vers := c.Headers[1].Indexes.StringsByTag(versionsTagId)

	deps := make(Dependencies, len(names))
	for i := 0; i < len(names); i++ {
		deps[i] = NewDependency(int(flgs[i]), names[i], 0, vers[i], "")
	}

	return deps
}

// String reassembles package metadata to form a standard rpm package name;
// including the package name, version, release and architecture.
func (c *PackageFile) String() string {
	return fmt.Sprintf("%s-%s-%s.%s", c.Name(), c.Version(), c.Release(), c.Architecture())
}

// For tag definitions, see:
// https://github.com/rpm-software-management/rpm/blob/master/lib/rpmtag.h

func (c *PackageFile) Name() string {
	return c.Headers[1].Indexes.StringByTag(1000)
}

func (c *PackageFile) Version() string {
	return c.Headers[1].Indexes.StringByTag(1001)
}

func (c *PackageFile) Release() string {
	return c.Headers[1].Indexes.StringByTag(1002)
}

func (c *PackageFile) Epoch() int64 {
	return c.Headers[1].Indexes.IntByTag(1003)
}

func (c *PackageFile) Requires() Dependencies {
	return c.dependencies(5041, 1048, 1049, 1050)
}

func (c *PackageFile) Provides() Dependencies {
	return c.dependencies(5042, 1112, 1047, 1113)
}

func (c *PackageFile) Conflicts() Dependencies {
	return c.dependencies(5044, 1053, 1054, 1055)
}

func (c *PackageFile) Obsoletes() Dependencies {
	return c.dependencies(5043, 1114, 1090, 1115)
}

func (c *PackageFile) Summary() []string {
	return c.Headers[1].Indexes.StringsByTag(1004)
}

func (c *PackageFile) Description() []string {
	return c.Headers[1].Indexes.StringsByTag(1005)
}

func (c *PackageFile) BuildTime() time.Time {
	return c.Headers[1].Indexes.TimeByTag(1006)
}

func (c *PackageFile) BuildHost() string {
	return c.Headers[1].Indexes.StringByTag(1007)
}

func (c *PackageFile) InstallTime() time.Time {
	return c.Headers[1].Indexes.TimeByTag(1008)
}

func (c *PackageFile) Size() int64 {
	return c.Headers[1].Indexes.IntByTag(1009)
}

func (c *PackageFile) Distribution() string {
	return c.Headers[1].Indexes.StringByTag(1010)
}

func (c *PackageFile) Vendor() string {
	return c.Headers[1].Indexes.StringByTag(1011)
}

func (c *PackageFile) GIFImage() []byte {
	return c.Headers[1].Indexes.BytesByTag(1012)
}

func (c *PackageFile) XPMImage() []byte {
	return c.Headers[1].Indexes.BytesByTag(1013)
}

func (c *PackageFile) License() string {
	return c.Headers[1].Indexes.StringByTag(1014)
}

func (c *PackageFile) PackageFiler() string {
	return c.Headers[1].Indexes.StringByTag(1015)
}

func (c *PackageFile) Groups() []string {
	return c.Headers[1].Indexes.StringsByTag(1016)
}

func (c *PackageFile) ChangeLog() []string {
	return c.Headers[1].Indexes.StringsByTag(1017)
}

func (c *PackageFile) Source() []string {
	return c.Headers[1].Indexes.StringsByTag(1018)
}

func (c *PackageFile) Patch() []string {
	return c.Headers[1].Indexes.StringsByTag(1019)
}

func (c *PackageFile) URL() string {
	return c.Headers[1].Indexes.StringByTag(1020)
}

func (c *PackageFile) OperatingSystem() string {
	return c.Headers[1].Indexes.StringByTag(1021)
}

func (c *PackageFile) Architecture() string {
	return c.Headers[1].Indexes.StringByTag(1022)
}

func (c *PackageFile) PreInstallScript() string {
	return c.Headers[1].Indexes.StringByTag(1023)
}

func (c *PackageFile) PostInstallScript() string {
	return c.Headers[1].Indexes.StringByTag(1024)
}

func (c *PackageFile) PreUninstallScript() string {
	return c.Headers[1].Indexes.StringByTag(1025)
}

func (c *PackageFile) PostUninstallScript() string {
	return c.Headers[1].Indexes.StringByTag(1026)
}

func (c *PackageFile) OldFilenames() []string {
	return c.Headers[1].Indexes.StringsByTag(1027)
}

func (c *PackageFile) Icon() []byte {
	return c.Headers[1].Indexes.BytesByTag(1043)
}

func (c *PackageFile) SourceRPM() string {
	return c.Headers[1].Indexes.StringByTag(1044)
}

func (c *PackageFile) RPMVersion() string {
	return c.Headers[1].Indexes.StringByTag(1064)
}

func (c *PackageFile) Platform() string {
	return c.Headers[1].Indexes.StringByTag(1132)
}
