package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/felo/felo-backend/api/felo"
	"github.com/felo/felo-backend/api/jsoncodec"
	"github.com/felo/felo-backend/internal/demo/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const exchangeName = "felo.events"

type App struct {
	cfg        config.Config
	rideDB     *pgxpool.Pool
	matchDB    *pgxpool.Pool
	walletDB   *pgxpool.Pool
	paymentDB  *pgxpool.Pool
	locationDB *pgxpool.Pool
	redis      *redis.Client
	rabbitConn *amqp091.Connection
	rabbitCh   *amqp091.Channel
	servers    []*grpc.Server
	closers    []func()
	wg         sync.WaitGroup
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	encoding.RegisterCodec(jsoncodec.Codec{})

	rideDB, err := pgxpool.New(ctx, cfg.RidePostgresDSN)
	if err != nil {
		return nil, err
	}
	matchDB, err := pgxpool.New(ctx, cfg.MatchPostgresDSN)
	if err != nil {
		return nil, err
	}
	walletDB, err := pgxpool.New(ctx, cfg.WalletPostgresDSN)
	if err != nil {
		return nil, err
	}
	paymentDB, err := pgxpool.New(ctx, cfg.PaymentPostgresDSN)
	if err != nil {
		return nil, err
	}
	locationDB, err := pgxpool.New(ctx, cfg.LocationPostgresDSN)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	rabbitConn, err := amqp091.Dial(cfg.RabbitURL)
	if err != nil {
		return nil, err
	}
	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		return nil, err
	}
	if err := rabbitCh.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}

	app := &App{
		cfg:        cfg,
		rideDB:     rideDB,
		matchDB:    matchDB,
		walletDB:   walletDB,
		paymentDB:  paymentDB,
		locationDB: locationDB,
		redis:      rdb,
		rabbitConn: rabbitConn,
		rabbitCh:   rabbitCh,
	}
	if err := app.initSchema(ctx); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.startConsumers(ctx); err != nil {
		return err
	}
	if err := a.startGRPCServers(); err != nil {
		return err
	}

	<-ctx.Done()

	for _, server := range a.servers {
		server.GracefulStop()
	}
	for _, closeFn := range a.closers {
		closeFn()
	}
	a.wg.Wait()

	a.rabbitCh.Close()
	a.rabbitConn.Close()
	a.redis.Close()
	a.rideDB.Close()
	a.matchDB.Close()
	a.walletDB.Close()
	a.paymentDB.Close()
	a.locationDB.Close()
	return nil
}

func (a *App) startGRPCServers() error {
	type endpoint struct {
		addr     string
		register func(*grpc.Server)
	}

	interceptor := a.authInterceptor()
	codec := grpc.ForceServerCodec(jsoncodec.Codec{})
	endpoints := []endpoint{
		{addr: a.cfg.RideAddr, register: func(s *grpc.Server) { felo.RegisterRideServiceServer(s, &rideService{app: a}) }},
		{addr: a.cfg.MatchingAddr, register: func(s *grpc.Server) { felo.RegisterMatchingServiceServer(s, &matchingService{app: a}) }},
		{addr: a.cfg.WalletAddr, register: func(s *grpc.Server) { felo.RegisterWalletServiceServer(s, &walletService{app: a}) }},
		{addr: a.cfg.PaymentAddr, register: func(s *grpc.Server) { felo.RegisterPaymentServiceServer(s, &paymentService{app: a}) }},
		{addr: a.cfg.LocationAddr, register: func(s *grpc.Server) { felo.RegisterLocationServiceServer(s, &locationService{app: a}) }},
	}

	for _, endpoint := range endpoints {
		lis, err := net.Listen("tcp", endpoint.addr)
		if err != nil {
			return err
		}
		server := grpc.NewServer(codec, grpc.UnaryInterceptor(interceptor))
		endpoint.register(server)
		a.servers = append(a.servers, server)
		a.closers = append(a.closers, func() { _ = lis.Close() })
		a.wg.Add(1)
		go func(l net.Listener, s *grpc.Server) {
			defer a.wg.Done()
			if err := s.Serve(l); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("grpc server stopped with error: %v", err)
			}
		}(lis, server)
	}

	return nil
}

func (a *App) authInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		values := md.Get("authorization")
		if len(values) == 0 || values[0] != "Bearer "+a.cfg.AuthToken {
			return nil, status.Error(codes.Unauthenticated, "invalid auth token")
		}
		return handler(ctx, req)
	}
}

