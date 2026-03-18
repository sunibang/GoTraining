# Distributed Tracing

---

## The Three Pillars of Observability

```mermaid
graph TB
    OBS["🔭 OBSERVABILITY<br/>Can you understand what your system is doing?"]

    LOGS["📋 LOGS<br/>Discrete events<br/>'Payment failed for acc_001'<br/>→ What happened"]
    METRICS["📊 METRICS<br/>Aggregated numbers over time<br/>error_rate=0.02 · p99_latency=340ms<br/>→ How much / how fast"]
    TRACES["🔗 TRACES<br/>The lifecycle of one request<br/>across multiple services<br/>→ Where time is spent"]

    OBS --> LOGS
    OBS --> METRICS
    OBS --> TRACES

    LOGS --- Q1["Query: find all errors in last 5m"]
    METRICS --- Q2["Alert: error rate > 1%"]
    TRACES --- Q3["Debug: why is checkout slow?"]
```

> Logs tell you **what**. Metrics tell you **how much**. Traces tell you **where**.

---

## Anatomy of a Distributed Trace

```mermaid
graph LR
    subgraph Request["Single User Request: POST /payments — trace_id: 4bf92f..."]
        direction TB
        ROOT["🔷 Root Span<br/>POST /payments<br/>Total: 210ms"]
        SPACE1[" "]
        SPACE2[" "]

        ROOT --> AUTH["🔷 Span: validate JWT<br/>3ms"]
        ROOT --> POLICY["🔷 Span: OPA policy check<br/>5ms"]
        ROOT --> SVC["🔷 Span: PaymentService.Process<br/>195ms"]

        SVC --> DB["🔷 Span: INSERT payments<br/>db: postgres<br/>180ms ← bottleneck!"]
        SVC --> NOTIFY["🔷 Span: NotificationService<br/>12ms"]
    end

    style SPACE1 fill:none,stroke:none
    style SPACE2 fill:none,stroke:none
```

> Without tracing, you see a 210ms P99. With tracing, you see the DB insert is 180ms. **Now you know what to fix.**

---

## OpenTelemetry: The Standard

```mermaid
graph TB
    subgraph OTel["OpenTelemetry SDK"]
        INSTR["📦 Instrumentation<br/>Wrap HTTP handlers<br/>Wrap DB calls · Wrap outbound HTTP"]
        PROP["🔗 Context Propagation<br/>traceparent header injected<br/>into every outbound call"]
        EXP["📤 Exporter<br/>OTLP → Jaeger / Tempo / Datadog / X-Ray"]
        INSTR --> PROP --> EXP
    end

    subgraph Services["Your Services"]
        S1["⚙️ API Gateway"]
        S2["⚙️ Payment Service"]
        S3["⚙️ Notification Service"]
        S1 -->|"traceparent: 00-4bf92f..."| S2
        S2 -->|"traceparent: 00-4bf92f..."| S3
    end

    Services --> OTel
    EXP --> BACKEND["🔍 Trace Backend<br/>Jaeger / Grafana Tempo"]
```

> One `trace_id` ties together spans from every service in the call chain.

---

## Trace Propagation Across Services

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client
    participant GW as ⚙️ API Gateway
    participant Pay as ⚙️ Payment Service
    participant Notify as ⚙️ Notification Service

    Client->>GW: POST /payments
    GW->>GW: Create root span<br/>trace_id: 4bf92f3577b34da6

    GW->>Pay: POST /internal/payments<br/>traceparent: 00-4bf92f3577b34da6-00f067aa0ba902b7-01
    Pay->>Pay: Create child span<br/>parent: 00f067aa0ba902b7

    Pay->>Notify: POST /notifications/send<br/>traceparent: 00-4bf92f3577b34da6-b7ad6b7169203331-01
    Notify->>Notify: Create child span<br/>parent: b7ad6b7169203331
    Notify-->>Pay: 202 Accepted

    Pay-->>GW: 201 Created
    GW-->>Client: 201 Created

    Note over GW,Notify: All spans share trace_id: 4bf92f...<br/>Full request journey reconstructed in Jaeger.
