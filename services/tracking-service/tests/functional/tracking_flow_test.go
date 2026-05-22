//go:build functional

package functional_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/felo/felo-backend/services/tracking-service/internal/domain"
	"github.com/felo/felo-backend/services/tracking-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestTrackingFunctional_StartTracking_PersistsSession(t *testing.T) {
	ctx := context.Background()
	db := openTrackingPG(t, getenv("FELO_TRACKING_PG_DSN", "postgres://felo:felo@127.0.0.1:54340/tracking_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initTrackingTables(t, db)

	svc := service.NewTrackingService(
		&pgSessionRepo{db: db},
		&pgRecordRepo{db: db},
		&noopTrackingPublisher{},
		fixedClock{now: time.Date(2026, 5, 5, 10, 0, 0, 0, time.UTC)},
		&seqIDGen{next: 1, prefix: "track-ft-"},
	)

	session, err := svc.StartTracking(ctx, service.StartTrackingInput{
		ShipmentID: "ship-ft-001",
		DriverID:   "driver-ft-001",
	})
	if err != nil {
		t.Fatalf("StartTracking() error = %v", err)
	}

	var status string
	var shipmentID, driverID string
	if err := db.QueryRow(ctx,
		"select session_id, shipment_id, driver_id, status from tracking_sessions where session_id=$1",
		session.ID,
	).Scan(&session.ID, &shipmentID, &driverID, &status); err != nil {
		t.Fatalf("query persisted session: %v", err)
	}
	if status != string(domain.StatusActive) {
		t.Fatalf("persisted status = %s, want %s", status, domain.StatusActive)
	}
	if shipmentID != "ship-ft-001" {
		t.Fatalf("persisted shipment_id = %s, want ship-ft-001", shipmentID)
	}
	if driverID != "driver-ft-001" {
		t.Fatalf("persisted driver_id = %s, want driver-ft-001", driverID)
	}
}

func TestTrackingFunctional_RecordLocation_PersistsRecord(t *testing.T) {
	ctx := context.Background()
	db := openTrackingPG(t, getenv("FELO_TRACKING_PG_DSN", "postgres://felo:felo@127.0.0.1:54340/tracking_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initTrackingTables(t, db)

	now := time.Date(2026, 5, 5, 10, 0, 0, 0, time.UTC)

	_, _ = db.Exec(ctx,
		`insert into tracking_sessions (session_id, shipment_id, driver_id, status, started_at, updated_at)
		 values ($1,$2,$3,$4,$5,$6)`,
		"track-ft-002", "ship-ft-002", "driver-ft-002",
		string(domain.StatusActive), now, now,
	)

	svc := service.NewTrackingService(
		&pgSessionRepo{db: db},
		&pgRecordRepo{db: db},
		&noopTrackingPublisher{},
		fixedClock{now: now},
		&seqIDGen{next: 1, prefix: "rec-ft-"},
	)

	record, err := svc.RecordLocation(ctx, service.RecordLocationInput{
		SessionID:  "track-ft-002",
		Coordinate: domain.Coordinate{Latitude: -6.2, Longitude: 106.8},
		Speed:      40.5,
		Heading:    180.0,
	})
	if err != nil {
		t.Fatalf("RecordLocation() error = %v", err)
	}

	var lat, lng, speed, heading float64
	var sessionID string
	if err := db.QueryRow(ctx,
		"select session_id, lat, lng, speed, heading from tracking_records where record_id=$1",
		record.ID,
	).Scan(&sessionID, &lat, &lng, &speed, &heading); err != nil {
		t.Fatalf("query persisted record: %v", err)
	}
	if sessionID != "track-ft-002" {
		t.Fatalf("persisted session_id = %s, want track-ft-002", sessionID)
	}
	if lat != -6.2 || lng != 106.8 {
		t.Fatalf("persisted coordinate = (%f,%f), want (-6.2,106.8)", lat, lng)
	}
}

func TestTrackingFunctional_StopTracking_UpdatesSession(t *testing.T) {
	ctx := context.Background()
	db := openTrackingPG(t, getenv("FELO_TRACKING_PG_DSN", "postgres://felo:felo@127.0.0.1:54340/tracking_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initTrackingTables(t, db)

	now := time.Date(2026, 5, 5, 10, 0, 0, 0, time.UTC)

	_, _ = db.Exec(ctx,
		`insert into tracking_sessions (session_id, shipment_id, driver_id, status, started_at, updated_at)
		 values ($1,$2,$3,$4,$5,$6)`,
		"track-ft-003", "ship-ft-003", "driver-ft-003",
		string(domain.StatusActive), now, now,
	)

	svc := service.NewTrackingService(
		&pgSessionRepo{db: db},
		&pgRecordRepo{db: db},
		&noopTrackingPublisher{},
		fixedClock{now: now.Add(1 * time.Hour)},
		&seqIDGen{},
	)

	session, err := svc.StopTracking(ctx, "track-ft-003")
	if err != nil {
		t.Fatalf("StopTracking() error = %v", err)
	}

	var status string
	var endedAt *time.Time
	if err := db.QueryRow(ctx,
		"select status, ended_at from tracking_sessions where session_id=$1",
		session.ID,
	).Scan(&status, &endedAt); err != nil {
		t.Fatalf("query persisted session: %v", err)
	}
	if status != string(domain.StatusCompleted) {
		t.Fatalf("persisted status = %s, want %s", status, domain.StatusCompleted)
	}
	if endedAt == nil {
		t.Fatal("persisted ended_at is nil, want non-nil")
	}
}

func TestTrackingFunctional_GetTrackingHistory_ReturnsRecords(t *testing.T) {
	ctx := context.Background()
	db := openTrackingPG(t, getenv("FELO_TRACKING_PG_DSN", "postgres://felo:felo@127.0.0.1:54340/tracking_db?sslmode=disable"))
	t.Cleanup(func() { db.Close() })

	initTrackingTables(t, db)

	now := time.Date(2026, 5, 5, 10, 0, 0, 0, time.UTC)
	ended := now.Add(2 * time.Hour)

	_, _ = db.Exec(ctx,
		`insert into tracking_sessions (session_id, shipment_id, driver_id, status, started_at, updated_at, ended_at)
		 values ($1,$2,$3,$4,$5,$6,$7)`,
		"track-ft-004", "ship-ft-004", "driver-ft-004",
		string(domain.StatusCompleted), now, ended, &ended,
	)

	for i := 0; i < 3; i++ {
		_, _ = db.Exec(ctx,
			`insert into tracking_records (record_id, session_id, lat, lng, speed, heading, recorded_at)
			 values ($1,$2,$3,$4,$5,$6,$7)`,
			"rec-ft-004-", i, "track-ft-004", -6.2, 106.8, float64(30+i*5), 180.0, now.Add(time.Duration(i)*time.Minute),
		)
	}

	svc := service.NewTrackingService(
		&pgSessionRepo{db: db},
		&pgRecordRepo{db: db},
		&noopTrackingPublisher{},
		fixedClock{now: now},
		&seqIDGen{},
	)

	records, err := svc.GetTrackingHistory(ctx, "track-ft-004")
	if err != nil {
		t.Fatalf("GetTrackingHistory() error = %v", err)
	}
	if len(records) != 3 {
		t.Fatalf("len(records) = %d, want 3", len(records))
	}
}

type pgSessionRepo struct{ db *pgxpool.Pool }

func (r *pgSessionRepo) Save(ctx context.Context, session domain.TrackingSession) error {
	_, err := r.db.Exec(ctx,
		`insert into tracking_sessions (session_id, shipment_id, driver_id, status, started_at, updated_at, ended_at)
		 values ($1,$2,$3,$4,$5,$6,$7)
		 on conflict (session_id) do update set
		   status=excluded.status,
		   updated_at=excluded.updated_at,
		   ended_at=excluded.ended_at`,
		session.ID, session.ShipmentID, session.DriverID, string(session.Status),
		session.StartedAt, session.UpdatedAt, session.EndedAt)
	return err
}

func (r *pgSessionRepo) GetByID(ctx context.Context, sessionID string) (domain.TrackingSession, error) {
	var session domain.TrackingSession
	var status string
	err := r.db.QueryRow(ctx,
		`select session_id, shipment_id, driver_id, status, started_at, updated_at, ended_at
		 from tracking_sessions where session_id=$1`, sessionID,
	).Scan(&session.ID, &session.ShipmentID, &session.DriverID, &status,
		&session.StartedAt, &session.UpdatedAt, &session.EndedAt)
	if err != nil {
		return domain.TrackingSession{}, err
	}
	session.Status = domain.TrackingStatus(status)
	return session, nil
}

type pgRecordRepo struct{ db *pgxpool.Pool }

func (r *pgRecordRepo) Save(ctx context.Context, record domain.TrackingRecord) error {
	_, err := r.db.Exec(ctx,
		`insert into tracking_records (record_id, session_id, lat, lng, speed, heading, recorded_at)
		 values ($1,$2,$3,$4,$5,$6,$7)`,
		record.ID, record.SessionID, record.Coordinate.Latitude, record.Coordinate.Longitude,
		record.Speed, record.Heading, record.RecordedAt)
	return err
}

func (r *pgRecordRepo) ListBySession(ctx context.Context, sessionID string) ([]domain.TrackingRecord, error) {
	rows, err := r.db.Query(ctx,
		`select record_id, session_id, lat, lng, speed, heading, recorded_at
		 from tracking_records where session_id=$1 order by recorded_at`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.TrackingRecord
	for rows.Next() {
		var rec domain.TrackingRecord
		if err := rows.Scan(&rec.ID, &rec.SessionID, &rec.Coordinate.Latitude,
			&rec.Coordinate.Longitude, &rec.Speed, &rec.Heading, &rec.RecordedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

type noopTrackingPublisher struct{}

func (noopTrackingPublisher) Publish(_ context.Context, _ domain.Event) error { return nil }

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

type seqIDGen struct {
	next   int
	prefix string
}

func (g *seqIDGen) NewID() string {
	if g.next == 0 {
		g.next = 1
	}
	id := fmt.Sprintf("%s%03d", g.prefix, g.next)
	g.next++
	return id
}

func initTrackingTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, `create table if not exists tracking_sessions (
		session_id text primary key,
		shipment_id text not null,
		driver_id text not null,
		status text not null,
		started_at timestamptz not null,
		updated_at timestamptz not null,
		ended_at timestamptz
	)`)
	if err != nil {
		t.Fatalf("init tracking_sessions table: %v", err)
	}
	_, err = db.Exec(ctx, `create table if not exists tracking_records (
		record_id text primary key,
		session_id text not null,
		lat double precision not null,
		lng double precision not null,
		speed double precision not null default 0,
		heading double precision not null default 0,
		recorded_at timestamptz not null
	)`)
	if err != nil {
		t.Fatalf("init tracking_records table: %v", err)
	}
}

func openTrackingPG(t *testing.T, dsn string) *pgxpool.Pool {
	t.Helper()
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("pgxpool.New() error = %v", err)
	}
	return db
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
