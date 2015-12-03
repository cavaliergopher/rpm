package rpm

import (
	"time"
)

// An IndexEntry is a rpm key/value tag stored in the package header.
type IndexEntry struct {
	Tag       int64
	Type      int64
	Offset    int64
	ItemCount int64
	Value     interface{}
}

// IndexEntries is an array of IndexEntry structs.
type IndexEntries []IndexEntry

const (
	IndexDataTypeNull int64 = iota
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

// Get returns a pointer to an IndexEntry with the given tag ID or nil if the
// tag is not found.
func (c IndexEntries) Get(tag int64) *IndexEntry {
	for _, e := range c {
		if e.Tag == tag {
			return &e
		}
	}

	return nil
}

// GetString returns the string value of an IndexEntry or an empty string if
// the tag is not found or has no value.
func (c IndexEntries) GetString(tag int64) string {
	i := c.Get(tag)
	if i == nil || i.Value == nil {
		return ""
	}

	s := i.Value.([]string)

	return s[0]
}

// GetStrings returns the slice of string values of an IndexEntry or nil if the
// tag is not found or has no value.
func (c IndexEntries) GetStrings(tag int64) []string {
	i := c.Get(tag)
	if i == nil || i.Value == nil {
		return nil
	}

	return i.Value.([]string)
}

// GetInt returns the int64 value of an IndexEntry or 0 if the tag is not found
// or has no value. Values with a lower range (E.g. int8) are cast as an int64.
func (c IndexEntries) GetInt(tag int64) int64 {
	i := c.Get(tag)
	if i != nil && i.Value != nil {
		switch i.Type {
		case IndexDataTypeChar, IndexDataTypeInt8:
			return int64(i.Value.(int8))

		case IndexDataTypeInt16:
			return int64(i.Value.(int16))

		case IndexDataTypeInt32:
			return int64(i.Value.(int32))

		case IndexDataTypeInt64:
			return int64(i.Value.(int64))
		}
	}

	return 0
}

// GetBytes returns the raw value of an IndexEntry or nil if the tag is not
// found or has no value.
func (c IndexEntries) GetBytes(tag int64) []byte {
	i := c.Get(tag)
	if i == nil || i.Value == nil {
		return nil
	}

	return i.Value.([]byte)
}

// GetTime returns the value of an IndexEntry as a Go native timestamp or Epoch
// if the tag is not found or has no value.
func (c IndexEntries) GetTime(tag int64) time.Time {
	i := c.Get(tag)
	if i == nil || i.Value == nil {
		return time.Unix(0, 0)
	}

	return time.Unix(int64(i.Value.(int32)), 0)
}
