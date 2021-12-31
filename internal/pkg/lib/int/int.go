package intUtils

func FindInArr(val int, arr []int) bool {
	for _, i := range arr {
		if val == i {
			return true
		}
	}

	return false
}

func FindUintInArr(val uint, arr []uint) bool {
	for _, i := range arr {
		if val == i {
			return true
		}
	}

	return false
}
