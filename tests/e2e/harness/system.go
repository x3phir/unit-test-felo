//go:build e2e

package harness

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/felo/felo-backend/api/felo"
	"github.com/felo/felo-backend/internal/demo/app"
	"github.com/felo/felo-backend/tests/e2e/config"
	"github.com/rabbitmq/amqp091-go"
)

type SystemUnderTest struct {
	Health   HealthChecker
	Ride     RideClient
	Matching MatchingClient
	Wallet   WalletClient
	Payment  PaymentClient
	Location LocationClient
	Events   EventObserver
}

func NewSystemUnderTest(cfg config.Config) (SystemUnderTest, error) {
	rideConn, err := app.Dial(cfg.RideGRPCAddr, cfg.AuthJWT)
	if err != nil {
		return SystemUnderTest{}, err
	}
	matchConn, err := app.Dial(cfg.MatchingGRPCAddr, cfg.AuthJWT)
	if err != nil {
		return SystemUnderTest{}, err
	}
	walletConn, err := app.Dial(cfg.WalletGRPCAddr, cfg.AuthJWT)
	if err != nil {
		return SystemUnderTest{}, err
	}
	paymentConn, err := app.Dial(cfg.PaymentGRPCAddr, cfg.AuthJWT)
	if err != nil {
		return SystemUnderTest{}, err
	}
	locationConn, err := app.Dial(cfg.LocationGRPCAddr, cfg.AuthJWT)
	if err != nil {
		return SystemUnderTest{}, err
	}
	observer, err := newRabbitObserver(cfg.RabbitURL)
	if err != nil {
		return SystemUnderTest{}, err
	}
	return SystemUnderTest{
		Health:   healthChecker{ride: felo.NewRideServiceClient(rideConn), matching: felo.NewMatchingServiceClient(matchConn), wallet: felo.NewWalletServiceClient(walletConn), payment: felo.NewPaymentServiceClient(paymentConn), location: felo.NewLocationServiceClient(locationConn)},
		Ride:     rideClient{client: felo.NewRideServiceClient(rideConn)},
		Matching: matchingClient{client: felo.NewMatchingServiceClient(matchConn)},
		Wallet:   walletClient{client: felo.NewWalletServiceClient(walletConn)},
		Payment:  paymentClient{client: felo.NewPaymentServiceClient(paymentConn)},
		Location: locationClient{client: felo.NewLocationServiceClient(locationConn)},
		Events:   observer,
	}, nil
}

type healthChecker struct {
	ride     felo.RideServiceClient
	matching felo.MatchingServiceClient
	wallet   felo.WalletServiceClient
	payment  felo.PaymentServiceClient
	location felo.LocationServiceClient
}

func (h healthChecker) Check(ctx context.Context) error {
	checks := []func(context.Context) error{
		func(ctx context.Context) error { _, err := h.ride.Ping(ctx, &felo.PingRequest{}); return err },
		func(ctx context.Context) error { _, err := h.matching.Ping(ctx, &felo.PingRequest{}); return err },
		func(ctx context.Context) error { _, err := h.wallet.Ping(ctx, &felo.PingRequest{}); return err },
		func(ctx context.Context) error { _, err := h.payment.Ping(ctx, &felo.PingRequest{}); return err },
		func(ctx context.Context) error { _, err := h.location.Ping(ctx, &felo.PingRequest{}); return err },
	}
	for _, check := range checks {
		if err := check(ctx); err != nil {
			return err
		}
	}
	return nil
}

type rideClient struct{ client felo.RideServiceClient }

func (c rideClient) RequestRide(ctx context.Context, req RequestRideRequest) (Ride, error) {
	resp, err := c.client.RequestRide(ctx, &felo.RequestRideRequest{CustomerID: req.CustomerID, Pickup: felo.Coordinate(req.Pickup), Destination: felo.Coordinate(req.Destination), FareEstimate: req.FareEstimate})
	if err != nil {
		return Ride{}, err
	}
	return Ride{ID: resp.ID, State: resp.State, Driver: resp.Driver}, nil
}

func (c rideClient) StartRide(ctx context.Context, rideID string) (Ride, error) {
	resp, err := c.client.StartRide(ctx, &felo.RideByIDRequest{RideID: rideID})
	if err != nil {
		return Ride{}, err
	}
	return Ride{ID: resp.ID, State: resp.State, Driver: resp.Driver}, nil
}

func (c rideClient) CompleteRide(ctx context.Context, rideID string) (Ride, error) {
	resp, err := c.client.CompleteRide(ctx, &felo.RideByIDRequest{RideID: rideID})
	if err != nil {
		return Ride{}, err
	}
	return Ride{ID: resp.ID, State: resp.State, Driver: resp.Driver}, nil
}

