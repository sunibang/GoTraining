# Structured Logging with slog

---

## Unstructured vs Structured Logs

```mermaid
graph TB
    subgraph Unstructured["❌ UNSTRUCTURED"]
        U1["2024-01-15 10:23:44 INFO user alice fetched account acc_001 in 12ms"]
        U2["2024-01-15 10:23:45 ERROR failed to fetch account acc_002: timeout"]
        U3["To find errors: grep 'ERROR' | awk '{print $5}' | sort | uniq -c"]
        U4["At 10M logs/day: good luck"]
        U1 --> U2 --> U3 --> U4
    end

    subgraph Structured["✅ STRUCTURED"]
        S1["{time:2024-01-15T10:23:44Z, level:INFO,<br/>msg:account fetched, user:alice,<br/>account_id:acc_001, latency_ms:12}"]
        S2["Query: level=ERROR AND latency_ms > 1000"]
        S3["Dashboard: P99 latency by endpoint"]
        S1 --> S2 --> S3
    end

    U4 ~~~ S1

```

> Unstructured: ❌ Impossible to Query at Scale <br/> Structured: ✅ Machine-Readable, Instantly Queryable

---

## log/slog: Go's Native Structured Logger

```mermaid
graph TB
    subgraph API["slog API"]
        direction LR
        INFO["slog.Info(msg, k, v, ...)"]
        WARN["slog.Warn(msg, k, v, ...)"]
        ERROR["slog.Error(msg, k, v, ...)"]
        DEBUG["slog.Debug(msg, k, v, ...)"]
    end

    subgraph Handlers["Pluggable Handlers"]
        TEXT["slog.NewTextHandler<br/>→ key=value (development)"]
        JSON["slog.NewJSONHandler<br/>→ {key:value} (production)"]
    end

    subgraph Sinks["Log Aggregators"]
        direction LR
        DD["Datadog"]
        ELK["Elasticsearch"]
        CW["CloudWatch"]
        SPLUNK["Splunk"]
    end

    API --> TEXT
    API --> JSON
    JSON -->|"JSON lines on stdout"| Sinks
```

> Switch from text to JSON with **one line** at startup. Zero changes to log call sites.

---

## JSON Handler: Production Configuration

```mermaid
graph LR
    INIT["slog.SetDefault(<br/>  slog.New(<br/>    slog.NewJSONHandler(<br/>      os.Stdout,<br/>      &slog.HandlerOptions{<br/>        Level: slog.LevelInfo,<br/>      },<br/>    ),<br/>  ),<br/>)"]

    OUTPUT["{time:2024-01-15T10:23:44.123Z,<br/>level:INFO,<br/>msg:payment processed,<br/>payment_id:pay_xyz,<br/>amount:250.00,<br/>user_id:user_alice}"]

    INIT -->|"stdout"| OUTPUT
    OUTPUT -->|"parsed by"| AGG["📊 Log Aggregator<br/>Datadog / CloudWatch / ELK"]
```

---

## Log Levels: When to Use Each

```mermaid
graph TB
    DEBUG["🔵 DEBUG<br/>Detailed internal state.<br/>Never in production.<br/>slog.Debug('cache miss', 'key', k)"]
    INFO["🟢 INFO<br/>Normal operational events.<br/>Request received, payment processed.<br/>slog.Info('payment processed', 'id', id)"]
    WARN["🟡 WARN<br/>Unexpected but handled.<br/>Retry succeeded, degraded mode.<br/>slog.Warn('upstream slow, using cache')"]
    ERROR["🔴 ERROR<br/>Something failed and needs attention.<br/>Database down, unhandled error.<br/>slog.Error('db query failed', 'err', err)"]

    DEBUG --> INFO --> WARN --> ERROR
```

---

## Logger with Context: Propagating Fields

```mermaid
sequenceDiagram
    autonumber
    participant MW as 🔀 Logging Middleware
    participant CTX as 🗂️ context.Context
    participant Handler as ⚙️ Handler
    participant Log as 📋 slog.Logger

    MW->>MW: Generate request_id = uuid
    MW->>CTX: Store slog.Logger with request_id in context
    MW->>Handler: ServeHTTP with enriched context

    Handler->>CTX: logger := loggerFromContext(r.Context())
    Handler->>Log: logger.Info("fetching account", "account_id", id)

    Note over Log: {time:..., level:INFO, msg:fetching account,<br/>request_id:req_abc, account_id:acc_001}

    Handler-->>MW: Response written
    MW->>Log: logger.Info("request complete", "status", 200, "latency_ms", 14)

    Note over MW,Log: Every log line from this request carries request_id.<br/>Find all logs for one request instantly.
```

---

## Structured Log Fields: What to Always Include

```mermaid
graph TB

    subgraph OnError["On Error"]
        E1["error — err.Error() string"]
        E2["stack — optional, for panics only"]
    end

    subgraph PerRequest["Per-Request"]
        R1["method — GET, POST, etc."]
        R2["path — /api/v1/accounts"]
        R3["status — 200, 404, 500"]
        R4["latency_ms — response time"]
        R5["user_id — authenticated identity"]
    end

    subgraph Always["Always Include"]
        A1["time — ISO 8601 timestamp"]
        A2["level — DEBUG / INFO / WARN / ERROR"]
        A3["msg — human-readable description"]
        A4["trace_id — correlate across services"]
        A5["request_id — correlate within a request"]
    end
```

---

## Logger-Per-Request Pattern

```mermaid
graph TB
    subgraph Middleware["Logging Middleware"]
        M1["reqLogger := slog.Default().With(<br/>  'request_id', requestID,<br/>  'method', r.Method,<br/>  'path', r.URL.Path,<br/>)"]
        M2["ctx := context.WithValue(r.Context(), logKey, reqLogger)"]
        M3["next.ServeHTTP(w, r.WithContext(ctx))"]
        M1 --> M2 --> M3
    end

    subgraph Handler["Handler"]
        H1["logger := loggerFromContext(r.Context())"]
        H2["logger.Info('account fetched', 'account_id', id)"]
        H3["// {request_id:req_abc, account_id:acc_001, ...}"]
        H1 --> H2 --> H3
    end

    Middleware --> Handler
```

> Every `With(...)` call returns a **new** logger — the original is unchanged. Fields accumulate down the call chain.

---

## slog vs fmt.Println

```mermaid
graph TB
    subgraph FMT["❌ fmt.Println — Never in Production"]
        F1["fmt.Println('payment processed: ' + id)"]
        F2["Unstructured string"]
        F3["No level, no timestamp, no fields"]
        F1 --> F2 --> F3
    end

    subgraph SLOG["✅ slog — Production Standard"]
        S1["slog.Info('payment processed', 'id', id, 'amount', amt)"]
        S2["JSON: {time:..., level:INFO, msg:..., id:pay_xyz, amount:250}"]
        S3["Queryable · Alertable · Dashboard-ready"]
        S1 --> S2 --> S3
    end
```

> `log/slog` is in the Go standard library since **1.21**. No external dependencies needed.