func (a *App) initSchema(ctx context.Context) error {
	queries := []struct {
		db  *pgxpool.Pool
		sql string
	}{
		{a.rideDB, `create table if not exists rides (
			id text primary key,
			customer_id text not null,
			driver_id text not null default '',
			pickup_lat double precision not null,
			pickup_lng double precision not null,
			dest_lat double precision not null,
			dest_lng double precision not null,
			fare bigint not null,
			state text not null,
			qr_code text not null default '',
			qr_expires_at timestamptz,
			qr_locked_driver text not null default '',
			created_at timestamptz not null,
			updated_at timestamptz not null
		)`},
		{a.matchDB, `create table if not exists drivers (
			id text primary key,
			status text not null,
			lat double precision not null,
			lng double precision not null
		);
		create table if not exists assignments (
			ride_id text primary key,
			driver_id text not null,
			assigned_at timestamptz not null
		)`},
		{a.walletDB, `create table if not exists wallets (
			owner_id text primary key,
			owner_type text not null,
			balance bigint not null,
			currency text not null,
			updated_at timestamptz not null
		);
		create table if not exists wallet_ledger (
			reference text primary key,
			owner_id text not null,
			delta bigint not null,
			reason text not null,
			created_at timestamptz not null
		)`},
		{a.paymentDB, `create table if not exists payments (
			event_id text primary key,
			ride_id text unique not null,
			customer_id text not null,
			amount bigint not null,
			currency text not null,
			status text not null,
			created_at timestamptz not null
		)`},
		{a.locationDB, `create table if not exists driver_locations (
			id bigserial primary key,
			driver_id text not null,
			lat double precision not null,
			lng double precision not null,
			recorded_at timestamptz not null
		)`},
	}

	for _, query := range queries {
		if _, err := query.db.Exec(ctx, query.sql); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) startConsumers(ctx context.Context) error {
	if err := a.consume(ctx, "matching-service", "ride.created.v1", a.handleRideCreated); err != nil {
		return err
	}
	if err := a.consume(ctx, "payment-service", "ride.completed.v1", a.handleRideCompleted); err != nil {
		return err
	}
	if err := a.consume(ctx, "wallet-service", "payment.completed.v1", a.handlePaymentCompleted); err != nil {
		return err
	}
	return nil
}

func (a *App) consume(ctx context.Context, queueName string, routingKey string, handler func(context.Context, felo.EventEnvelope) error) error {
	ch, err := a.rabbitConn.Channel()
	if err != nil {
		return err
	}
	if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind(queueName, routingKey, exchangeName, false, nil); err != nil {
		return err
	}
	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	a.closers = append(a.closers, func() { _ = ch.Close() })
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				var event felo.EventEnvelope
				if err := json.Unmarshal(msg.Body, &event); err != nil {
					_ = msg.Nack(false, false)
					continue
				}
				if err := handler(ctx, event); err != nil {
					log.Printf("event handler error for %s: %v", event.Name, err)
					_ = msg.Nack(false, false)
					continue
				}
				_ = msg.Ack(false)
			}
		}
	}()
	return nil
}

func (a *App) publish(ctx context.Context, name string, key string, payload map[string]string) error {
	body, err := json.Marshal(felo.EventEnvelope{Name: name, Key: key, Payload: payload})
	if err != nil {
		return err
	}
	return a.rabbitCh.PublishWithContext(ctx, exchangeName, name, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
		Timestamp:   time.Now(),
	})
}

