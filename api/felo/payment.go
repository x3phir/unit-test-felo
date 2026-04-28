package felo

import (
	"context"

	"google.golang.org/grpc"
)

type PaymentStatusRequest struct {
	RideID string `json:"ride_id"`
}

type PaymentStatusResponse struct {
	Status string `json:"status"`
}

type PaymentServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	GetPaymentStatus(context.Context, *PaymentStatusRequest) (*PaymentStatusResponse, error)
}

type PaymentServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	GetPaymentStatus(ctx context.Context, in *PaymentStatusRequest, opts ...grpc.CallOption) (*PaymentStatusResponse, error)
}

type paymentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentServiceClient(cc grpc.ClientConnInterface) PaymentServiceClient {
	return &paymentServiceClient{cc: cc}
}

func (c *paymentServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/felo.payment.v1.PaymentService/Ping", in, out, opts...)
	return out, err
}

func (c *paymentServiceClient) GetPaymentStatus(ctx context.Context, in *PaymentStatusRequest, opts ...grpc.CallOption) (*PaymentStatusResponse, error) {
	out := new(PaymentStatusResponse)
	err := c.cc.Invoke(ctx, "/felo.payment.v1.PaymentService/GetPaymentStatus", in, out, opts...)
	return out, err
}

func RegisterPaymentServiceServer(s grpc.ServiceRegistrar, srv PaymentServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "felo.payment.v1.PaymentService",
		HandlerType: (*PaymentServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			unaryMethod("Ping", srv.Ping),
			unaryMethod("GetPaymentStatus", srv.GetPaymentStatus),
		},
	}, srv)
}