func (c rideClient) GenerateNowQR(ctx context.Context, req GenerateNowQRRequest) (QRSession, error) {
	resp, err := c.client.GenerateNowQR(ctx, &felo.GenerateNowQRRequest{CustomerID: req.CustomerID, Destination: felo.Coordinate(req.Destination), FareEstimate: req.FareEstimate})
	if err != nil {
		return QRSession{}, err
	}
	expiresAt, _ := time.Parse(time.RFC3339, resp.ExpiresAt)
	return QRSession{TripID: resp.TripID, QRCode: resp.QRCode, ExpiresAt: expiresAt, DriverLock: resp.LockedDriver}, nil
}

func (c rideClient) ScanNowQR(ctx context.Context, qrCode string, driverID string) (QRSession, error) {
	resp, err := c.client.ScanNowQR(ctx, &felo.ScanNowQRRequest{QRCode: qrCode, DriverID: driverID})
	if err != nil {
		return QRSession{}, err
	}
	expiresAt, _ := time.Parse(time.RFC3339, resp.ExpiresAt)
	return QRSession{TripID: resp.TripID, QRCode: resp.QRCode, ExpiresAt: expiresAt, DriverLock: resp.LockedDriver}, nil
}

func (c rideClient) AcceptNowQR(ctx context.Context, tripID string, driverID string) (Ride, error) {
	resp, err := c.client.AcceptNowQR(ctx, &felo.AcceptNowQRRequest{TripID: tripID, DriverID: driverID})
	if err != nil {
		return Ride{}, err
	}
	return Ride{ID: resp.ID, State: resp.State, Driver: resp.Driver}, nil
}

type matchingClient struct{ client felo.MatchingServiceClient }

func (c matchingClient) GetAssignment(ctx context.Context, rideID string) (string, error) {
	resp, err := c.client.GetAssignment(ctx, &felo.AssignmentRequest{RideID: rideID})
	if err != nil {
		return "", err
	}
	return resp.DriverID, nil
}

type walletClient struct{ client felo.WalletServiceClient }

func (c walletClient) GetDriverBalance(ctx context.Context, driverID string) (WalletBalance, error) {
	resp, err := c.client.GetDriverBalance(ctx, &felo.BalanceRequest{OwnerID: driverID})
	if err != nil {
		return WalletBalance{}, err
	}
	return WalletBalance{OwnerID: resp.OwnerID, Balance: resp.Balance, Currency: resp.Currency}, nil
}

func (c walletClient) GetCustomerBalance(ctx context.Context, customerID string) (WalletBalance, error) {
	resp, err := c.client.GetCustomerBalance(ctx, &felo.BalanceRequest{OwnerID: customerID})
	if err != nil {
		return WalletBalance{}, err
	}
	return WalletBalance{OwnerID: resp.OwnerID, Balance: resp.Balance, Currency: resp.Currency}, nil
}

type paymentClient struct{ client felo.PaymentServiceClient }

func (c paymentClient) GetPaymentStatus(ctx context.Context, rideID string) (string, error) {
	resp, err := c.client.GetPaymentStatus(ctx, &felo.PaymentStatusRequest{RideID: rideID})
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}

type locationClient struct{ client felo.LocationServiceClient }

func (c locationClient) ReportDriverLocation(ctx context.Context, sample LocationSample) error {
	_, err := c.client.ReportDriverLocation(ctx, &felo.LocationSample{DriverID: sample.DriverID, Position: felo.Coordinate(sample.Position), RecordedAt: sample.RecordedAt.Format(time.RFC3339)})
	return err
}

func (c locationClient) GetDriverHistory(ctx context.Context, driverID string, from time.Time, to time.Time) ([]LocationSample, error) {
	resp, err := c.client.GetDriverHistory(ctx, &felo.LocationHistoryRequest{DriverID: driverID, From: from.Format(time.RFC3339), To: to.Format(time.RFC3339)})
	if err != nil {
		return nil, err
	}
	samples := make([]LocationSample, 0, len(resp.Samples))
	for _, sample := range resp.Samples {
		recordedAt, _ := time.Parse(time.RFC3339, sample.RecordedAt)
		samples = append(samples, LocationSample{DriverID: sample.DriverID, Position: Coordinate(sample.Position), RecordedAt: recordedAt})
	}
	return samples, nil
}

type rabbitObserver struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
	msgs <-chan amqp091.Delivery
}

func newRabbitObserver(url string) (*rabbitObserver, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	queue, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return nil, err
	}
	if err := ch.QueueBind(queue.Name, "#", "felo.events", false, nil); err != nil {
		return nil, err
	}
	msgs, err := ch.Consume(queue.Name, "", true, true, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &rabbitObserver{conn: conn, ch: ch, msgs: msgs}, nil
}

func (o *rabbitObserver) WaitForEvent(ctx context.Context, name string, key string) (Event, error) {
	for {
		select {
		case <-ctx.Done():
			return Event{}, ctx.Err()
		case msg, ok := <-o.msgs:
			if !ok {
				return Event{}, fmt.Errorf("event stream closed")
			}
			var event felo.EventEnvelope
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				continue
			}
			if event.Name == name && event.Key == key {
				return Event{Name: event.Name, Key: event.Key, Payload: event.Payload}, nil
			}
		}
	}
}