func (a *App) handleRideCreated(ctx context.Context, event felo.EventEnvelope) error {
	rideID := event.Key
	var pickupLat, pickupLng float64
	if _, err := fmt.Sscanf(event.Payload["pickup"], "%f,%f", &pickupLat, &pickupLng); err != nil {
		return err
	}

	rows, err := a.matchDB.Query(ctx, `select id, lat, lng from drivers where status='available'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	bestDriver := ""
	bestDistance := math.MaxFloat64
	for rows.Next() {
		var driverID string
		var lat, lng float64
		if err := rows.Scan(&driverID, &lat, &lng); err != nil {
			return err
		}
		distance := roughDistanceKM(pickupLat, pickupLng, lat, lng)
		if distance < bestDistance && distance <= 5.5 {
			bestDistance = distance
			bestDriver = driverID
		}
	}
	if bestDriver == "" {
		return nil
	}

	if _, err := a.matchDB.Exec(ctx, `insert into assignments (ride_id, driver_id, assigned_at) values ($1, $2, $3)
		on conflict (ride_id) do update set driver_id=excluded.driver_id, assigned_at=excluded.assigned_at`, rideID, bestDriver, time.Now()); err != nil {
		return err
	}
	if _, err := a.rideDB.Exec(ctx, `update rides set driver_id=$2 where id=$1`, rideID, bestDriver); err != nil {
		return err
	}
	return a.publish(ctx, "driver.matched.v1", rideID, map[string]string{"ride_id": rideID, "driver_id": bestDriver})
}

func (a *App) handleRideCompleted(ctx context.Context, event felo.EventEnvelope) error {
	customerID := event.Payload["customer_id"]
	rideID := event.Key
	driverID := event.Payload["driver_id"]
	amount := parseInt64(event.Payload["amount"])
	currency := event.Payload["currency"]
	if currency == "" {
		currency = "IDR"
	}

	balance, err := a.getWalletBalance(ctx, customerID)
	if err != nil {
		return err
	}
	if balance < amount {
		if _, err := a.paymentDB.Exec(ctx, `insert into payments (event_id, ride_id, customer_id, amount, currency, status, created_at)
			values ($1,$2,$3,$4,$5,$6,$7) on conflict (event_id) do nothing`, rideID, rideID, customerID, amount, currency, "failed", time.Now()); err != nil {
			return err
		}
		return a.publish(ctx, "payment.failed.v1", rideID, map[string]string{"ride_id": rideID, "driver_id": driverID})
	}

	if err := a.applyWalletDelta(ctx, customerID, -amount, rideID+"-debit", "ride_payment"); err != nil {
		return err
	}
	if _, err := a.paymentDB.Exec(ctx, `insert into payments (event_id, ride_id, customer_id, amount, currency, status, created_at)
		values ($1,$2,$3,$4,$5,$6,$7) on conflict (event_id) do nothing`, rideID, rideID, customerID, amount, currency, "completed", time.Now()); err != nil {
		return err
	}
	return a.publish(ctx, "payment.completed.v1", rideID, map[string]string{
		"ride_id":    rideID,
		"driver_id":  driverID,
		"customer_id": customerID,
		"amount":     fmt.Sprintf("%d", amount),
		"currency":   currency,
	})
}

func (a *App) handlePaymentCompleted(ctx context.Context, event felo.EventEnvelope) error {
	driverID := event.Payload["driver_id"]
	amount := parseInt64(event.Payload["amount"])
	return a.applyWalletDelta(ctx, driverID, amount, event.Key+"-credit", "driver_settlement")
}

func (a *App) getWalletBalance(ctx context.Context, ownerID string) (int64, error) {
	var balance int64
	err := a.walletDB.QueryRow(ctx, `select balance from wallets where owner_id=$1`, ownerID).Scan(&balance)
	return balance, err
}

func (a *App) applyWalletDelta(ctx context.Context, ownerID string, delta int64, reference string, reason string) error {
	tx, err := a.walletDB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var exists string
	err = tx.QueryRow(ctx, `select reference from wallet_ledger where reference=$1`, reference).Scan(&exists)
	if err == nil {
		return tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	var balance int64
	if err := tx.QueryRow(ctx, `select balance from wallets where owner_id=$1 for update`, ownerID).Scan(&balance); err != nil {
		return err
	}
	next := balance + delta
	if next < 0 {
		return status.Error(codes.FailedPrecondition, "insufficient balance")
	}

	if _, err := tx.Exec(ctx, `update wallets set balance=$2, updated_at=$3 where owner_id=$1`, ownerID, next, time.Now()); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `insert into wallet_ledger (reference, owner_id, delta, reason, created_at) values ($1,$2,$3,$4,$5)`,
		reference, ownerID, delta, reason, time.Now()); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func roughDistanceKM(lat1, lng1, lat2, lng2 float64) float64 {
	dlat := (lat1 - lat2) * 111
	dlng := (lng1 - lng2) * 111
	return math.Sqrt(dlat*dlat + dlng*dlng)
}

func parseInt64(value string) int64 {
	var parsed int64
	fmt.Sscanf(value, "%d", &parsed)
	return parsed
}

type rideService struct{ app *App }

func (s *rideService) Ping(context.Context, *felo.PingRequest) (*felo.PingResponse, error) {
	return &felo.PingResponse{Status: "ok"}, nil
}

func (s *rideService) RequestRide(ctx context.Context, req *felo.RequestRideRequest) (*felo.Ride, error) {
	rideID := fmt.Sprintf("ride-%d", time.Now().UnixNano())
	now := time.Now().UTC()
	if _, err := s.app.rideDB.Exec(ctx, `insert into rides (id, customer_id, pickup_lat, pickup_lng, dest_lat, dest_lng, fare, state, created_at, updated_at)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, rideID, req.CustomerID, req.Pickup.Latitude, req.Pickup.Longitude, req.Destination.Latitude, req.Destination.Longitude, req.FareEstimate, "matching", now, now); err != nil {
		return nil, err
	}
	if err := s.app.publish(ctx, "ride.created.v1", rideID, map[string]string{
		"ride_id":    rideID,
		"customer_id": req.CustomerID,
		"pickup":     fmt.Sprintf("%f,%f", req.Pickup.Latitude, req.Pickup.Longitude),
		"amount":     fmt.Sprintf("%d", req.FareEstimate),
		"currency":   "IDR",
	}); err != nil {
		return nil, err
	}
	return &felo.Ride{ID: rideID, State: "matching"}, nil
}

