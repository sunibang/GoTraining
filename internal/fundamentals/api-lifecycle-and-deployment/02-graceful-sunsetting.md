# Graceful Sunsetting

---

## Deprecation Is a Process, Not an Event

```mermaid
graph TB
    NOW["📅 TODAY<br/>Announce deprecation<br/>Inject Deprecation + Sunset headers"]
    GRACE["⏳ GRACE PERIOD (e.g. 6 months)<br/>Both old and new endpoints live<br/>Clients migrate at their own pace"]
    SUNSET["🌅 SUNSET DATE<br/>Old endpoint returns 410 Gone<br/>Link header points to replacement"]
    DONE["✅ DONE<br/>Old route removed from codebase<br/>No version debt"]

    NOW --> GRACE --> SUNSET --> DONE
```

> Clients — and their automated tooling — need explicit signals. Headers are that signal.

---

## Standard Deprecation Headers

```mermaid
graph TB
    subgraph Response["HTTP Response — deprecated endpoint"]
        SPACE1["  "]
        H1["Deprecation: true<br/>Machine-readable flag for clients and SDK tooling"]
        H2["Sunset: Tue, 01 Jul 2025 00:00:00 GMT<br/>Exact date the endpoint stops responding"]
        H3["Link: &lt;https://docs.example.com/migration&gt;; rel=sunset<br/>Points to migration guide — automated tools follow this"]
    end

    subgraph Tooling["What Automated Tooling Does With These Headers"]
        SPACE2["  "]
        T1["🤖 SDK generators read Deprecation: true<br/>Emit compile-time warnings"]
        T2["📊 API gateways track sunset dates<br/>Alert ops teams before removal"]
        T3["🔗 Link rel=sunset enables<br/>one-click migration guide discovery"]
    end

    style SPACE1 fill:none,stroke:none
    style SPACE2 fill:none,stroke:none
    Response --> Tooling
```

---

## Deprecation Middleware Flow

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client
    participant MW as 🔧 deprecationMiddleware
    participant Handler as ⚙️ Old v1 Handler
    participant Docs as 📄 Migration Docs

    Note over Client,Docs: Old v1 endpoint is still functional — but clearly signalled

    Client->>MW: GET /api/v1/accounts/acc_001
    MW->>Handler: Forward request
    Handler-->>MW: 200 OK {"id":"acc_001",...}

    MW->>MW: Inject headers<br/>Deprecation: true<br/>Sunset: Tue, 01 Jul 2025 00:00:00 GMT<br/>Link: https://docs.example.com/migration?rel=sunset

    MW-->>Client: 200 OK + deprecation headers

    Note over Client,Docs: Client SDK detects Deprecation: true — logs migration warning

    Client->>Docs: Follow Link header → read migration guide
```

---

## Middleware as a Reusable Decorator

```mermaid
graph TB
    subgraph Deprecated["🪦 Deprecated Endpoints"]
        OLD1["GET /api/v1/accounts/{id}"]
        OLD2["POST /api/v1/accounts"]
        OLD3["GET /api/v1/payments"]
    end

    subgraph Active["✅ Active Endpoints"]
        NEW1["GET /api/v2/accounts/{id}"]
        NEW2["POST /api/v2/accounts"]
        NEW3["GET /api/v2/payments"]
    end

    MW["🔧 deprecationMiddleware(sunsetDate, docURL)<br/>Wraps any handler — zero coupling to business logic"]

    OLD1 --> MW
    OLD2 --> MW
    OLD3 --> MW

    MW -->|"adds headers then delegates"| HANDLER["⚙️ Original Handler<br/>Unchanged — still serves responses"]
```

> The middleware is a thin decorator. The handler is unaware it is being deprecated. Separation of concerns.

---

## After Sunset: 410 Gone

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client (not yet migrated)
    participant API as 🚪 API

    Note over Client,API: After sunset date — endpoint removed

    Client->>API: GET /api/v1/accounts/acc_001
    API-->>Client: 410 Gone<br/>code: ENDPOINT_REMOVED<br/>message: Use /api/v2/accounts<br/>Link: docs.example.com/migration

    Note over Client,API: 410 is permanent — do not retry, migrate now
```

> `410 Gone` is the correct status after sunset. `404 Not Found` implies it might come back. `410` is permanent.

---

## Tracking Who Is Still on v1

```mermaid
graph TB
    subgraph Signals["📊 Migration Progress Signals"]
        S1["📈 v1 request volume trending to zero?<br/>Sunset is safe to proceed"]
        S2["📊 v1 request volume still high?<br/>Extend the grace period — clients are blocked"]
        S3["🔍 Which clients are still calling v1?<br/>User-Agent + API key logs reveal who to contact"]
    end

    subgraph Actions["⚙️ Actions Before Removal"]
        A1["Email API key owners still hitting v1"]
        A2["Post deprecation notices in developer portal"]
        A3["Set hard cutoff — communicate clearly"]
    end

    Signals --> Actions
```

> Never remove a version without first checking traffic. Logs tell you when it is safe. Until then, keep the grace period open.
