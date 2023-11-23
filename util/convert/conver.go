package convert

func To[T any, R any](source T, funcHandle func(source T) R) R {
	return funcHandle(source)
}

func ToList[T any, R any](source []T, funcHandle func(source []T) []R) (bool, []R) {

	if source == nil || len(source) == 0 {
		return false, nil
	}

	return true, funcHandle(source)
}