func (s *rideService) StartRide(ctx context.Context, req *felo.RideByIDRequest) (*felo.Ride, error) {
	var state string
	if err := s.app.rideDB.QueryRow(ctx, `select state from rides where id=$1`, req.RideID).Scan(&state); err != nil {
		return nil, err
	}
	if state != "matching" && state != "arrived" && state != "on_ride" {
		return nil, status.Error(codes.FailedPrecondition, "ride cannot start")
	}
	if _, err := s.app.rideDB.Exec(ctx, `update rides set state='on_ride', updated_at=$2 where id=$1`, req.RideID, time.Now().UTC()); err != nil {
		return nil, err
	}
	return &felo.Ride{ID: req.RideID, State: "on_ride"}, nil
}

func (s *rideService) CompleteRide(ctx context.Context, req *felo.RideByIDRequest) (*felo.Ride, error) {
	var customerID, driverID string
	var amount int64
	if err := s.app.rideDB.QueryRow(ctx, `select customer_id, driver_id, fare from rides where id=$1`, req.RideID).Scan(&customerID, &driverID, &amount); err != nil {
		return nil, err
	}
	if _, err := s.app.rideDB.Exec(ctx, `update rides set state='completed', updated_at=$2 where id=$1`, req.RideID, time.Now().UTC()); err != nil {
		return nil, err
	}
	if err := s.app.publish(ctx, "ride.completed.v1", req.RideID, map[string]string{
		"ride_id":    req.RideID,
		"customer_id": customerID,
		"driver_id":  driverID,
		"amount":     fmt.Sprintf("%d", amount),
		"currency":   "IDR",
	}); err != nil {
		return nil, err
	}
	return &felo.Ride{ID: req.RideID, State: "completed", Driver: driverID}, nil
}

func (s *rideService) GenerateNowQR(ctx context.Context, req *felo.GenerateNowQRRequest) (*felo.QRSession, error) {
	rideID := fmt.Sprintf("ride-%d", time.Now().UnixNano())
	qrCode := fmt.Sprintf("qr-%d", time.Now().UnixNano())
	now := time.Now().UTC()
	expiry := now.Add(10 * time.Minute)
	if _, err := s.app.rideDB.Exec(ctx, `insert into rides (id, customer_id, pickup_lat, pickup_lng, dest_lat, dest_lng, fare, state, qr_code, qr_expires_at, created_at, updated_at)
		values ($1,$2,0,0,$3,$4,$5,$6,$7,$8,$9,$10)`, rideID, req.CustomerID, req.Destination.Latitude, req.Destination.Longitude, req.FareEstimate, "qr_generated", qrCode, expiry, now, now); err != nil {
		return nil, err
	}
	return &felo.QRSession{TripID: rideID, QRCode: qrCode, ExpiresAt: expiry.Format(time.RFC3339)}, nil
}

