package felo

import (
	"context"

	"google.golang.org/grpc"
)

type AssignmentRequest struct {
	RideID string `json:"ride_id"`
}

type AssignmentResponse struct {
	DriverID string `json:"driver_id"`
}

type MatchingServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	GetAssignment(context.Context, *AssignmentRequest) (*AssignmentResponse, error)
}

type MatchingServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	GetAssignment(ctx context.Context, in *AssignmentRequest, opts ...grpc.CallOption) (*AssignmentResponse, error)
}

type matchingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMatchingServiceClient(cc grpc.ClientConnInterface) MatchingServiceClient {
	return &matchingServiceClient{cc: cc}
}

func (c *matchingServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/felo.matching.v1.MatchingService/Ping", in, out, opts...)
	return out, err
}

func (c *matchingServiceClient) GetAssignment(ctx context.Context, in *AssignmentRequest, opts ...grpc.CallOption) (*AssignmentResponse, error) {
	out := new(AssignmentResponse)
	err := c.cc.Invoke(ctx, "/felo.matching.v1.MatchingService/GetAssignment", in, out, opts...)
	return out, err
}

func RegisterMatchingServiceServer(s grpc.ServiceRegistrar, srv MatchingServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "felo.matching.v1.MatchingService",
		HandlerType: (*MatchingServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			unaryMethod("Ping", srv.Ping),
			unaryMethod("GetAssignment", srv.GetAssignment),
		},
	}, srv)
}
