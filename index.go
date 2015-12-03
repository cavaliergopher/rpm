package rpm

import (
	"time"
)

type IndexEntry struct {
	Tag       int64
	Type      int64
	Offset    int64
	ItemCount int64
	Value     interface{}
}

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

func (c IndexEntries) Get(tag int64) *IndexEntry {
	for _, e := range c {
		if e.Tag == tag {
			return &e
		}
	}

	return nil
}

func (c IndexEntries) GetString(tag int64) string {
	i := c.Get(tag)
	if i == nil || i.Value == nil {
		return ""
	}

	s := i.Value.([]string)

	return s[0]
}

func (c IndexEntries) GetStringArray(tag int64) []string {
	i := c.Get(tag)
	if i == nil || i.Value == nil {
		return nil
	}

	return i.Value.([]string)
}

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

func (c IndexEntries) GetBytes(tag int64) []byte {
	i := c.Get(tag)
	if i == nil || i.Value == nil {
		return nil
	}

	return i.Value.([]byte)
}

func (c IndexEntries) GetTime(tag int64) time.Time {
	i := c.Get(tag)
	if i == nil || i.Value == nil {
		return time.Unix(0, 0)
	}

	return time.Unix(int64(i.Value.(int32)), 0)
}
