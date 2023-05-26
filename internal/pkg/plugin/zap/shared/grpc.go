package zapShared

import (
	"context"
	zapProto "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/proto"
	zapService "github.com/easysoft/zentaoatf/internal/pkg/plugin/zap/service"
)

// GRPCClient is an implementation of ZapInterface that talks over RPC.
type GRPCClient struct{ client zapProto.ZapClient }

func (m *GRPCClient) Put(key string, value []byte) error {
	_, err := m.client.Put(context.Background(), &zapProto.PutRequest{
		Key:   key,
		Value: value,
	})
	return err
}

func (m *GRPCClient) Get(key string) ([]byte, error) {
	resp, err := m.client.Get(context.Background(), &zapProto.GetRequest{
		Key: key,
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

type GRPCServer struct {
	zapProto.UnimplementedZapServer
	Impl zapService.ZapInterface
}

func (m *GRPCServer) Put(ctx context.Context, req *zapProto.PutRequest) (*zapProto.Empty, error) {
	return &zapProto.Empty{}, m.Impl.Put(req.Key, req.Value)
}

func (m *GRPCServer) Get(ctx context.Context, req *zapProto.GetRequest) (*zapProto.GetResponse, error) {
	v, err := m.Impl.Get(req.Key)
	return &zapProto.GetResponse{Value: v}, err
}
