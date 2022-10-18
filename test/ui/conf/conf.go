package conf

func SetExitOnError(val bool) {
	ExitAllOnError = val
}

func GetExitOnError() bool {
	return ExitAllOnError
}
func SetShowErr(val bool) {
	ShowErr = val
}

func GetShowErr() bool {
	return ShowErr
}

func DisableErr() {
	ShowErr = false
	ExitAllOnError = false
}

func EnableErr() {
	ShowErr = true
	ExitAllOnError = true
}
