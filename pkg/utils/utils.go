package utils

// Must panics if e is not nil.
// Do not use in production!
// This is just a helper for the examples or test environment.
func Must[T any](t T, e error) T {
	if e != nil {
		panic(e)
	}

	return t
}
