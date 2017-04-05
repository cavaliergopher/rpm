package rpm

import (
	"sort"
	"time"
)

// Header index value data types.
const (
	IndexDataTypeNull int = iota
	IndexDataTypeChar
	IndexDataTypeInt8
	IndexDataTypeInt16
	IndexDataTypeInt32
	IndexDataTypeInt64
	IndexDataTypeString
	IndexDataTypeBinary
	IndexDataTypeStringArray
	IndexDataTypeI8NString
)

// An IndexEntry is a rpm key/value tag stored in the package header.
type IndexEntry struct {
	Tag       int
	Type      int
	Offset    int
	ItemCount int
	Value     interface{}
}

// IndexList enables the fast retrieval of Indexes in O(log(n)). This will most
// frequently be faster than a O(1) map lookup as the number of indexes
// is usually small enough for the constant-time complexity of a map to be
// disadvantageous.
type IndexList interface {
	Add(...*IndexEntry)
	IndexByTag(tag int) *IndexEntry
	StringByTag(tag int) string
	StringsByTag(tag int) []string
	IntByTag(tag int) int64
	IntsByTag(tag int) []int64
	BytesByTag(tag int) []byte
	TimeByTag(tag int) time.Time
	TimesByTag(tag int) []time.Time
}

// IndexEntries is a list of IndexEntry structs.
type indexList struct {
	indexes []*IndexEntry
}

// NewIndexList returns an implementation of IndexList.
func NewIndexList() IndexList {
	return &indexList{
		indexes: make([]*IndexEntry, 0),
	}
}

func (c *indexList) Len() int {
	return len(c.indexes)
}

func (c *indexList) Less(i, j int) bool {
	return c.indexes[i].Tag < c.indexes[j].Tag
}

func (c *indexList) Swap(i, j int) {
	c.indexes[i], c.indexes[j] = c.indexes[j], c.indexes[i]
}

// Add adds the given indexes to the index list. This is an expensive O(n log n)
// and is only meant to be called once.
func (c *indexList) Add(indexes ...*IndexEntry) {
	c.indexes = append(c.indexes, indexes...)
	sort.Sort(c)
}

// IndexByTag returns a pointer to an IndexEntry with the given tag ID or nil if
// the tag is not found.
func (c *indexList) IndexByTag(tag int) *IndexEntry {
	b := 0
	e := len(c.indexes)
	var m int
	var ix *IndexEntry

	//fmt.Printf("Searching for %d\n", tag)

loop:
	for {
		m = b + ((e - b) / 2)
		ix = c.indexes[m]

		//fmt.Printf("\tTrying b: %d, m: %d, e: %d, l: %d:", b, m, e, len(c.indexes))

		switch {
		case ix.Tag == tag:
			//fmt.Printf("found!\n")
			return ix

		case ix.Tag < tag:
			//fmt.Printf("%d too low\n", ix.Tag)
			if b == m {
				break loop
			}
			b = m

		case ix.Tag > tag:
			//fmt.Printf("%d too high\n", ix.Tag)
			if e == m {
				break loop
			}
			e = m
		}
	}

	//fmt.Printf("%d not found\n", tag)
	return nil

	/*
		// binary search slice of indexes, sorted by tag
		var f func(x []*IndexEntry) *IndexEntry
		f = func(x []*IndexEntry) *IndexEntry {
			m := int(len(x) / 2)
			ix := x[m]

			//fmt.Printf("\tRange: %d -> %d (%d) - ", x[0].Tag, x[len(x)-1].Tag, len(x))
			switch {
			case ix.Tag == tag:
				//fmt.Printf("Found at %d!\n", m)
				return ix

			case len(x) == 1:
				//fmt.Printf("Not found...\n")
				return nil

			case ix.Tag > tag:
				//fmt.Printf("Trying left of %d 0->%d\n", ix.Tag, m)
				return f(x[:m])

			case ix.Tag < tag:
				if m == len(x)-1 {
					//fmt.Printf("Not found...\n")
					return nil
				}

				//fmt.Printf("Trying right of %d %d->%d\n", ix.Tag, m+1, len(x))
				return f(x[m+1:])

			default:
				panic("this should never happen")
			}
		}

		//fmt.Printf("Searching for %d in %d items\n", tag, len(c.indexes))
		return f(c.indexes)
	*/

	/*
		for _, i := range c.indexes {
			if i.Tag == tag {
				return i
			}
		}

		return nil
	*/
}