```

---

## Prometheus: The Four Golden Signals

```mermaid
graph TB
    subgraph Golden["The Four Golden Signals"]
        LATENCY["⏱️ LATENCY<br/>How long requests take<br/>Histogram: http_request_duration_seconds"]
        TRAFFIC["📈 TRAFFIC<br/>How much demand<br/>Counter: http_requests_total"]
        ERRORS["❌ ERRORS<br/>Rate of failed requests<br/>Counter: http_errors_total"]
        SATURATION["💾 SATURATION<br/>How full resources are<br/>Gauge: go_goroutines"]
    end

    LATENCY ~~~ TRAFFIC
    TRAFFIC ~~~ ERRORS
    ERRORS ~~~ SATURATION
```

---

## Prometheus Metric Types

```mermaid
graph TB
    COUNTER["**Counter**<br/>Only goes up · Resets on restart<br/>Use: requests · errors · payments processed<br/>http_requests_total{method='GET',status='200'}"]
    GAUGE["**Gauge**<br/>Goes up and down<br/>Use: active connections · queue depth · memory<br/>active_connections 42"]
    HISTOGRAM["**Histogram**<br/>Samples observations into buckets<br/>Use: request duration · payload size<br/>http_duration_seconds_bucket{le='0.1'} 4523"]

    COUNTER ~~~ GAUGE
    GAUGE ~~~ HISTOGRAM
```

---

## Metrics + Tracing Middleware Flow

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client
    participant MW as 📊 Metrics Middleware
    participant TMW as 🔗 Tracing Middleware
    participant Handler as ⚙️ Handler
    participant Prom as 📈 Prometheus

    Client->>MW: POST /api/v1/payments<br/>traceparent: 00-4bf92f...

    MW->>MW: Start timer
    MW->>TMW: Forward request

    TMW->>TMW: Extract traceparent header
    TMW->>TMW: Create child span
    TMW->>TMW: Inject span into context

    TMW->>Handler: Forward with trace context
    Handler->>Handler: span := trace.SpanFromContext(r.Context())
    Handler->>Handler: span.AddEvent("payment.validated")
    Handler-->>TMW: 201 Created

    TMW->>TMW: span.End() — record duration
    TMW-->>MW: Pass response

    MW->>Prom: http_requests_total{method=POST,status=201}.Inc()
    MW->>Prom: http_request_duration_seconds.Observe(elapsed)
    MW-->>Client: 201 Created

    Note over Client,Prom: Prometheus scrapes /metrics every 15s.<br/>Trace exported async to Jaeger.
```

---

## The Full Observability Stack

```mermaid
graph TB
    APP["⚙️ Go API"]

    APP -->|"stdout JSON lines"| LOGS["📋 Log Aggregator<br/>Datadog / CloudWatch"]
    APP -->|"GET /metrics"| PROM["📊 Prometheus<br/>Scrapes every 15s"]
    APP -->|"OTLP gRPC"| JAEGER["🔗 Jaeger / Grafana Tempo<br/>Distributed Traces"]

    PROM --> GRAF["📈 Grafana<br/>Dashboards + Alerts"]
    LOGS --> GRAF
    JAEGER --> GRAF
```

> One Grafana dashboard — logs, metrics, and traces linked by `trace_id`. Full system visibility.

---

## Correlating Logs and Traces

```mermaid
sequenceDiagram
    autonumber
    participant MW as 🔀 Logging + Tracing Middleware
    participant CTX as 🗂️ context.Context
    participant Handler as ⚙️ Handler
    participant Log as 📋 slog.Logger

    MW->>MW: Extract / create trace_id from traceparent header
    MW->>MW: Create request_id = uuid
    MW->>CTX: logger := slog.Default().With(<br/>  "trace_id", traceID,<br/>  "request_id", requestID,<br/>)
    MW->>Handler: ServeHTTP with context

    Handler->>CTX: logger := loggerFromContext(r.Context())
    Handler->>Log: logger.Info("payment validated", "payment_id", id)

    Note over Log: {trace_id:4bf92f..., request_id:req_abc,<br/>msg:payment validated, payment_id:pay_xyz}

    Note over MW,Log: Same trace_id in both your logs and Jaeger.<br/>Click trace → pivot to logs. Click log → pivot to trace.
```

> Inject `trace_id` into every log line. Logs and traces become **navigable together** in Grafana.