func (s *rideService) ScanNowQR(ctx context.Context, req *felo.ScanNowQRRequest) (*felo.QRSession, error) {
	var rideID, lockedDriver string
	var expiry time.Time
	if err := s.app.rideDB.QueryRow(ctx, `select id, qr_expires_at, qr_locked_driver from rides where qr_code=$1`, req.QRCode).Scan(&rideID, &expiry, &lockedDriver); err != nil {
		return nil, err
	}
	if time.Now().UTC().After(expiry) {
		return nil, status.Error(codes.FailedPrecondition, "qr expired")
	}
	if lockedDriver != "" && lockedDriver != req.DriverID {
		return nil, status.Error(codes.FailedPrecondition, "qr locked by another driver")
	}
	if _, err := s.app.rideDB.Exec(ctx, `update rides set qr_locked_driver=$2, updated_at=$3 where id=$1`, rideID, req.DriverID, time.Now().UTC()); err != nil {
		return nil, err
	}
	return &felo.QRSession{TripID: rideID, QRCode: req.QRCode, ExpiresAt: expiry.Format(time.RFC3339), LockedDriver: req.DriverID}, nil
}

func (s *rideService) AcceptNowQR(ctx context.Context, req *felo.AcceptNowQRRequest) (*felo.Ride, error) {
	var qrLockedDriver string
	if err := s.app.rideDB.QueryRow(ctx, `select qr_locked_driver from rides where id=$1`, req.TripID).Scan(&qrLockedDriver); err != nil {
		return nil, err
	}
	if qrLockedDriver != req.DriverID {
		return nil, status.Error(codes.PermissionDenied, "driver does not hold qr lock")
	}
	if _, err := s.app.rideDB.Exec(ctx, `update rides set state='on_ride', driver_id=$2, updated_at=$3 where id=$1`, req.TripID, req.DriverID, time.Now().UTC()); err != nil {
		return nil, err
	}
	return &felo.Ride{ID: req.TripID, State: "on_ride", Driver: req.DriverID}, nil
}

type matchingService struct{ app *App }

func (s *matchingService) Ping(context.Context, *felo.PingRequest) (*felo.PingResponse, error) {
	return &felo.PingResponse{Status: "ok"}, nil
}

func (s *matchingService) GetAssignment(ctx context.Context, req *felo.AssignmentRequest) (*felo.AssignmentResponse, error) {
	var driverID string
	err := s.app.matchDB.QueryRow(ctx, `select driver_id from assignments where ride_id=$1`, req.RideID).Scan(&driverID)
	if err != nil {
		return &felo.AssignmentResponse{DriverID: ""}, nil
	}
	return &felo.AssignmentResponse{DriverID: driverID}, nil
}

type walletService struct{ app *App }

func (s *walletService) Ping(context.Context, *felo.PingRequest) (*felo.PingResponse, error) {
	return &felo.PingResponse{Status: "ok"}, nil
}

func (s *walletService) GetDriverBalance(ctx context.Context, req *felo.BalanceRequest) (*felo.WalletBalance, error) {
	return s.getBalance(ctx, req.OwnerID)
}

func (s *walletService) GetCustomerBalance(ctx context.Context, req *felo.BalanceRequest) (*felo.WalletBalance, error) {
	return s.getBalance(ctx, req.OwnerID)
}

func (s *walletService) getBalance(ctx context.Context, ownerID string) (*felo.WalletBalance, error) {
	var balance int64
	var currency string
	if err := s.app.walletDB.QueryRow(ctx, `select balance, currency from wallets where owner_id=$1`, ownerID).Scan(&balance, &currency); err != nil {
		return nil, err
	}
	return &felo.WalletBalance{OwnerID: ownerID, Balance: balance, Currency: currency}, nil
}

type paymentService struct{ app *App }

func (s *paymentService) Ping(context.Context, *felo.PingRequest) (*felo.PingResponse, error) {
	return &felo.PingResponse{Status: "ok"}, nil
}

func (s *paymentService) GetPaymentStatus(ctx context.Context, req *felo.PaymentStatusRequest) (*felo.PaymentStatusResponse, error) {
	var statusValue string
	if err := s.app.paymentDB.QueryRow(ctx, `select status from payments where ride_id=$1`, req.RideID).Scan(&statusValue); err != nil {
		return &felo.PaymentStatusResponse{Status: ""}, nil
	}
	return &felo.PaymentStatusResponse{Status: statusValue}, nil
}

type locationService struct{ app *App }

func (s *locationService) Ping(context.Context, *felo.PingRequest) (*felo.PingResponse, error) {
	return &felo.PingResponse{Status: "ok"}, nil
}

