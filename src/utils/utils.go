package utils

import "io/ioutil"

func ReadFile(filePth string) []byte {
	buf, err := ioutil.ReadFile(filePth)
	if err != nil {
		return nil
	}

	return buf
}
