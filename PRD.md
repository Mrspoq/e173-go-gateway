# E173 GO VOICE GATEWAY – PRODUCT REQUIREMENTS v0.2 (Concise)

**Version:** 0.2
**Date:** 2025-06-16

## Purpose
Provide an end-to-end platform that lets a small team run 100–250 Huawei E173 modems as a GSM voice/SMS termination farm with minimal manual work and full commercial, technical and compliance tooling.

## Actors
*   **Super-Admin** – owns the company.
*   **NOC Manager** – runs day-to-day ops, views finance but cannot withdraw funds.
*   **Technician** – swaps SIMs, reboots hubs, sees tech metrics only.
*   **Customer** – views CDR, balance, tops up, pulls API key.
*   **API Client** – machine user (wholesale customer).

## High-Level Architecture
Go (Gin) + HTMX UI, Postgres, Redis (caching & rate-limiting), Asterisk 18 + chan_dongle, optional Janus/gRPC API gateway, Prometheus + Grafana stack, Docker compose file for every component, optional K8s Helm chart.

## Functional Requirements

### 4.1 Hardware & Modem Management
*   Auto-discover new ttyUSB pairs; persist IMEI→logical-name map.
*   USB-hub port power-cycle command (uhubctl / python-usb) from UI & API.
*   Alarm if RSSI < −95 dBm, temperature > 50 °C, current draw > 450 mA.

### 4.2 SIM Inventory & Policies
*   **Fields:** ICCID, IMSI, phone-number, operator, status, activation-date, expiry-date, last-recharge, current-credit, lifetime-minutes, max-minutes-per-day, rotation-cooldown-min.
*   **Workflows:**
    *   Auto USSD/SMS recharge via YAML scenario; PIN-stock counter decremented.
    *   Auto-rotate SIM when policy limit hit or after X short calls.
    *   Flag SIM “replace” when three consecutive registration failures.

### 4.3 Call Routing & Termination
*   Multi-level LCR: customer-rate-deck → internal-cost-deck → SIM/operator group.
*   CLI strategy: random from customer pool, fixed, or per-destination override.
*   Call-duration shaper (Gaussian mean N seconds, σ configurable).
*   Real-time call cap enforcement: per-SIM, per-operator, per-customer.
*   SIP ingest (customer) → Asterisk dial-plan → dongle → PSTN.

### 4.4 SMS / USSD
*   Send/receive SMS, long SMS concatenation, delivery report handling.
*   USSD gateway (check balance, bundle purchase).

### 4.5 Anti-Spam / Fraud / Quality
*   Short-call detector (< 6 s) → score system → auto-blacklist.
*   ASR, ACD, PDD monitors with threshold alarms.
*   High-spend burst rule: if spend > X in Y min → block customer.
*   Optional bot-voice classifier (open-source TensorFlow lite model).

### 4.6 Customer & Billing
*   Pre-paid wallet, multi-currency, taxation rules.
*   Per-second billing + rounding mode per destination.
*   Credit-limit enforcement; optional post-paid with credit-limit.
*   Invoice PDF generation, Stripe & crypto payment hooks.
*   Customer self-care portal: balance, recharge, CDR CSV, rate download, API key regen.

### 4.7 User, Security & Compliance
*   Role-based ACL, 2FA (TOTP), password policy.
*   Audit trail table for CRUD on modems, SIMs, rates, users.
*   GDPR/data-retention config: purge CDR after N months, optional audio recording.
*   TLS everywhere, signed JWT with 12 h TTL, refresh tokens, IP allow list.

### 4.8 Observability & DevOps
*   Metrics exporter: modem_online, sim_balance, asr, acd, cpu, mem.
*   Alertmanager rules: low SIM credit, high CPU, high failed calls.
*   Structured JSON logs shipped to Loki/ELK.
*   Automatic nightly Postgres dump + S3.
*   Dockerfile per service; docker-compose.yml for single-node; Helm chart for HA.
*   Blue/green deployment script (systemd + symlink swap).

### 4.9 APIs
*   REST (JSON) and Web-Socket push. gRPC optional.
*   OpenAPI spec auto-generated.
*   Rate-limit 100 req/s per customer with Redis token-bucket.

## Non-Functional
*   **Throughput:** 60 concurrent calls, 20 000 CDR/day.
*   **Latency:** < 250 ms call setup internal, < 1 s UI P90.
*   **Availability:** 99.9 % monthly.
*   **Security:** pass OWASP top-10, AMI reachable only via localhost.

## Milestones
*   **M1** – Tech MVP (modem detect, make/receive call, basic UI) – 2 w
*   **M2** – SIM inventory + auto-recharge – 2 w
*   **M3** – LCR + per-second billing – 3 w
*   **M4** – Anti-spam + alerting – 2 w
*   **M5** – Customer portal + payments – 3 w
*   **M6** – Docker/K8s, HA, full monitoring – 2 w

## Risks & Mitigation
*   **Rapid SIM blocking** – humanisation & call-duration shaper.
*   **USB hub failures** – dual-PSU + spare hubs on site.
*   **Legal changes** – modular compliance layer.
*   **Payment disputes** – automated CDR evidence exports.

*This v0.2 PRD outlines the major aspects for operating, monetising and securing a VoIP-termination gateway. It should be used to guide development and track progress against milestones.*
