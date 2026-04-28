package felo

import (
	"context"

	"google.golang.org/grpc"
)

type LocationHistoryRequest struct {
	DriverID string `json:"driver_id"`
	From     string `json:"from"`
	To       string `json:"to"`
}

type LocationHistoryResponse struct {
	Samples []LocationSample `json:"samples"`
}

type LocationServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	ReportDriverLocation(context.Context, *LocationSample) (*PingResponse, error)
	GetDriverHistory(context.Context, *LocationHistoryRequest) (*LocationHistoryResponse, error)
}

type LocationServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	ReportDriverLocation(ctx context.Context, in *LocationSample, opts ...grpc.CallOption) (*PingResponse, error)
	GetDriverHistory(ctx context.Context, in *LocationHistoryRequest, opts ...grpc.CallOption) (*LocationHistoryResponse, error)
}

type locationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLocationServiceClient(cc grpc.ClientConnInterface) LocationServiceClient {
	return &locationServiceClient{cc: cc}
}

func (c *locationServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/felo.location.v1.LocationService/Ping", in, out, opts...)
	return out, err
}

func (c *locationServiceClient) ReportDriverLocation(ctx context.Context, in *LocationSample, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/felo.location.v1.LocationService/ReportDriverLocation", in, out, opts...)
	return out, err
}

func (c *locationServiceClient) GetDriverHistory(ctx context.Context, in *LocationHistoryRequest, opts ...grpc.CallOption) (*LocationHistoryResponse, error) {
	out := new(LocationHistoryResponse)
	err := c.cc.Invoke(ctx, "/felo.location.v1.LocationService/GetDriverHistory", in, out, opts...)
	return out, err
}

func RegisterLocationServiceServer(s grpc.ServiceRegistrar, srv LocationServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "felo.location.v1.LocationService",
		HandlerType: (*LocationServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			unaryMethod("Ping", srv.Ping),
			unaryMethod("ReportDriverLocation", srv.ReportDriverLocation),
			unaryMethod("GetDriverHistory", srv.GetDriverHistory),
		},
	}, srv)
}