func (s *locationService) ReportDriverLocation(ctx context.Context, req *felo.LocationSample) (*felo.PingResponse, error) {
	recordedAt, err := time.Parse(time.RFC3339, req.RecordedAt)
	if err != nil {
		return nil, err
	}
	if _, err := s.app.locationDB.Exec(ctx, `insert into driver_locations (driver_id, lat, lng, recorded_at) values ($1,$2,$3,$4)`,
		req.DriverID, req.Position.Latitude, req.Position.Longitude, recordedAt); err != nil {
		return nil, err
	}
	payload, _ := json.Marshal(req)
	if err := s.app.redis.Set(ctx, "latest:"+req.DriverID, string(payload), 0).Err(); err != nil {
		return nil, err
	}
	return &felo.PingResponse{Status: "ok"}, nil
}

func (s *locationService) GetDriverHistory(ctx context.Context, req *felo.LocationHistoryRequest) (*felo.LocationHistoryResponse, error) {
	from, err := time.Parse(time.RFC3339, req.From)
	if err != nil {
		return nil, err
	}
	to, err := time.Parse(time.RFC3339, req.To)
	if err != nil {
		return nil, err
	}
	rows, err := s.app.locationDB.Query(ctx, `select driver_id, lat, lng, recorded_at from driver_locations where driver_id=$1 and recorded_at between $2 and $3 order by recorded_at`, req.DriverID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var samples []felo.LocationSample
	for rows.Next() {
		var sample felo.LocationSample
		var recordedAt time.Time
		if err := rows.Scan(&sample.DriverID, &sample.Position.Latitude, &sample.Position.Longitude, &recordedAt); err != nil {
			return nil, err
		}
		sample.RecordedAt = recordedAt.Format(time.RFC3339)
		samples = append(samples, sample)
	}
	return &felo.LocationHistoryResponse{Samples: samples}, nil
}

func Dial(addr string, token string) (*grpc.ClientConn, error) {
	return grpc.DialContext(context.Background(), addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(jsoncodec.Codec{})),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
}

func (a *App) Seed(ctx context.Context, root string) error {
	type customerSeed struct {
		ID            string `json:"id"`
		WalletBalance int64  `json:"wallet_balance"`
		Currency      string `json:"currency"`
		Status        string `json:"status"`
	}
	type driverSeed struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
	}
	type locationSeed struct {
		DriverID string `json:"driver_id"`
		Samples  []struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			RecordedAt string `json:"recorded_at"`
		} `json:"samples"`
	}

	var customers []customerSeed
	if err := readJSON(filepath.Join(root, "customers.json"), &customers); err != nil {
		return err
	}
	for _, customer := range customers {
		if _, err := a.walletDB.Exec(ctx, `insert into wallets (owner_id, owner_type, balance, currency, updated_at)
			values ($1,$2,$3,$4,$5)
			on conflict (owner_id) do update set balance=excluded.balance, currency=excluded.currency, updated_at=excluded.updated_at`,
			customer.ID, "customer", customer.WalletBalance, customer.Currency, time.Now().UTC()); err != nil {
			return err
		}
	}

	var drivers []driverSeed
	if err := readJSON(filepath.Join(root, "drivers.json"), &drivers); err != nil {
		return err
	}
	for _, driver := range drivers {
		if _, err := a.matchDB.Exec(ctx, `insert into drivers (id, status, lat, lng)
			values ($1,$2,$3,$4)
			on conflict (id) do update set status=excluded.status, lat=excluded.lat, lng=excluded.lng`,
			driver.ID, driver.Status, driver.Location.Latitude, driver.Location.Longitude); err != nil {
			return err
		}
		if _, err := a.walletDB.Exec(ctx, `insert into wallets (owner_id, owner_type, balance, currency, updated_at)
			values ($1,$2,$3,$4,$5)
			on conflict (owner_id) do nothing`, driver.ID, "driver", int64(0), "IDR", time.Now().UTC()); err != nil {
			return err
		}
	}

	var locations []locationSeed
	if err := readJSON(filepath.Join(root, "locations.json"), &locations); err != nil {
		return err
	}
	for _, item := range locations {
		for _, sample := range item.Samples {
			recordedAt, err := time.Parse(time.RFC3339, sample.RecordedAt)
			if err != nil {
				return err
			}
			if _, err := a.locationDB.Exec(ctx, `insert into driver_locations (driver_id, lat, lng, recorded_at) values ($1,$2,$3,$4)`,
				item.DriverID, sample.Latitude, sample.Longitude, recordedAt); err != nil {
				return err
			}
		}
	}
	return nil
}

func readJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
