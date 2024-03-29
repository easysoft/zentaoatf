package zapService

import (
	"fmt"
	"os"
)

// implement ZapInterface
type ZapService struct{}

func (ZapService) Put(key string, value []byte) error {
	value = []byte(fmt.Sprintf("%s\nWritten from plugin-go-grpc", string(value)))
	return os.WriteFile("zap_"+key, value, 0644)
}

func (ZapService) Get(key string) ([]byte, error) {
	return os.ReadFile("zap_" + key)
}
