CREATE SCHEMA IF NOT EXISTS warehouse;

SET search_path TO warehouse;

CREATE TABLE IF NOT EXISTS dim_date (
    date_key INTEGER PRIMARY KEY,
    full_date DATE NOT NULL UNIQUE,
    day INTEGER NOT NULL,
    day_name VARCHAR(16) NOT NULL,
    week_of_year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    month_name VARCHAR(16) NOT NULL,
    quarter INTEGER NOT NULL,
    year INTEGER NOT NULL,
    is_weekend BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS dim_time (
    time_key INTEGER PRIMARY KEY,
    hour INTEGER NOT NULL,
    minute INTEGER NOT NULL,
    second INTEGER NOT NULL DEFAULT 0,
    time_bucket VARCHAR(16) NOT NULL,
    is_peak_hour BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS dim_user (
    user_key BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    user_hash VARCHAR(128) NOT NULL,
    user_type VARCHAR(32) NOT NULL DEFAULT 'customer',
    registered_date_key INTEGER REFERENCES dim_date(date_key),
    city VARCHAR(128),
    verification_status VARCHAR(32) NOT NULL DEFAULT 'unknown',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    valid_from TIMESTAMPTZ NOT NULL,
    valid_to TIMESTAMPTZ,
    is_current BOOLEAN NOT NULL DEFAULT TRUE,
    data_quality_flag VARCHAR(64)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_dim_user_current_user_id
    ON dim_user(user_id)
    WHERE is_current;

CREATE INDEX IF NOT EXISTS ix_dim_user_registered_date_key
    ON dim_user(registered_date_key);

CREATE TABLE IF NOT EXISTS dim_driver (
    driver_key BIGSERIAL PRIMARY KEY,
    driver_id VARCHAR(64) NOT NULL,
    driver_hash VARCHAR(128) NOT NULL,
    vehicle_type VARCHAR(32) NOT NULL DEFAULT 'unknown',
    vehicle_brand VARCHAR(64),
    driver_status VARCHAR(32) NOT NULL DEFAULT 'unknown',
    join_date_key INTEGER REFERENCES dim_date(date_key),
    city VARCHAR(128),
    rating_bucket VARCHAR(32) NOT NULL DEFAULT 'unknown',
    valid_from TIMESTAMPTZ NOT NULL,
    valid_to TIMESTAMPTZ,
    is_current BOOLEAN NOT NULL DEFAULT TRUE,
    data_quality_flag VARCHAR(64)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_dim_driver_current_driver_id
    ON dim_driver(driver_id)
    WHERE is_current;

CREATE INDEX IF NOT EXISTS ix_dim_driver_join_date_key
    ON dim_driver(join_date_key);

CREATE TABLE IF NOT EXISTS dim_merchant (
    merchant_key BIGSERIAL PRIMARY KEY,
    merchant_id VARCHAR(64) NOT NULL,
    merchant_name_masked VARCHAR(256),
    merchant_category VARCHAR(128) NOT NULL DEFAULT 'unknown',
    city VARCHAR(128),
    district VARCHAR(128),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    trust_score_bucket VARCHAR(32) NOT NULL DEFAULT 'unknown',
    valid_from TIMESTAMPTZ NOT NULL,
    valid_to TIMESTAMPTZ,
    is_current BOOLEAN NOT NULL DEFAULT TRUE,
    data_quality_flag VARCHAR(64)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_dim_merchant_current_merchant_id
    ON dim_merchant(merchant_id)
    WHERE is_current;

CREATE TABLE IF NOT EXISTS dim_location_zone (
    zone_key BIGSERIAL PRIMARY KEY,
    zone_code VARCHAR(128) NOT NULL UNIQUE,
    city VARCHAR(128) NOT NULL DEFAULT 'unknown',
    district VARCHAR(128) NOT NULL DEFAULT 'unknown',
    province VARCHAR(128) NOT NULL DEFAULT 'unknown',
    country VARCHAR(128) NOT NULL DEFAULT 'Indonesia',
    latitude_center NUMERIC(10, 7),
    longitude_center NUMERIC(10, 7),
    zone_type VARCHAR(32) NOT NULL DEFAULT 'unknown'
);

CREATE INDEX IF NOT EXISTS ix_dim_location_zone_city_district
    ON dim_location_zone(city, district);

CREATE TABLE IF NOT EXISTS dim_service_type (
    service_type_key SERIAL PRIMARY KEY,
    service_code VARCHAR(64) NOT NULL UNIQUE,
    service_group VARCHAR(32) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS dim_payment_method (
    payment_method_key SERIAL PRIMARY KEY,
    payment_method_code VARCHAR(32) NOT NULL UNIQUE,
    payment_method_group VARCHAR(32) NOT NULL,
    requires_validation BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS dim_status (
    status_key SERIAL PRIMARY KEY,
    status_domain VARCHAR(32) NOT NULL,
    status_code VARCHAR(64) NOT NULL,
    status_group VARCHAR(32) NOT NULL,
    is_final_status BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (status_domain, status_code)
);

CREATE TABLE IF NOT EXISTS dim_cancel_reason (
    cancel_reason_key SERIAL PRIMARY KEY,
    reason_code VARCHAR(64) NOT NULL UNIQUE,
    reason_group VARCHAR(32) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS dim_promo (
    promo_key SERIAL PRIMARY KEY,
    promo_id VARCHAR(64) NOT NULL UNIQUE,
    promo_code VARCHAR(64) NOT NULL,
    promo_type VARCHAR(32) NOT NULL,
    campaign_name VARCHAR(256),
    start_date_key INTEGER REFERENCES dim_date(date_key),
    end_date_key INTEGER REFERENCES dim_date(date_key)
);

CREATE INDEX IF NOT EXISTS ix_dim_promo_date_range
    ON dim_promo(start_date_key, end_date_key);

CREATE TABLE IF NOT EXISTS dim_vehicle (
    vehicle_key SERIAL PRIMARY KEY,
    vehicle_type VARCHAR(32) NOT NULL UNIQUE,
    capacity INTEGER,
    service_eligibility VARCHAR(32) NOT NULL DEFAULT 'all'
);

CREATE TABLE IF NOT EXISTS fact_trip (
    trip_key BIGSERIAL PRIMARY KEY,
    trip_id VARCHAR(64) NOT NULL UNIQUE,
    request_date_key INTEGER NOT NULL REFERENCES dim_date(date_key),
    request_time_key INTEGER NOT NULL REFERENCES dim_time(time_key),
    user_key BIGINT NOT NULL REFERENCES dim_user(user_key),
    driver_key BIGINT REFERENCES dim_driver(driver_key),
    service_type_key INTEGER NOT NULL REFERENCES dim_service_type(service_type_key),
    pickup_zone_key BIGINT NOT NULL REFERENCES dim_location_zone(zone_key),
    dropoff_zone_key BIGINT NOT NULL REFERENCES dim_location_zone(zone_key),
    payment_method_key INTEGER NOT NULL REFERENCES dim_payment_method(payment_method_key),
    status_key INTEGER NOT NULL REFERENCES dim_status(status_key),
    cancel_reason_key INTEGER REFERENCES dim_cancel_reason(cancel_reason_key),
    promo_key INTEGER REFERENCES dim_promo(promo_key),
    vehicle_key INTEGER REFERENCES dim_vehicle(vehicle_key),
    fare_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    base_fare_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    distance_km NUMERIC(10, 3),
    duration_minute NUMERIC(10, 2),
    matching_duration_second NUMERIC(10, 2),
    waiting_duration_minute NUMERIC(10, 2),
    discount_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    platform_fee_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    driver_earning_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    is_completed SMALLINT NOT NULL DEFAULT 0,
    is_cancelled SMALLINT NOT NULL DEFAULT 0,
    is_felo_now SMALLINT NOT NULL DEFAULT 0,
    data_quality_flag VARCHAR(64) NOT NULL DEFAULT 'OK'
);

CREATE INDEX IF NOT EXISTS ix_fact_trip_date_service
    ON fact_trip(request_date_key, service_type_key);

CREATE INDEX IF NOT EXISTS ix_fact_trip_user_key
    ON fact_trip(user_key);

CREATE INDEX IF NOT EXISTS ix_fact_trip_driver_key
    ON fact_trip(driver_key);

CREATE INDEX IF NOT EXISTS ix_fact_trip_pickup_zone_key
    ON fact_trip(pickup_zone_key);

CREATE TABLE IF NOT EXISTS fact_trip_status_event (
    trip_status_event_key BIGSERIAL PRIMARY KEY,
    trip_id VARCHAR(64) NOT NULL,
    event_date_key INTEGER NOT NULL REFERENCES dim_date(date_key),
    event_time_key INTEGER NOT NULL REFERENCES dim_time(time_key),
    user_key BIGINT NOT NULL REFERENCES dim_user(user_key),
    driver_key BIGINT REFERENCES dim_driver(driver_key),
    status_key INTEGER NOT NULL REFERENCES dim_status(status_key),
    event_sequence_number INTEGER NOT NULL,
    duration_from_previous_second NUMERIC(10, 2),
    event_count INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS ix_fact_trip_status_event_trip_id
    ON fact_trip_status_event(trip_id);

CREATE INDEX IF NOT EXISTS ix_fact_trip_status_event_date_status
    ON fact_trip_status_event(event_date_key, status_key);

CREATE TABLE IF NOT EXISTS fact_food_order (
    food_order_key BIGSERIAL PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL UNIQUE,
    order_date_key INTEGER NOT NULL REFERENCES dim_date(date_key),
    order_time_key INTEGER NOT NULL REFERENCES dim_time(time_key),
    user_key BIGINT NOT NULL REFERENCES dim_user(user_key),
    merchant_key BIGINT NOT NULL REFERENCES dim_merchant(merchant_key),
    driver_key BIGINT REFERENCES dim_driver(driver_key),
    service_type_key INTEGER NOT NULL REFERENCES dim_service_type(service_type_key),
    merchant_zone_key BIGINT NOT NULL REFERENCES dim_location_zone(zone_key),
    dropoff_zone_key BIGINT NOT NULL REFERENCES dim_location_zone(zone_key),
    payment_method_key INTEGER NOT NULL REFERENCES dim_payment_method(payment_method_key),
    status_key INTEGER NOT NULL REFERENCES dim_status(status_key),
    cancel_reason_key INTEGER REFERENCES dim_cancel_reason(cancel_reason_key),
    promo_key INTEGER REFERENCES dim_promo(promo_key),
    gross_merchandise_value NUMERIC(18, 2) NOT NULL DEFAULT 0,
    item_count INTEGER NOT NULL DEFAULT 0,
    delivery_fee_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    service_fee_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    total_paid_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    merchant_earning_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    driver_earning_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    platform_fee_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    preparation_duration_minute NUMERIC(10, 2),
    delivery_duration_minute NUMERIC(10, 2),
    is_completed SMALLINT NOT NULL DEFAULT 0,
    is_cancelled SMALLINT NOT NULL DEFAULT 0,
    is_cod_locked SMALLINT NOT NULL DEFAULT 0,
    data_quality_flag VARCHAR(64) NOT NULL DEFAULT 'OK'
);

CREATE INDEX IF NOT EXISTS ix_fact_food_order_date_service
    ON fact_food_order(order_date_key, service_type_key);

CREATE INDEX IF NOT EXISTS ix_fact_food_order_user_key
    ON fact_food_order(user_key);

CREATE INDEX IF NOT EXISTS ix_fact_food_order_merchant_key
    ON fact_food_order(merchant_key);

CREATE INDEX IF NOT EXISTS ix_fact_food_order_driver_key
    ON fact_food_order(driver_key);

CREATE TABLE IF NOT EXISTS fact_food_order_item (
    food_order_item_key BIGSERIAL PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL,
    order_date_key INTEGER NOT NULL REFERENCES dim_date(date_key),
    merchant_key BIGINT NOT NULL REFERENCES dim_merchant(merchant_key),
    item_id VARCHAR(64) NOT NULL,
    item_name_masked VARCHAR(256),
    item_category VARCHAR(128) NOT NULL DEFAULT 'unknown',
    quantity INTEGER NOT NULL DEFAULT 0,
    unit_price_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    subtotal_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(18, 2) NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS ix_fact_food_order_item_order_id
    ON fact_food_order_item(order_id);

CREATE INDEX IF NOT EXISTS ix_fact_food_order_item_date_merchant
    ON fact_food_order_item(order_date_key, merchant_key);

CREATE INDEX IF NOT EXISTS ix_fact_food_order_item_item_id
    ON fact_food_order_item(item_id);

CREATE TABLE IF NOT EXISTS fact_shipment (
    shipment_key BIGSERIAL PRIMARY KEY,
    shipment_id VARCHAR(64) NOT NULL UNIQUE,
    send_order_id VARCHAR(64) NOT NULL,
    request_date_key INTEGER NOT NULL REFERENCES dim_date(date_key),
    request_time_key INTEGER NOT NULL REFERENCES dim_time(time_key),
    sender_user_key BIGINT NOT NULL REFERENCES dim_user(user_key),
    driver_key BIGINT REFERENCES dim_driver(driver_key),
    service_type_key INTEGER NOT NULL REFERENCES dim_service_type(service_type_key),
    pickup_zone_key BIGINT NOT NULL REFERENCES dim_location_zone(zone_key),
    dropoff_zone_key BIGINT NOT NULL REFERENCES dim_location_zone(zone_key),
    payment_method_key INTEGER NOT NULL REFERENCES dim_payment_method(payment_method_key),
    payer_type VARCHAR(16) NOT NULL DEFAULT 'sender',
    status_key INTEGER NOT NULL REFERENCES dim_status(status_key),
    cancel_reason_key INTEGER REFERENCES dim_cancel_reason(cancel_reason_key),
    distance_km NUMERIC(10, 3),
    package_weight_kg NUMERIC(10, 3),
    delivery_fee_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    insurance_fee_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    total_paid_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    platform_fee_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    driver_earning_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    pickup_duration_minute NUMERIC(10, 2),
    delivery_duration_minute NUMERIC(10, 2),
    is_completed SMALLINT NOT NULL DEFAULT 0,
    is_cancelled SMALLINT NOT NULL DEFAULT 0,
    has_proof_of_delivery SMALLINT NOT NULL DEFAULT 0,
    data_quality_flag VARCHAR(64) NOT NULL DEFAULT 'OK'
);

CREATE INDEX IF NOT EXISTS ix_fact_shipment_date_service
    ON fact_shipment(request_date_key, service_type_key);

CREATE INDEX IF NOT EXISTS ix_fact_shipment_sender_user_key
    ON fact_shipment(sender_user_key);

CREATE INDEX IF NOT EXISTS ix_fact_shipment_driver_key
    ON fact_shipment(driver_key);

CREATE INDEX IF NOT EXISTS ix_fact_shipment_pickup_zone_key
    ON fact_shipment(pickup_zone_key);

CREATE TABLE IF NOT EXISTS fact_matching_attempt (
    matching_attempt_key BIGSERIAL PRIMARY KEY,
    match_request_id VARCHAR(64) NOT NULL,
    domain_order_id VARCHAR(64) NOT NULL,
    attempt_date_key INTEGER NOT NULL REFERENCES dim_date(date_key),
    attempt_time_key INTEGER NOT NULL REFERENCES dim_time(time_key),
    service_type_key INTEGER NOT NULL REFERENCES dim_service_type(service_type_key),
    user_key BIGINT NOT NULL REFERENCES dim_user(user_key),
    driver_key BIGINT NOT NULL REFERENCES dim_driver(driver_key),
    pickup_zone_key BIGINT NOT NULL REFERENCES dim_location_zone(zone_key),
    status_key INTEGER NOT NULL REFERENCES dim_status(status_key),
    attempt_number INTEGER NOT NULL,
    matching_duration_second NUMERIC(10, 2),
    driver_distance_to_pickup_km NUMERIC(10, 3),
    is_success SMALLINT NOT NULL DEFAULT 0,
    is_rejected SMALLINT NOT NULL DEFAULT 0,
    is_timeout SMALLINT NOT NULL DEFAULT 0,
    data_quality_flag VARCHAR(64) NOT NULL DEFAULT 'OK'
);

CREATE INDEX IF NOT EXISTS ix_fact_matching_attempt_request_id
    ON fact_matching_attempt(match_request_id);

CREATE INDEX IF NOT EXISTS ix_fact_matching_attempt_date_service
    ON fact_matching_attempt(attempt_date_key, service_type_key);

CREATE INDEX IF NOT EXISTS ix_fact_matching_attempt_driver_key
    ON fact_matching_attempt(driver_key);

CREATE TABLE IF NOT EXISTS fact_payment (
    payment_key BIGSERIAL PRIMARY KEY,
    payment_id VARCHAR(64) NOT NULL UNIQUE,
    invoice_id VARCHAR(64),
    domain_order_id VARCHAR(64) NOT NULL,
    idempotency_key VARCHAR(128),
    payment_date_key INTEGER NOT NULL REFERENCES dim_date(date_key),
    payment_time_key INTEGER NOT NULL REFERENCES dim_time(time_key),
    user_key BIGINT NOT NULL REFERENCES dim_user(user_key),
    service_type_key INTEGER NOT NULL REFERENCES dim_service_type(service_type_key),
    payment_method_key INTEGER NOT NULL REFERENCES dim_payment_method(payment_method_key),
    status_key INTEGER NOT NULL REFERENCES dim_status(status_key),
    amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    paid_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    failed_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    refund_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    payment_processing_second NUMERIC(10, 2),
    is_success SMALLINT NOT NULL DEFAULT 0,
    is_failed SMALLINT NOT NULL DEFAULT 0,
    is_refunded SMALLINT NOT NULL DEFAULT 0,
    data_quality_flag VARCHAR(64) NOT NULL DEFAULT 'OK'
);

CREATE INDEX IF NOT EXISTS ix_fact_payment_date_service
    ON fact_payment(payment_date_key, service_type_key);

CREATE INDEX IF NOT EXISTS ix_fact_payment_method_status
    ON fact_payment(payment_method_key, status_key);

CREATE INDEX IF NOT EXISTS ix_fact_payment_idempotency_key
    ON fact_payment(idempotency_key);

CREATE INDEX IF NOT EXISTS ix_fact_payment_domain_order_id
    ON fact_payment(domain_order_id);

CREATE TABLE IF NOT EXISTS fact_wallet_transaction (
    wallet_transaction_key BIGSERIAL PRIMARY KEY,
    wallet_transaction_id VARCHAR(64) NOT NULL UNIQUE,
    wallet_id VARCHAR(64) NOT NULL,
    reference_id VARCHAR(64),
    transaction_date_key INTEGER NOT NULL REFERENCES dim_date(date_key),
    transaction_time_key INTEGER NOT NULL REFERENCES dim_time(time_key),
    user_key BIGINT REFERENCES dim_user(user_key),
    driver_key BIGINT REFERENCES dim_driver(driver_key),
    service_type_key INTEGER NOT NULL REFERENCES dim_service_type(service_type_key),
    transaction_type VARCHAR(32) NOT NULL,
    credit_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    debit_amount NUMERIC(18, 2) NOT NULL DEFAULT 0,
    balance_after_amount NUMERIC(18, 2),
    is_settlement SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS ix_fact_wallet_transaction_date_service
    ON fact_wallet_transaction(transaction_date_key, service_type_key);

CREATE INDEX IF NOT EXISTS ix_fact_wallet_transaction_wallet_id
    ON fact_wallet_transaction(wallet_id);

CREATE INDEX IF NOT EXISTS ix_fact_wallet_transaction_user_driver
    ON fact_wallet_transaction(user_key, driver_key);

CREATE TABLE IF NOT EXISTS fact_feedback (
    feedback_key BIGSERIAL PRIMARY KEY,
    feedback_id VARCHAR(64) NOT NULL UNIQUE,
    domain_order_id VARCHAR(64) NOT NULL,
    feedback_date_key INTEGER NOT NULL REFERENCES dim_date(date_key),
    feedback_time_key INTEGER NOT NULL REFERENCES dim_time(time_key),
    user_key BIGINT NOT NULL REFERENCES dim_user(user_key),
    driver_key BIGINT REFERENCES dim_driver(driver_key),
    merchant_key BIGINT REFERENCES dim_merchant(merchant_key),
    service_type_key INTEGER NOT NULL REFERENCES dim_service_type(service_type_key),
    rating_value NUMERIC(3, 2) NOT NULL,
    review_count INTEGER NOT NULL DEFAULT 1,
    has_text_review SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS ix_fact_feedback_date_service
    ON fact_feedback(feedback_date_key, service_type_key);

CREATE INDEX IF NOT EXISTS ix_fact_feedback_driver_key
    ON fact_feedback(driver_key);

CREATE INDEX IF NOT EXISTS ix_fact_feedback_merchant_key
    ON fact_feedback(merchant_key);

CREATE TABLE IF NOT EXISTS fact_geofence_event (
    geofence_event_key BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(64) NOT NULL UNIQUE,
    domain_order_id VARCHAR(64) NOT NULL,
    event_date_key INTEGER NOT NULL REFERENCES dim_date(date_key),
    event_time_key INTEGER NOT NULL REFERENCES dim_time(time_key),
    user_key BIGINT REFERENCES dim_user(user_key),
    driver_key BIGINT REFERENCES dim_driver(driver_key),
    merchant_key BIGINT REFERENCES dim_merchant(merchant_key),
    service_type_key INTEGER NOT NULL REFERENCES dim_service_type(service_type_key),
    location_zone_key BIGINT NOT NULL REFERENCES dim_location_zone(zone_key),
    rule_code VARCHAR(64) NOT NULL,
    event_count INTEGER NOT NULL DEFAULT 1,
    distance_violation_km NUMERIC(10, 3),
    risk_score NUMERIC(5, 2),
    is_blocked SMALLINT NOT NULL DEFAULT 0,
    is_otp_required SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS ix_fact_geofence_event_date_service
    ON fact_geofence_event(event_date_key, service_type_key);

CREATE INDEX IF NOT EXISTS ix_fact_geofence_event_rule_code
    ON fact_geofence_event(rule_code);

CREATE INDEX IF NOT EXISTS ix_fact_geofence_event_location_zone_key
    ON fact_geofence_event(location_zone_key);
