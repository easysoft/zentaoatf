package action

func SwitchWorkDir(dir string) error {
	Set("workDir", dir, false)

	return nil
}
