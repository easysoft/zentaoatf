package zapService

type ZapInterface interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
}