// StringByTag returns the string value of an IndexEntry or an empty string if
// the tag is not found or has no value.
func (c *indexList) StringByTag(tag int) string {
	i := c.IndexByTag(tag)
	if i == nil || i.Value == nil {
		return ""
	}

	s := i.Value.([]string)

	return s[0]
}

// StringsByTag returns the slice of string values of an IndexEntry or nil if
// the tag is not found or has no value.
func (c *indexList) StringsByTag(tag int) []string {
	i := c.IndexByTag(tag)
	if i == nil || i.Value == nil {
		return nil
	}

	return i.Value.([]string)
}

// IntsByTag returns the int64 values of an IndexEntry or nil if the tag is not
// found or has no value. Values with a lower range (E.g. int8) are cast as an
// int64.
func (c *indexList) IntsByTag(tag int) []int64 {
	ix := c.IndexByTag(tag)
	if ix != nil && ix.Value != nil {
		vals := make([]int64, ix.ItemCount)

		for i := 0; i < int(ix.ItemCount); i++ {
			switch ix.Type {
			case IndexDataTypeChar, IndexDataTypeInt8:
				vals[i] = int64(ix.Value.([]int8)[i])

			case IndexDataTypeInt16:
				vals[i] = int64(ix.Value.([]int16)[i])

			case IndexDataTypeInt32:
				vals[i] = int64(ix.Value.([]int32)[i])

			case IndexDataTypeInt64:
				vals[i] = ix.Value.([]int64)[i]
			}
		}

		return vals
	}

	return nil
}

// IntByTag returns the int64 value of an IndexEntry or 0 if the tag is not found
// or has no value. Values with a lower range (E.g. int8) are cast as an int64.
func (c *indexList) IntByTag(tag int) int64 {
	i := c.IndexByTag(tag)
	if i != nil && i.Value != nil {
		switch i.Type {
		case IndexDataTypeChar, IndexDataTypeInt8:
			return int64(i.Value.([]int8)[0])

		case IndexDataTypeInt16:
			return int64(i.Value.([]int16)[0])

		case IndexDataTypeInt32:
			return int64(i.Value.([]int32)[0])

		case IndexDataTypeInt64:
			return int64(i.Value.([]int64)[0])
		}
	}

	return 0
}

// BytesByTag returns the raw value of an IndexEntry or nil if the tag is not
// found or has no value.
func (c *indexList) BytesByTag(tag int) []byte {
	i := c.IndexByTag(tag)
	if i == nil || i.Value == nil {
		return nil
	}

	return i.Value.([]byte)
}

// TimesByTag returns the value of an IndexEntry as a slice of Go native
// timestamps or nil if the tag is not found or has no value.
func (c *indexList) TimesByTag(tag int) []time.Time {
	ix := c.IndexByTag(tag)

	if ix == nil || ix.Value == nil {
		return nil
	}

	vals := make([]time.Time, ix.ItemCount)
	for i := 0; i < ix.ItemCount; i++ {
		vals[i] = time.Unix(int64(ix.Value.([]int32)[i]), 0)
	}

	return vals
}

// TimeByTag returns the value of an IndexEntry as a Go native timestamp or
// zero-time if the tag is not found or has no value.
func (c *indexList) TimeByTag(tag int) time.Time {
	vals := c.TimesByTag(tag)
	if vals == nil || len(vals) == 0 {
		return time.Time{}
	}

	return vals[0]
}
