package conf

func SetExitOnError(val bool) {
	ExitAllOnError = val
}

func GetExitOnError() bool {
	return ExitAllOnError
}
