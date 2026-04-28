package felo

import (
	"context"

	"google.golang.org/grpc"
)

type RequestRideRequest struct {
	CustomerID   string     `json:"customer_id"`
	Pickup       Coordinate `json:"pickup"`
	Destination  Coordinate `json:"destination"`
	FareEstimate int64      `json:"fare_estimate"`
}

type RideByIDRequest struct {
	RideID string `json:"ride_id"`
}

type GenerateNowQRRequest struct {
	CustomerID   string     `json:"customer_id"`
	Destination  Coordinate `json:"destination"`
	FareEstimate int64      `json:"fare_estimate"`
}

type ScanNowQRRequest struct {
	QRCode   string `json:"qr_code"`
	DriverID string `json:"driver_id"`
}

type AcceptNowQRRequest struct {
	TripID    string `json:"trip_id"`
	DriverID  string `json:"driver_id"`
}

type PingRequest struct{}

type PingResponse struct {
	Status string `json:"status"`
}

type RideServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	RequestRide(context.Context, *RequestRideRequest) (*Ride, error)
	StartRide(context.Context, *RideByIDRequest) (*Ride, error)
	CompleteRide(context.Context, *RideByIDRequest) (*Ride, error)
	GenerateNowQR(context.Context, *GenerateNowQRRequest) (*QRSession, error)
	ScanNowQR(context.Context, *ScanNowQRRequest) (*QRSession, error)
	AcceptNowQR(context.Context, *AcceptNowQRRequest) (*Ride, error)
}

type RideServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	RequestRide(ctx context.Context, in *RequestRideRequest, opts ...grpc.CallOption) (*Ride, error)
	StartRide(ctx context.Context, in *RideByIDRequest, opts ...grpc.CallOption) (*Ride, error)
	CompleteRide(ctx context.Context, in *RideByIDRequest, opts ...grpc.CallOption) (*Ride, error)
	GenerateNowQR(ctx context.Context, in *GenerateNowQRRequest, opts ...grpc.CallOption) (*QRSession, error)
	ScanNowQR(ctx context.Context, in *ScanNowQRRequest, opts ...grpc.CallOption) (*QRSession, error)
	AcceptNowQR(ctx context.Context, in *AcceptNowQRRequest, opts ...grpc.CallOption) (*Ride, error)
}

type rideServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRideServiceClient(cc grpc.ClientConnInterface) RideServiceClient {
	return &rideServiceClient{cc: cc}
}

func (c *rideServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/felo.ride.v1.RideService/Ping", in, out, opts...)
	return out, err
}

func (c *rideServiceClient) RequestRide(ctx context.Context, in *RequestRideRequest, opts ...grpc.CallOption) (*Ride, error) {
	out := new(Ride)
	err := c.cc.Invoke(ctx, "/felo.ride.v1.RideService/RequestRide", in, out, opts...)
	return out, err
}

func (c *rideServiceClient) StartRide(ctx context.Context, in *RideByIDRequest, opts ...grpc.CallOption) (*Ride, error) {
	out := new(Ride)
	err := c.cc.Invoke(ctx, "/felo.ride.v1.RideService/StartRide", in, out, opts...)
	return out, err
}

func (c *rideServiceClient) CompleteRide(ctx context.Context, in *RideByIDRequest, opts ...grpc.CallOption) (*Ride, error) {
	out := new(Ride)
	err := c.cc.Invoke(ctx, "/felo.ride.v1.RideService/CompleteRide", in, out, opts...)
	return out, err
}

func (c *rideServiceClient) GenerateNowQR(ctx context.Context, in *GenerateNowQRRequest, opts ...grpc.CallOption) (*QRSession, error) {
	out := new(QRSession)
	err := c.cc.Invoke(ctx, "/felo.ride.v1.RideService/GenerateNowQR", in, out, opts...)
	return out, err
}

func (c *rideServiceClient) ScanNowQR(ctx context.Context, in *ScanNowQRRequest, opts ...grpc.CallOption) (*QRSession, error) {
	out := new(QRSession)
	err := c.cc.Invoke(ctx, "/felo.ride.v1.RideService/ScanNowQR", in, out, opts...)
	return out, err
}

func (c *rideServiceClient) AcceptNowQR(ctx context.Context, in *AcceptNowQRRequest, opts ...grpc.CallOption) (*Ride, error) {
	out := new(Ride)
	err := c.cc.Invoke(ctx, "/felo.ride.v1.RideService/AcceptNowQR", in, out, opts...)
	return out, err
}

func RegisterRideServiceServer(s grpc.ServiceRegistrar, srv RideServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "felo.ride.v1.RideService",
		HandlerType: (*RideServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			unaryMethod("Ping", srv.Ping),
			unaryMethod("RequestRide", srv.RequestRide),
			unaryMethod("StartRide", srv.StartRide),
			unaryMethod("CompleteRide", srv.CompleteRide),
			unaryMethod("GenerateNowQR", srv.GenerateNowQR),
			unaryMethod("ScanNowQR", srv.ScanNowQR),
			unaryMethod("AcceptNowQR", srv.AcceptNowQR),
		},
	}, srv)
}
