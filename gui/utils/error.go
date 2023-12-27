package utils

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Must[T any](val T, err error) T {
	HandleErr(err)
	return val
}
