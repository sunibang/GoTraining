# Idempotency

---

## What Is Idempotency?

```mermaid
graph TB
    DEF["**Idempotent operation:**<br/>Calling it once produces the same result<br/>as calling it N times"]

    subgraph Idempotent["✅ Idempotent"]
        I1["GET  /accounts/acc_123<br/>→ same account every time"]
        I2["PUT  /accounts/acc_123 body={name: 'Savings'}<br/>→ same final state"]
        I3["DELETE /accounts/acc_123<br/>→ deleted whether called once or twice"]
    end

    subgraph NotIdempotent["❌ NOT Idempotent by default"]
        N1["POST /payments body={amount: 250.00}<br/>→ creates a NEW payment each call"]
        N2["POST /notifications/send<br/>→ sends a NEW email each call"]
    end

    DEF --> Idempotent
    DEF --> NotIdempotent
```

> If a retry is safe and produces no extra side effects, the operation is idempotent.

---

## HTTP Methods and Idempotency

```mermaid
graph LR
    GET["**GET**<br/>📖 Read-only<br/>✅ Idempotent<br/>✅ Cacheable"]
    HEAD["**HEAD**<br/>📋 Metadata only<br/>✅ Idempotent"]
    PUT["**PUT**<br/>🔄 Full replace<br/>✅ Idempotent"]
    DELETE["**DELETE**<br/>🗑️ Remove<br/>✅ Idempotent"]
    PATCH["**PATCH**<br/>✏️ Partial update<br/>⚠️ Depends on implementation"]
    POST["**POST**<br/>➕ Create / action<br/>❌ NOT idempotent<br/>Use Idempotency-Key header"]
    OPTIONS["**OPTIONS**<br/>🔍 Preflight / CORS<br/>✅ Idempotent"]

    GET ~~~ HEAD ~~~ PUT ~~~ DELETE ~~~ PATCH ~~~ POST ~~~ OPTIONS
```

---

## The Double-Charge Problem

```mermaid
sequenceDiagram
    autonumber
    participant App as 📱 Mobile App
    participant API as 🚪 API
    participant DB as 🐘 Database

    Note over App,DB: ❌ WITHOUT idempotency key

    App->>API: POST /payments {"amount": 250.00}
    API->>DB: INSERT INTO payments ...
    DB-->>API: payment_id = pay_001
    Note over App,API: ⏱️ Network timeout — app never receives response

    App->>API: POST /payments {"amount": 250.00}  ← retry
    API->>DB: INSERT INTO payments ...
    DB-->>API: payment_id = pay_002
    API-->>App: 201 Created {"id": "pay_002"}

    Note over App,DB: 💥 Customer charged TWICE
```

---

## The Fix: Idempotency-Key Header

```mermaid
sequenceDiagram
    autonumber
    participant App as 📱 Mobile App
    participant API as 🚪 API
    participant Cache as ⚡ Redis / DB
    participant DB as 🐘 Database

    Note over App,DB: ✅ WITH Idempotency-Key

    App->>API: POST /payments {"amount": 250.00}<br/>Idempotency-Key: key_abc123
    API->>Cache: Check key_abc123 → not found
    API->>DB: INSERT INTO payments ...
    DB-->>API: payment_id = pay_001
    API->>Cache: Store key_abc123 → {"id":"pay_001","status":201}
    API-->>App: 201 Created {"id": "pay_001"}

    Note over App,API: ⏱️ Network timeout — app retries

    App->>API: POST /payments {"amount": 250.00}<br/>Idempotency-Key: key_abc123  ← SAME key
    API->>Cache: Check key_abc123 → FOUND
    API-->>App: 201 Created {"id": "pay_001"}  ← cached response

    Note over App,DB: ✅ No duplicate. Same result. Safe to retry.
```

> The key must be generated **client-side** before the request is sent, and reused on every retry of the same logical operation.

---

## Idempotency Key Lifecycle

```mermaid
graph TB
    GEN["📱 Client generates key<br/>UUID or similar<br/>e.g. key_abc123"]

    FIRST["1️⃣ First request<br/>Key not in store → process normally<br/>Store key + response"]

    RETRY["🔄 Retry (timeout / 5xx)<br/>Key already in store → return cached response<br/>No side effects"]

    EXPIRE["⏱️ Key expires<br/>After TTL (e.g. 24 hours)<br/>Client must generate new key for new operations"]

    ERROR["❌ Key conflict<br/>Same key, different body → 422 Unprocessable<br/>Client bug: one key per logical operation"]

    GEN --> FIRST --> RETRY --> EXPIRE
    FIRST --> ERROR
```

---

## DELETE Idempotency: Already Gone Is Fine

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client
    participant API as 🚪 API
    participant DB as 🐘 Database

    Client->>API: DELETE /accounts/acc_123
    API->>DB: DELETE WHERE id = acc_123
    DB-->>API: 1 row deleted
    API-->>Client: 204 No Content

    Note over Client,API: ⏱️ Client times out, retries

    Client->>API: DELETE /accounts/acc_123  ← retry
    API->>DB: DELETE WHERE id = acc_123
    DB-->>API: 0 rows deleted (already gone)
    API-->>Client: 204 No Content  ← same response

    Note over Client,API: ✅ Resource is gone either way. Safe.
```

> Return `204` (not `404`) on a repeat DELETE. The desired state — "this resource must not exist" — has been achieved.

---

## PATCH: Conditional Idempotency

```mermaid
graph TB

    subgraph Idem["✅ PATCH — Idempotent"]
        Q1["PATCH /accounts/acc_123<br/>body: {'name': 'New Savings Account'}"]
        Q2["Call 1: name updated"]
        Q3["Call 2: same name, same result"]
        Q1 --> Q2 --> Q3
    end

    subgraph NonIdem["❌ PATCH — NOT Idempotent"]
        P1["PATCH /accounts/acc_123<br/>body: {'balance': {'increment': 100}}"]
        P2["Call 1: balance = 500 → 600"]
        P3["Call 2: balance = 600 → 700"]
        P1 --> P2 --> P3
    end

    Q3 ~~~ P1
```

> A PATCH that **sets an absolute value** is idempotent. A PATCH that **increments** is not — use an Idempotency-Key.

---

## Idempotency at Financial Institutions: What Gets a Key

```mermaid
graph LR
    MUST["🔴 **MUST use Idempotency-Key**"]
    SHOULD["🟡 **SHOULD use Idempotency-Key**"]
    NOTNEED["🟢 **Does NOT need Idempotency-Key**"]

    MUST --> M1["POST /payments"]
    MUST --> M2["POST /transfers"]
    MUST --> M3["POST /direct-debits"]
    MUST --> M4["POST /notifications/send"]

    SHOULD --> S1["POST /accounts (onboarding)"]
    SHOULD --> S2["POST /loans/apply"]

    NOTNEED --> N1["GET (all read operations)"]
    NOTNEED --> N2["PUT /accounts/{id}"]
    NOTNEED --> N3["DELETE /accounts/{id}"]
```

> Every operation with a **financial side effect** must be idempotent. Retries happen. Networks fail. Design for it.
