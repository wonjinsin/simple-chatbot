package shared

// Ptr returns a pointer to the given value
func Ptr[T any](v T) *T {
	return &v
}

// ValueOr returns the dereferenced value if ptr is not nil, otherwise returns the default value
func ValueOr[T any](ptr *T, defaultVal T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}

// IsNil checks if a pointer is nil
func IsNil[T any](ptr *T) bool {
	return ptr == nil
}
