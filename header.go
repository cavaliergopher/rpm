package rpm

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

const r_MaxHeaderSize = 33554432

// A Header stores metadata about an rpm package.
type Header struct {
	Version int
	Tags    map[int]*Tag
}

// GetTag returns the tag with the given identifier.
//
// Nil is returned if the specified tag does not exist or the header is nil.
func (c *Header) GetTag(id int) *Tag {
	if c == nil || len(c.Tags) == 0 {
		return nil
	}
	return c.Tags[id]
}

type rpmHeader [16]byte

func (b rpmHeader) Magic() []byte   { return b[:3] }
func (b rpmHeader) Version() int    { return int(b[3]) }
func (b rpmHeader) IndexCount() int { return int(binary.BigEndian.Uint32(b[8:12])) }
func (b rpmHeader) Size() int       { return int(binary.BigEndian.Uint32(b[12:16])) }

type rpmIndex [16]byte

func (b rpmIndex) Tag() int        { return int(binary.BigEndian.Uint32(b[:4])) }
func (b rpmIndex) Type() TagType   { return TagType(binary.BigEndian.Uint32(b[4:8])) }
func (b rpmIndex) Offset() int     { return int(binary.BigEndian.Uint32(b[8:12])) }
func (b rpmIndex) ValueCount() int { return int(binary.BigEndian.Uint32(b[12:16])) }

// readHeader reads an RPM package file header structure from r.
func readHeader(r io.Reader, pad bool) (*Header, error) {
	// decode the header structure header
	var hdrBytes rpmHeader
	if _, err := r.Read(hdrBytes[:]); err != nil {
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
			hdrBytes.Size(),
		)
	}

	// decode the index
	indexBytes := make([]rpmIndex, hdrBytes.IndexCount())
	for i := 0; i < len(indexBytes); i++ {
		if _, err := r.Read(indexBytes[i][:]); err != nil {
			return nil, err
		}
		if indexBytes[i].Offset() >= hdrBytes.Size() {
			return nil, errorf(
				"offset of index %d is out of range: %s",
				i,
				indexBytes[i].Offset(),
			)
		}
	}

	// decode the store
	tags := make(map[int]*Tag, len(indexBytes))
	buf := make([]byte, hdrBytes.Size())
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	for i, ix := range indexBytes {
		if ix.ValueCount() < 1 {
			return nil, errorf("invalid value count for index %d: %d", i, ix.ValueCount())
		}
		o := ix.Offset()
		var v interface{}
		switch ix.Type() {
		case TagTypeBinary, TagTypeChar, TagTypeInt8:
			if o+ix.ValueCount() > len(buf) {
				switch ix.Type() {
				case TagTypeBinary:
					return nil, errorf("binary value for index %d is out of range", i+1)
				case TagTypeChar:
					return nil, errorf("uint8 value for index %d is out of range", i+1)
				case TagTypeInt8:
					return nil, errorf("int8 value for index %d is out of range", i+1)
				}
				return nil, errorf("value for index %d is out of range", i+1)
			}
			a := make([]byte, ix.ValueCount())
			copy(a, buf[o:o+ix.ValueCount()])
			v = a

		case TagTypeInt16:
			a := make([]int64, ix.ValueCount())
			for v := 0; v < ix.ValueCount(); v++ {
				if o+2 > len(buf) {
					return nil, errorf("int16 value for index %d is out of range", i+1)
				}
				a[v] = int64(binary.BigEndian.Uint16(buf[o : o+2]))
				o += 2
			}
			v = a

		case TagTypeInt32:
			a := make([]int64, ix.ValueCount())
			for v := 0; v < ix.ValueCount(); v++ {
				if o+4 > len(buf) {
					return nil, errorf("int32 value for index %d is out of range", i+1)
				}
				a[v] = int64(binary.BigEndian.Uint32(buf[o : o+4]))
				o += 4
			}

			v = a

		case TagTypeInt64:
			a := make([]int64, ix.ValueCount())
			for v := 0; v < ix.ValueCount(); v++ {
				if o+8 > len(buf) {
					// TODO: better errors
					return nil, errorf("int64 value for index %d is out of range", i+1)
				}
				a[v] = int64(binary.BigEndian.Uint64(buf[o : o+8]))
				o += 8
			}
			v = a

		case TagTypeString, TagTypeStringArray, TagTypeI18NString:
			// allow at least one byte per string
			if o+ix.ValueCount() > len(buf) {
				return nil, errorf("[]string value for index %d is out of range", i+1)
			}
			a := make([]string, ix.ValueCount())
			for s := 0; s < ix.ValueCount(); s++ {
				// calculate string length
				var j int
				for j = 0; (o+j) < len(buf) && buf[o+j] != 0; j++ {
				}
				if j == len(buf) {
					return nil, errorf("string value for index %d is out of range", i+1)
				}
				a[s] = string(buf[o : o+j])
				o += j + 1
			}
			v = a

		case TagTypeNull:
			// nothing to do here

		default:
			// unknown data type
			return nil, errorf("unknown index data type: %0X", ix.Type())
		}
		tags[ix.Tag()] = &Tag{
			ID:    ix.Tag(),
			Type:  ix.Type(),
			Value: v,
		}
	}

	// pad to next header
	padding := int64(8-(hdrBytes.Size()%8)) % 8
	if pad && padding != 0 {
		if _, err := io.CopyN(ioutil.Discard, r, padding); err != nil {
			return nil, err
		}
	}

	return &Header{
		Version: hdrBytes.Version(),
		Tags:    tags,
	}, nil
}
