package conf

func SetExitOnError(val bool) {
	ExitOnError = val
}

func GetExitOnError() bool {
	return ExitOnError
}
