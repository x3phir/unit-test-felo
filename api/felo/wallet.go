package felo

import (
	"context"

	"google.golang.org/grpc"
)

type BalanceRequest struct {
	OwnerID string `json:"owner_id"`
}

type WalletServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	GetDriverBalance(context.Context, *BalanceRequest) (*WalletBalance, error)
	GetCustomerBalance(context.Context, *BalanceRequest) (*WalletBalance, error)
}

type WalletServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	GetDriverBalance(ctx context.Context, in *BalanceRequest, opts ...grpc.CallOption) (*WalletBalance, error)
	GetCustomerBalance(ctx context.Context, in *BalanceRequest, opts ...grpc.CallOption) (*WalletBalance, error)
}

type walletServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWalletServiceClient(cc grpc.ClientConnInterface) WalletServiceClient {
	return &walletServiceClient{cc: cc}
}

func (c *walletServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/felo.wallet.v1.WalletService/Ping", in, out, opts...)
	return out, err
}

func (c *walletServiceClient) GetDriverBalance(ctx context.Context, in *BalanceRequest, opts ...grpc.CallOption) (*WalletBalance, error) {
	out := new(WalletBalance)
	err := c.cc.Invoke(ctx, "/felo.wallet.v1.WalletService/GetDriverBalance", in, out, opts...)
	return out, err
}

func (c *walletServiceClient) GetCustomerBalance(ctx context.Context, in *BalanceRequest, opts ...grpc.CallOption) (*WalletBalance, error) {
	out := new(WalletBalance)
	err := c.cc.Invoke(ctx, "/felo.wallet.v1.WalletService/GetCustomerBalance", in, out, opts...)
	return out, err
}

func RegisterWalletServiceServer(s grpc.ServiceRegistrar, srv WalletServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "felo.wallet.v1.WalletService",
		HandlerType: (*WalletServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			unaryMethod("Ping", srv.Ping),
			unaryMethod("GetDriverBalance", srv.GetDriverBalance),
			unaryMethod("GetCustomerBalance", srv.GetCustomerBalance),
		},
	}, srv)
}
