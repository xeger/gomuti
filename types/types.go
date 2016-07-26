package types

// Widen numeric types to the largest possible representation of that type.
func widen(v interface{}) interface{} {
	switch v.(type) {
	case int8:
		return int64(v.(int8))
	case int16:
		return int64(v.(int16))
	case int32:
		return int64(v.(int32))
	case int64:
		return int64(v.(int64))
	case uint8:
		return uint64(v.(uint8))
	case uint16:
		return uint64(v.(uint16))
	case uint32:
		return uint64(v.(uint32))
	case uint64:
		return uint64(v.(uint64))
	case float32:
		return float64(v.(float32))
	default:
		return v
	}
}
