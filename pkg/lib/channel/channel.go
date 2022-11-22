package channelUtils

func IsChanClose(ch chan int) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
