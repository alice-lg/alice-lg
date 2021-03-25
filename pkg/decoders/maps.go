package decoders

// MapGet retrievs a key from an expected map
// it falls back if the input is not a map
// or the key was not found.
func MapGet(m interface{}, key string, fallback interface{}) interface{} {
	smap, ok := m.(map[string]interface{})
	if !ok {
		return fallback
	}
	val, ok := smap[key]
	if !ok {
		return fallback
	}
	return val
}
