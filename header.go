package rpm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// A Header stores metadata about a rpm package.
type Header struct {
	Version    int
	IndexCount int
	Length     int
	Indexes    IndexEntries
}

// Headers is an array of Header structs.
type Headers []Header

// ReadPackageHeader reads an RPM package file header structure from the given
// io.Reader.
//
// This function should only be used if you intend to read a package header
// structure in isolation.
func ReadPackageHeader(r io.Reader) (*Header, error) {
	// read the "header structure header"
	header := make([]byte, 16)
	n, err := r.Read(header)
	if err != nil {
		return nil, fmt.Errorf("Error reading package header: %v", err)
	}

	if n != 16 {
		return nil, fmt.Errorf("Error reading package header: only %d bytes returned", n)
	}

	// check magic number
	if 0 != bytes.Compare(header[:3], []byte{0x8E, 0xAD, 0xE8}) {
		return nil, fmt.Errorf("Bad magic number in package header")
	}

	// translate header
	h := &Header{}
	h.Version = int(header[3])
	h.IndexCount = int(binary.BigEndian.Uint32(header[8:12]))
	h.Length = int(binary.BigEndian.Uint32(header[12:16]))
	h.Indexes = make(IndexEntries, h.IndexCount)

	// read indexes
	indexLength := 16 * h.IndexCount
	indexes := make([]byte, indexLength)
	n, err = r.Read(indexes)
	if err != nil {
		return nil, fmt.Errorf("Error reading index entries for header: %v", err)
	}

	if n != indexLength {
		return nil, fmt.Errorf("Error reading index entries for header: only %d bytes returned", n)
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
		return nil, fmt.Errorf("Error reading header store: %v", err)
	}

	if n != h.Length {
		return nil, fmt.Errorf("Error reading header store: only %d bytes returned", n)
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

	// calculate location of next header by padding to a multiple of 8
	o := 8 - int(math.Mod(float64(h.Length), 8))

	// seek to next header
	if o > 0 && o < 8 {
		pad := make([]byte, o)
		n, err = r.Read(pad)
		if err != nil {
			return nil, fmt.Errorf("Error seeking beyond header padding of %d bytes: %v", o, err)
		}

		if n != o {
			return nil, fmt.Errorf("Error seeking beyond header padding of %d bytes: only %d bytes returned", o, n)
		}
	}

	return h, nil
}
