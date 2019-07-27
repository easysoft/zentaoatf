package action

func SwitchWorkDir(dir string) error {
	Set("workDir", dir, true)

	return nil
}
