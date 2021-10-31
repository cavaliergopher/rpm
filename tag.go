package rpm

const (
	TagTypeNull TagType = iota
	TagTypeChar
	TagTypeInt8
	TagTypeInt16
	TagTypeInt32
	TagTypeInt64
	TagTypeString
	TagTypeBinary
	TagTypeStringArray
	TagTypeI18NString
)

var tagTypeNames = []string{
	"NULL",
	"CHAR",
	"INT8",
	"INT16",
	"INT32",
	"INT64",
	"STRING",
	"BIN",
	"STRING_ARRAY",
	"I18NSTRING",
}

// TagType describes the data type of a tag's value.
type TagType int

func (i TagType) String() string {
	if i > 0 && int(i) < len(tagTypeNames) {
		return tagTypeNames[i]
	}
	return "UNKNOWN"
}

// Tag is an rpm header entry and its associated data value. Once the data type
// is known, use the associated value method to retrieve the tag value.
//
// All Tag methods will return their zero value if the underlying data type is
// a different type or if the tag is nil.
type Tag struct {
	ID    int
	Type  TagType
	Value interface{}
}

// StringSlice returns a slice of strings or nil if the index is not a string
// slice value.
//
// Use StringSlice for all STRING, STRING_ARRAY and I18NSTRING data types.
func (c *Tag) StringSlice() []string {
	if c == nil || c.Value == nil {
		return nil
	}
	if v, ok := c.Value.([]string); ok {
		return v
	}
	return nil
}

// String returns a string or an empty string if the index is not a string
// value.
//
// Use String for all STRING, STRING_ARRAY and I18NSTRING data types.
//
// This is not intended to implement fmt.Stringer. To format the tag using its
// identifier, use Tag.ID. To format the tag's value, use Tag.Value.
func (c *Tag) String() string {
	v := c.StringSlice()
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// Int64Slice returns a slice of int64s or nil if the index is not a numerical
// slice value. All integer types are cast to int64.
//
// Use Int64Slice for all INT16, INT32 and INT64 data types.
func (c *Tag) Int64Slice() []int64 {
	if c == nil || c.Value == nil {
		return nil
	}
	if v, ok := c.Value.([]int64); ok {
		return v
	}
	return nil
}

// Int64 returns an int64 if the index is not a numerical value. All integer
// types are cast to int64.
//
// Use Int64 for all INT16, INT32 and INT64 data types.
func (c *Tag) Int64() int64 {
	v := c.Int64Slice()
	if len(v) > 0 {
		return v[0]
	}
	return 0
}

// Bytes returns a slice of bytes or nil if the index is not a byte slice value.
//
// Use Bytes for all CHAR, INT8 and BIN data types.
func (c *Tag) Bytes() []byte {
	if c == nil || c.Value == nil {
		return nil
	}
	if v, ok := c.Value.([]byte); ok {
		return v
	}
	return nil
}
