# FELO Mobile App — Product Requirements Document

| Field | Value |
|---|---|
| Version | 1.0 — MVP |
| Status | Draft — For Review |
| Last Updated | April 2026 |
| Platform | iOS (Swift) + Android (Kotlin) |
| Audience | PM, Backend Engineer, Mobile Engineer, AI Agent |

---

## Table of Contents

1. [Overview](#1-overview)
2. [Goals](#2-goals)
3. [User Stories](#3-user-stories)
4. [Functional Requirements](#4-functional-requirements)
   - 4.1 FELO-City Regular Ride Flow
   - 4.2 FELO-Now QR Instant Ride
   - 4.3 FELO Wallet & Settlement
   - 4.5 Dynamic Pricing Engine
   - 4.6 Driver Verification (KYC)
   - 4.7 Order Fiktif Mitigation
5. [Non-Functional Requirements](#5-non-functional-requirements)
   - 5.1 Performance
   - 5.2 Notification Delivery SLA
   - 5.3 Reliability & Availability
   - 5.4 Security & GPS Data Retention
   - 5.5 Platform Requirements
   - 5.6 Real-Time Fallback Strategy
   - 5.7 Microservices Deployment Strategy
   - 5.8 Analytics & Funnel Tracking
6. [Edge Cases & Error Handling](#6-edge-cases--error-handling)
- [Appendix A: Event Reference](#appendix-a-event-reference)
- [Appendix B: Cross-Context Rules Summary](#appendix-b-cross-context-rules-summary)

---

## 1. Overview

FELO is a hybrid super app platform designed for urban mobility and everyday services. The MVP release includes three verticals:

- **FELO-City (Ride Hailing)** — Core Priority
- **FELO-Food (Food Delivery)** — Limited MVP
- **FELO-Send (Package Delivery)** — Limited MVP

FELO-City is the primary focus for reliability and performance. FELO-Food and FELO-Send share the same infrastructure but operate with limited feature sets in MVP.

### 1.1 Product Suite

| Domain | Name | MVP Status |
|---|---|---|
| Ride Hailing | FELO-City | ✅ In Scope — Full MVP (Core Priority) |
| Ride Hailing | FELO-Now (QR Instant Ride) | ✅ In Scope — MVP Submodule |
| Financial | Wallet & Settlement | ✅ In Scope — FELO Wallet Only |
| Trust & Safety | Order Fiktif Mitigation | ✅ In Scope — Full Implementation |
| Food Delivery | FELO-Food | ✅ In Scope — Limited MVP |
| Package Delivery | FELO-Send | ✅ In Scope — Limited MVP |

### 1.2 Single Account Ecosystem

All FELO services share one unified account system:
- One user account works across all services.
- One driver profile covers ride, food, and send operations.
- One FELO Wallet is used for all transactions.

### 1.3 Service Architecture Summary

> **AI AGENT NOTE:** This PRD describes mobile-facing behavior. All service boundaries, event schemas, and cross-context communication rules are governed by **FELO Context Map v1**. Agents must reference the context map for backend service ownership.

| Context | Services | Role |
|---|---|---|
| Ride Mobility (CORE) | ride-service, matching-service, tracking-service, pricing-service, location-service | Trip lifecycle & GPS |
| Identity | user-service, driver-service, auth-service | Auth & profiles |
| Financial | wallet-service, payment-service, invoice-service | All monetary flow |
| Experience | feedback-service, notification-service | Ratings & notifications |

---

## 2. Goals

### 2.1 Product Goals

| # | Goal | Success Metric |
|---|---|---|
| G1 | Deliver a functional ride-hailing MVP (FELO-City) | Ride request → match → complete in < 3 min average |
| G2 | Enable instant QR-based rides (FELO-Now) | QR scan → ride start in < 60 seconds |
| G3 | Provide a secure, unified wallet for all payments | Zero failed settlements on completed rides |
| G4 | Mitigate fraudulent and fictitious orders | < 1% order fraud rate at launch |
| G5 | Establish scalable architecture for super app expansion | Food & Send launch ready with no core service re-architecture |

### 2.2 Business Goals

- Build driver and customer trust through real-time location transparency.
- Prevent revenue leakage from fake rides and fraudulent trip completions.
- Create a unified identity and wallet ecosystem that locks users into the FELO platform.
- Support rapid market entry with a lean, high-reliability MVP.

### 2.3 Out of Scope for MVP

- Merchant POS management interface.
- Promo and discount engine.
- In-app customer support chat.
- Multi-language support (Indonesian only for MVP).

---

## 3. User Stories

### 3.1 Customer — FELO-City Regular Ride

| ID | Story | Acceptance Criteria |
|---|---|---|
| US-C01 | As a customer, I want to input my pickup and destination so I can get a price estimate before ordering. | Pricing inquiry screen shown before order submission. Fare calculated by pricing-service. |
| US-C02 | As a customer, I want to be matched with a nearby driver within 5 seconds so I don't wait long. | Match confirmed ≤ 5 sec. If no match in 15 sec, user is notified and may retry. |
| US-C03 | As a customer, I want to track my driver's live location so I know when they will arrive. | Driver GPS position updates every 3 seconds on the tracking screen. |
| US-C04 | As a customer, I want to pay using my FELO Wallet so I don't need to carry cash. | Wallet balance shown before order. Deducted immediately on ride completion. |
| US-C05 | As a customer, I want to cancel my order with a clear cancellation policy. | Cancel allowed before driver arrival. Cancellation fee applies if driver is already en_route. |
| US-C06 | As a customer, I want to rate my driver after each trip. | Rating screen shown after ride completion. Minimum 1 star required to dismiss. |

### 3.2 Customer — FELO-Now QR Instant Ride

| ID | Story | Acceptance Criteria |
|---|---|---|
| US-N01 | As a customer, I want to set my destination first and then generate a QR code so a nearby driver can pick me up instantly. | Destination required before QR generation. QR generated by ride-service and displayed on screen. |
| US-N02 | As a customer, I want to see fare details before the driver validates the ride. | Driver's scan triggers fare display. Customer and driver agree verbally before driver taps Accept. |
| US-N03 | As a customer, I want my QR to expire if not scanned within a reasonable time. | QR expires after 10 minutes. Expired QR shows refresh option. |

### 3.3 Driver

| ID | Story | Acceptance Criteria |
|---|---|---|
| US-D01 | As a driver, I want to go online/offline instantly so I can control my availability. | Toggle syncs with driver-service ≤ 2 sec. Offline = no new match requests. |
| US-D02 | As a driver, I want to receive ride requests with pickup distance and fare so I can decide to accept or reject. | Request card shows distance, fare, and destination. Accept/Reject within 20 seconds. |
| US-D03 | As a driver, I want to scan a QR to accept a FELO-Now ride. | QR scanner opens from driver home screen. Scan shows: destination, estimated fare, customer info. |
| US-D04 | As a driver, I want my earnings to be settled immediately after trip completion. | Settlement triggered by ride.completed.v1 event. Balance reflected in FELO Wallet ≤ 30 seconds. |
| US-D05 | As a driver, I want to see my daily earnings summary. | Earnings tab shows daily total, per-trip breakdown, and settlement history. |

### 3.4 Customer — Order Fiktif Mitigation (Trust & Safety)

> **SCOPE NOTE:** Order Fiktif Mitigation is currently scoped for **FELO-Food** and **FELO-Send** (delivery services). It is **NOT** applied to FELO-City ride-hailing in the MVP, as the ride passenger is physically present with the driver.

| ID | Story | Acceptance Criteria |
|---|---|---|
| US-F01 | As a customer ordering for myself, I want to confirm I am physically present at the delivery location by sharing my live location. | Order cannot be placed unless live location sharing is active. Location must be within reasonable radius of delivery address. |
| US-F02 | As a customer ordering for a friend, I want to input my friend's phone number so they can validate receipt. | Phone number field shown for 'Order for Friend' mode. Validation notification sent to recipient via app, WhatsApp, or SMS. |
| US-F03 | As a recipient, I want to accept or decline an incoming delivery addressed to me. | Recipient receives deep-link notification. Tap Accept confirms presence. Decline notifies sender. |
| US-F04 | As a recipient without the FELO app, I want to be notified via WhatsApp or SMS. | System falls back to WhatsApp Business API, then SMS if WhatsApp unavailable. Fallback link opens web validation page. |
| US-F05 | As a system, I want the recipient to share their live location to prevent secondary-number fraud. | After recipient accepts, live location sharing is required and continuously monitored during delivery. Driver sees recipient's live location. |

---

## 4. Functional Requirements

### 4.1 FELO-City Regular Ride Flow

#### 4.1.1 Trip State Machine

> **AI AGENT GUARDRAIL:** The trip state machine lives exclusively in **ride-service**. No other service may modify trip state directly. All state transitions must publish the corresponding domain event to RabbitMQ.

| State | Trigger | Published Event | Next State |
|---|---|---|---|
| `pricing_inquiry` | Customer submits pickup + destination | — | `requested` |
| `requested` | Customer confirms fare and orders | `ride.created.v1` | `matching` |
| `matching` | matching-service finds driver | `driver.matched.v1` | `en_route` |
| `en_route` | Driver accepts and heads to pickup | — | `arrived` |
| `arrived` | Driver marks arrival at pickup | — | `on_ride` |
| `on_ride` | Customer boards, trip begins | `ride.started.v1` | `completed` |
| `completed` | Driver marks arrival at destination | `ride.completed.v1` | — |
| `cancelled` | Customer or system cancels | `ride.cancelled.v1` | — |

#### 4.1.2 Matching Requirements

- Matching target: driver found and confirmed within **5 seconds**.
- matching-service consumes GPS stream from tracking-service and driver availability from driver-service.
- If no driver found within 15 seconds, notify customer and offer retry.
- Rejected rides are re-queued immediately to the next nearest available driver.
- matching-service does **NOT** own driver availability state — it reads from driver-service.

#### 4.1.3 Matching Retry Behavior

If no driver is found within 15 seconds:

- System automatically retries matching up to **3 times**.
- Each retry **expands the search radius** incrementally.
- Customer is informed: *"Still searching for drivers..."*
- After max retries, user is prompted to retry manually.

#### 4.1.4 Real-Time Tracking

- tracking-service handles live WebSocket/MQTT GPS stream.
- Driver GPS position broadcast to customer every **3 seconds** while trip is active.
- location-service responsible for routing, ETA calculation, and reverse-geocoding.
- tracking-service is the single source of truth for live GPS data.

---

### 4.2 FELO-Now QR Instant Ride

> **AI AGENT GUARDRAIL:** FELO-Now is a **submodule of ride-service**. It is NOT a separate microservice. The standard ride state machine applies starting from `on_ride`. FELO-Now bypasses the matching phase entirely.

#### 4.2.1 Customer Flow

1. Customer opens the FELO-Now tab.
2. Customer sets destination address. **QR cannot be generated without a destination.**
3. System calls pricing-service to calculate estimated fare.
4. Customer taps **"Generate QR"**. ride-service creates a trip draft and returns a QR code.
5. QR code is displayed full-screen with destination and estimated fare visible.
6. QR is valid for **10 minutes**. Countdown timer shown. Auto-refresh option on expiry.
7. Customer shows QR to driver for scanning.

#### 4.2.2 Driver Flow

1. Driver taps **"Scan QR"** on the driver home screen.
2. Native camera QR scanner opens.
3. On successful scan, trip details card displays: customer name, destination, estimated fare, distance.
4. Driver discusses fare and destination verbally with the passenger.
5. Driver taps **"Accept & Start Ride"** to validate. Trip state transitions to `on_ride`.
6. Driver taps **"Decline"** if terms are not agreed. QR remains valid for the customer to try another driver.

#### 4.2.3 FELO-Now Trip State Extension

| State | Description |
|---|---|
| `qr_generated` | QR created by ride-service, waiting for a driver to scan |
| `scanned` | Driver has scanned the QR; details displayed to driver |
| `accepted` | Driver taps Accept & Start Ride; QR is consumed |
| `on_ride` | Trip started; standard ride state machine applies from here |

**QR Locking Rules:**
- Only **one driver** may hold a QR in `scanned` state at a time.
- When a driver scans a QR, the QR is **locked** for that driver.
- Any other scan attempt while QR is locked returns: *"This QR is currently being reviewed by another driver."*
- If the locking driver declines, the QR is **unlocked** and available for another driver to scan.

#### 4.2.4 FELO-Now Driver Cancellation

- Driver may cancel **before** ride starts (`accepted` state) without penalty.
- Upon cancellation:
  - QR is unlocked and returned to `qr_generated` state.
  - Customer may present QR to another driver.
  - Customer is notified: *"The driver declined your ride. Please show your QR to another driver."*
- System logs cancellation frequency per driver for fraud and abuse detection.

---

### 4.3 FELO Wallet & Settlement#### 4.3.1 Wallet Operations — MVP Scope

| Operation | Actor | Trigger |
|---|---|---|
| Top-up wallet | Customer | Manual top-up via payment gateway |
| Ride payment deduction | Customer | `ride.completed.v1` event |
| Driver earning credit (instant) | Driver | `ride.completed.v1` event — credited immediately, no holding period |
| Food order payment deduction | Customer | `order.completed.v1` event |
| Driver/courier earning credit (instant) | Driver / Courier | `order.completed.v1` event — credited immediately |
| Send package payment deduction | Customer | `shipment.delivered.v1` event |
| Courier earning credit (instant) | Courier | `shipment.delivered.v1` event — credited immediately |
| View balance & history | Customer / Driver | On-demand from wallet screen |

#### 4.3.2 Settlement Rules

- Settlement executes **instantly and directly** upon the service completion event (`ride.completed.v1`, `order.completed.v1`, `shipment.delivered.v1`).
- Earnings are credited **immediately into the driver's / courier's FELO Wallet** — there is no holding period, escrow, or manual disbursement step.
- wallet-service is the **single source of truth** for all balance ledgers.
- payment-service orchestrates payment flow but does **not** own the ledger.
- All transactions recorded in invoice-service for audit trail.
- Driver / courier wallet balance is updated within **30 seconds** of the completion event.

#### 4.3.3 Idempotency & Financial Safety

> **AI AGENT GUARDRAIL:** wallet-service must be fully idempotent. Any logic that processes a completion event without checking for prior settlement is a critical bug.

- Each transaction must include a unique `transaction_id`.
- wallet-service must enforce idempotency using `transaction_id` — duplicate submissions return the original result without re-processing.
- Duplicate events (e.g., duplicate `ride.completed.v1`) must **not** result in duplicate balance updates.
- wallet-service must validate that a `ride_id` / `order_id` / `shipment_id` has **not** been settled previously before processing.

---

### 4.5 Dynamic Pricing Engine

> **AI AGENT GUARDRAIL:** All fare calculations must go through pricing-service exclusively. No other service may compute fares independently. Pricing must be deterministic — the same inputs must always produce the same output for audit reproducibility.

#### 4.5.1 Pricing Inputs

pricing-service must support dynamic pricing based on the following factors:

| Factor | Description |
|---|---|
| Distance | Trip distance in kilometers |
| Duration | Estimated trip duration in minutes |
| Demand level | Number of active ride requests in the area |
| Supply level | Number of available drivers in the area |

#### 4.5.2 Surge Pricing Rules

- Surge pricing multiplier is applied when the **demand / supply ratio exceeds a defined threshold**.
- The multiplier is calculated as: `base_fare × surge_multiplier`.
- Surge multiplier and threshold values are configurable by ops (not hardcoded).
- Surge pricing must be **transparently displayed** to the customer before order confirmation.

#### 4.5.3 Pricing Calculation Points

Pricing must be calculated at two points in the trip lifecycle:

1. **Before order confirmation** — shown to customer as the estimated fare. Customer must see and accept the fare before the ride is placed.
2. **At ride completion** — final fare is recalculated based on actual distance and duration. A final adjustment is allowed within a defined tolerance (e.g., ±10% of the estimated fare).

#### 4.5.4 Audit & Reproducibility

- All pricing calculations must be **deterministic and reproducible**.
- Each fare calculation must be logged with its full input set (distance, duration, demand, supply, multiplier, timestamp) for audit purposes.
- If a customer disputes a fare, the exact inputs and formula used must be retrievable from the invoice-service log.

---

### 4.6 Driver Verification (KYC)

Driver must complete the following before going online:

| Requirement | Details |
|---|---|
| Identity verification | Upload KTP (National ID). Reviewed and approved by ops. |
| Vehicle verification | Upload vehicle registration (STNK) and photo of the vehicle. |
| Profile photo | Clear face photo. Must match KTP photo. |

- Driver **cannot** toggle online until all three verification items are approved.
- Verification status is owned by driver-service.
- Driver app shows a locked state with pending verification progress if not yet approved.
- Re-verification is triggered if document expiry is detected or ops flags the account.

---

### 4.7 Order Fiktif Mitigation

> **DESIGN CONTEXT:** This feature addresses fraudulent orders where a user places an order with no real recipient present — by listing a fake address, using their own secondary phone number as the recipient, or having no one at the delivery location. This mitigation applies to delivery-based services (FELO-Food, FELO-Send) and is integrated at the **checkout stage**.

#### 4.7.1 Ordering Mode Selection

At the checkout screen, the customer must select one of two modes before the order can be submitted:

| Dimension | "Kirim untuk Saya Sendiri" | "Kirim untuk Teman" |
|---|---|---|
| Payment Methods | FELO Wallet | FELO Wallet (preferred) or COD with conditions |
| Location Requirement | Customer must share live location throughout delivery | Recipient must share live location after accepting order |
| Phone Validation | Not required | Recipient phone number required |
| Notification to Recipient | N/A | App push, WhatsApp, or SMS |
| Recipient Acceptance | N/A | Required — order held until accepted or timed out |

#### 4.7.2 Mode A: Kirim untuk Saya Sendiri

1. Customer selects **"Kirim untuk Saya Sendiri"**.
2. System checks if location permission is granted. If denied, user is prompted to enable it. **Order cannot proceed without location.**
3. Customer completes checkout. Order placed.
4. Customer's live location is continuously shared with the driver during delivery.
5. If customer's live location sharing stops mid-delivery, driver receives an alert. System flags the trip for review.

#### 4.7.3 Mode B: Kirim untuk Teman

1. Customer selects **"Kirim untuk Teman"**.
2. Customer inputs recipient's phone number. Phone number is validated (format check + must not be the same as the customer's own number).
3. Payment: FELO Wallet is the default. COD is available **only if** recipient phone validation is completed.
4. System sends a validation notification to the recipient:
   - If recipient has FELO app → push notification with deep link to accept/decline screen.
   - If recipient does not have FELO app → WhatsApp message via WhatsApp Business API with web validation link.
   - If WhatsApp fails → SMS fallback with web validation link.
5. Recipient must tap **"Accept"** within the validation timeout window (default: **10 minutes**). Order is held in `pending` state until accepted.
6. If recipient declines or timeout occurs, customer is notified and order is cancelled.
7. After recipient accepts, **recipient's live location sharing is activated**. Recipient must grant location permission to complete the acceptance flow.
8. Driver can see recipient's live location on the tracking screen throughout the delivery.
9. If recipient's live location stops mid-delivery, driver receives an alert. System flags the trip for review.

#### 4.7.4 Fraud Risk Scoring

> **DESIGN NOTE:** The distance check is a risk signal, not a hard block. Orders from suspicious patterns are flagged for ops review, not automatically cancelled, to avoid false positives on legitimate use cases (e.g., a user ordering for a family member in the same household).

Recipient validation includes a multi-signal risk score:

| Signal | Check | Action |
|---|---|---|
| Distance check | Recipient GPS vs. sender GPS at time of acceptance | If < 50 meters → flag as high risk |
| Device fingerprint | Recipient's device fingerprint vs. sender's device fingerprint | If match → flag as high risk |
| Phone uniqueness | Recipient phone number vs. sender phone number | If match → block order (hard rule) |

**Risk Score Behavior:**
- If distance < 50 meters: order is **NOT blocked**. Order is flagged as `high_risk` and risk score is increased.
- If device fingerprint matches: order is **NOT blocked**. Risk score is increased.
- If phone number matches sender's own number: order is **blocked** immediately (this is the only hard block).
- Orders with a high cumulative risk score are routed to **ops manual review queue**.
- Repeated high-risk patterns from the same account trigger account-level review.

---

## 5. Non-Functional Requirements

### 5.1 Performance

| ID | Requirement | Target | Measurement |
|---|---|---|---|
| NFR-P01 | Driver matching speed | ≤ 5 seconds from `ride.created.v1` | 95th percentile in matching-service logs |
| NFR-P02 | Live GPS tracking update frequency | Every 3 seconds while trip is active | Measured at tracking-service WebSocket/MQTT emission rate |
| NFR-P03 | Wallet settlement latency | ≤ 30 seconds post `ride.completed.v1` | Time from event publish to wallet balance reflection |
| NFR-P04 | QR generation latency | ≤ 2 seconds from destination input | Client-measured from API call to QR display |
| NFR-P05 | App launch to home screen | ≤ 3 seconds on mid-range device | Cold start measured on reference devices |

### 5.2 Notification Delivery SLA

| Notification Type | Primary Channel | Max Delivery Time | Fallback |
|---|---|---|---|
| Recipient order validation | FELO Push Notification | ≤ 5 seconds | WhatsApp → SMS |
| WhatsApp validation fallback | WhatsApp Business API | ≤ 30 seconds | SMS |
| SMS validation fallback | SMS Gateway | ≤ 60 seconds | Ops manual review |
| Ride status updates | FELO Push Notification | ≤ 3 seconds | In-app polling |
| Settlement confirmation | FELO Push Notification | ≤ 30 seconds after completion event | In-app wallet refresh |

### 5.3 Reliability & Availability

- Core ride and matching services: **99.9% uptime SLA**.
- Wallet and settlement services: **99.95% uptime SLA** (financial critical path).
- Notification service: best-effort with multi-channel fallback (push → WhatsApp → SMS).
- Graceful degradation: if matching-service is degraded, app shows "Service temporarily unavailable" rather than crashing.

### 5.4 Security

- All API communication over **HTTPS / TLS 1.3**.
- JWT-based auth with short-lived access tokens and refresh token rotation.
- Recipient phone number is hashed and not displayed in full in the UI.
- QR codes are single-use and **cryptographically signed** by ride-service.

**GPS Data Retention:**
- Raw GPS data is stored for **48 hours** for dispute resolution.
- Aggregated trip data is stored long-term for analytics and fraud detection.
- Raw GPS data is not accessible to any party outside an active trip or open dispute.

### 5.5 Platform Requirements

| Dimension | iOS | Android |
|---|---|---|
| Language | Swift 5.9+ | Kotlin 1.9+ |
| Min OS Version | iOS 15.0+ | Android 8.0 (API 26)+ |
| Location Services | CoreLocation (Always / WhenInUse) | FusedLocationProviderClient |
| Push Notifications | APNs | Firebase Cloud Messaging (FCM) |
| QR Scanner | AVFoundation / Vision | Google ML Kit / ZXing |
| Maps / Routing | MapKit + location-service API | Google Maps SDK + location-service API |
| Real-time Transport | WebSocket (URLSessionWebSocketTask) | WebSocket (OkHttp / Ktor) |

---

### 5.6 Real-Time Fallback Strategy

**Primary transport:**
- WebSocket real-time updates for live GPS tracking and trip state changes.

**Fallback (if WebSocket connection fails):**
- In-app polling every **5 seconds** automatically activates.
- Client displays a subtle degraded-mode indicator.
- WebSocket reconnection is attempted in the background with exponential backoff.

**Push notifications are used for critical state changes only:**
- Driver matched, ride started, ride completed, payment settled, order fiktif validation requests.
- Push is not used for continuous GPS updates (WebSocket / polling handles this).

---

### 5.7 Microservices Deployment Strategy (MVP Constraint)

> **AI AGENT NOTE:** Logical service separation is enforced at the code and API level. Physical separation is optional during MVP to reduce operational overhead.

- All services **may** be deployed on a single cluster or VPS during MVP.
- Logical separation is enforced: each service has its own codebase, API contract, and data store.
- Inter-service communication remains via REST API or RabbitMQ events — **no direct database sharing**.
- The architecture must support migration to physically separated deployments (e.g., Kubernetes multi-node) without requiring API contract changes.

---

### 5.8 Analytics & Funnel Tracking

The system must track the following events for the ride funnel:

| Event | Description |
|---|---|
| `analytics.ride.requested` | Customer submitted a ride request |
| `analytics.ride.matched` | Driver successfully matched |
| `analytics.ride.started` | Trip entered `on_ride` state |
| `analytics.ride.completed` | Trip completed successfully |
| `analytics.ride.cancelled` | Trip cancelled at any state |

**Required Metrics:**
- Conversion rate: `requested → completed`
- Average matching time (ms)
- Cancellation rate by cancellation state (pre-match, en_route, on_ride)

**Rules:**
- All analytics events must be published to the event stream (RabbitMQ or dedicated analytics bus).
- Analytics events are **separate** from domain events — they must not be used for triggering business logic.
- Analytics data must be queryable for ops dashboards within 5 minutes of event publication.

---

## 6. Edge Cases & Error Handling

> **AI AGENT NOTE:** All edge cases below are **mandatory handling requirements** for MVP. Each case includes the expected system behavior, user-facing message, and the service responsible for detection.

---

### EC-01: Driver Goes Offline Mid-Ride

| Field | Detail |
|---|---|
| **Trigger** | Driver's device loses connection or driver-service receives offline signal while trip is in `on_ride` state. |
| **Responsible Service** | tracking-service detects GPS stream loss; ride-service monitors trip state. |
| **System Behavior** | 1. tracking-service emits `driver.connection.lost` event after 15-second no-signal window. 2. ride-service holds trip in `on_ride` state (does not auto-complete). 3. notification-service alerts customer that driver connection is unstable. 4. If driver reconnects within 5 minutes, trip resumes normally. 5. If driver does not reconnect, ops team is alerted. Trip is flagged for manual review and customer is offered a refund. |
| **User-Facing Message** | *"Your driver's connection is unstable. Please wait — we're monitoring the situation. If the issue persists, you will be contacted by our support team."* |

---

### EC-02: Customer Cancels After Driver Is En Route

| Field | Detail |
|---|---|
| **Trigger** | Customer taps Cancel while trip state is `en_route`. |
| **Responsible Service** | ride-service processes cancellation. wallet-service handles the fee. |
| **System Behavior** | 1. Cancellation fee calculated by pricing-service (based on distance driver has already travelled). 2. Fee shown to customer before final confirmation. Customer must confirm to proceed. 3. `ride.cancelled.v1` event published. 4. Driver notified immediately. Driver's trip is closed. 5. Cancellation fee deducted from customer's FELO Wallet. 6. Driver receives partial compensation. |
| **User-Facing Message** | *"Your driver is already on their way. A cancellation fee of Rp [X] will be charged. Are you sure you want to cancel?"* |

---

### EC-03: GPS / Location Permission Denied by User

| Field | Detail |
|---|---|
| **Trigger** | User denies location permission at OS prompt, or revokes permission from app settings. |
| **Responsible Service** | Mobile client (iOS / Android). No backend involvement. |
| **System Behavior (Regular Ride)** | 1. App detects no location permission on home screen load. 2. Full-screen permission request screen shown with explanation. 3. Core ride features are disabled until permission is granted. 4. Deep-link to OS Settings provided. |
| **System Behavior (Order Fiktif — Self Mode)** | Order cannot proceed. Checkout blocked with inline error. Same permission prompt shown. |
| **User-Facing Message** | *"Location access is required to use FELO. Please enable it in your phone settings to continue."* |

---

### EC-04: WhatsApp / SMS Validation Not Delivered to Recipient

| Field | Detail |
|---|---|
| **Trigger** | WhatsApp Business API and SMS gateway both fail to deliver within their respective SLA windows (30s / 60s). |
| **Responsible Service** | notification-service manages delivery orchestration and fallback logic. |
| **System Behavior** | 1. notification-service marks the delivery as failed after max retry attempts. 2. Sender is notified that validation could not be sent. 3. Customer is offered: (a) re-enter correct phone number, or (b) cancel the order. 4. Order remains in `pending` state. No charge is made. 5. Incident logged for ops review of gateway health. |
| **User-Facing Message** | *"We couldn't reach your recipient. Please check the phone number and try again, or cancel the order."* |

---

### EC-05: FELO-Now QR Expired Before Driver Scans

| Field | Detail |
|---|---|
| **Trigger** | QR code reaches its 10-minute TTL before any driver scans it. |
| **Responsible Service** | ride-service invalidates expired QR tokens. Mobile client shows countdown timer. |
| **System Behavior** | 1. QR is visually dimmed/greyed out on customer screen at expiry. 2. "QR Expired" label displayed. 3. "Generate New QR" button shown — destination and fare pre-filled. 4. Expired QR attempt by a driver returns error: *"This QR code has expired. Ask the passenger to generate a new one."* |
| **User-Facing Message** | *"Your QR code has expired. Tap below to generate a new one — your destination is already saved."* |

---

### EC-06: Recipient Never Accepts Validation (Timeout)

| Field | Detail |
|---|---|
| **Trigger** | Recipient does not tap Accept within the 10-minute validation window. |
| **Responsible Service** | notification-service tracks acceptance state. ride-service / order-service holds the order. |
| **System Behavior** | 1. At T+10 min, order transitions to `validation_expired` state. 2. Sender is notified via push notification. 3. Order is automatically cancelled. No charge applied. 4. Sender may re-place the order with a corrected phone number. 5. A reminder notification is sent at T+5 min warning the sender of pending expiry. |
| **User-Facing Message (Sender at T+5)** | *"Your recipient hasn't accepted yet. The order will be cancelled in 5 minutes if they don't respond."* |
| **User-Facing Message (Sender at T+10)** | *"Your order was cancelled because the recipient did not confirm in time. No charge was applied."* |

---

### EC-07: Shared Location Stops Mid-Journey

| Field | Detail |
|---|---|
| **Trigger** | Customer (self mode) or recipient (friend mode) stops sharing live location while delivery is active. |
| **Responsible Service** | tracking-service detects location stream loss. notification-service alerts relevant parties. |
| **System Behavior** | 1. tracking-service detects no location update for 30+ seconds from customer/recipient. 2. Driver receives in-app alert: *"Recipient location is unavailable."* 3. Sender receives push notification requesting re-enablement of location sharing. 4. Trip is flagged in the operations dashboard for monitoring. 5. If location is not restored within 5 minutes, trip is escalated to ops review. 6. Driver may choose to continue or report the issue via the trip menu. |
| **User-Facing Message (Recipient/Customer)** | *"Your live location has stopped. Please re-enable location sharing to keep your delivery on track."* |
| **User-Facing Message (Driver)** | *"The recipient's location is currently unavailable. Please proceed to the delivery address or contact support if needed."* |

---

## Appendix A: Event Reference

> All domain events are published to RabbitMQ. Services must consume events via the message broker. **Direct database access across service boundaries is strictly forbidden.**

| Event | Published By | Consumed By |
|---|---|---|
| `ride.created.v1` | ride-service | matching-service, notification-service |
| `driver.matched.v1` | matching-service | ride-service, notification-service |
| `ride.started.v1` | ride-service | tracking-service, notification-service |
| `ride.completed.v1` | ride-service | wallet-service, invoice-service, notification-service, feedback-service |
| `ride.cancelled.v1` | ride-service | wallet-service, notification-service |
| `order.created.v1` | order-service | notification-service, merchant-service |
| `order.completed.v1` | order-service | wallet-service, invoice-service, notification-service, feedback-service |
| `order.cancelled.v1` | order-service | wallet-service, notification-service |
| `shipment.created.v1` | send-order-service | notification-service |
| `shipment.picked_up.v1` | shipment-service | tracking-service, notification-service |
| `shipment.delivered.v1` | shipment-service | wallet-service, invoice-service, notification-service, feedback-service |
| `payment.completed.v1` | payment-service | wallet-service, invoice-service |
| `wallet.balance.updated.v1` | wallet-service | notification-service, user-service |
| `settlement.completed.v1` | wallet-service | notification-service, driver-service |
| `feedback.submitted.v1` | feedback-service | driver-service, notification-service |

---

## Appendix B: Cross-Context Rules Summary

> **MANDATORY — AI AGENT GUARDRAIL:** These rules are absolute. Any generated code, API design, or service interaction that violates these rules must be rejected and revised.

1. No service may read or write to another service's database directly.
2. All cross-context communication uses REST (synchronous) or RabbitMQ events (asynchronous).
3. **ride-service** is the single source of truth for trip lifecycle and state machine.
4. **wallet-service** is the single source of truth for all balance ledgers.
5. **tracking-service** owns GPS source of truth for live streaming.
6. **matching-service** reads (does not own) GPS stream and driver availability.
7. **driver-service** owns live driver availability status (Online/Offline).
8. **FELO-Now** is a submodule of ride-service. It is not a separate service.

---

*End of Document — FELO PRD v1.0 — Confidential*