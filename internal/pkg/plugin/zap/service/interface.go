package zapService

// ZapInterface is the interface that we're exposing as a plugin.
type ZapInterface interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
}
