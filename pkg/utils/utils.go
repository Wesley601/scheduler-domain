package utils

import "strconv"

// Must panics if e is not nil.
// Do not use in production!
// This is just a helper for the examples or test environment.
func Must[T any](t T, e error) T {
	if e != nil {
		panic(e)
	}

	return t
}

func ParseOptionalIntQueryParam(p string) (int, error) {
	if p == "" {
		return 1, nil
	}

	result, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}

	return result, nil
}
