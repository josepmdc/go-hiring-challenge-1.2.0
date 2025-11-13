package fn

// Map iterates over the slice, applying the provided function on each item.
// If the provided slice is nil, nil will be returned
// If the provided slice is zero, a zero slice willl be returned
func Map[T, R any](items []T, callback func(item T) R) []R {
	if items == nil {
		return nil
	}

	result := make([]R, 0, len(items))
	for _, item := range items {
		result = append(result, callback(item))
	}
	return result
}

func Ptr[T comparable](item T) *T {
	return &item
}
