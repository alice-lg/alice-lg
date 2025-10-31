package decoders

// MapGet retrieves a key from an expected map
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

// MapGetString retrieves a key from a map and
// asserts its type is a string. Otherwise fallback
// will be returned.
func MapGetString(m interface{}, key string, fallback string) string {
	val := MapGet(m, key, fallback)
	return val.(string)
}

// MapGetBool will retrieve a boolean value
// for a given key.
func MapGetBool(m interface{}, key string, fallback bool) bool {
	val := MapGet(m, key, fallback)
	return val.(bool)
}
