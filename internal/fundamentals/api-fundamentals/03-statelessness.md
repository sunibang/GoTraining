# Statelessness

---

## What Is Statelessness?

```mermaid
graph TB
    DEF["**Stateless API:**<br/>Every request contains ALL the information<br/>needed to process it.<br/>The server remembers nothing between calls."]

    WHAT["Every request carries:<br/>🪪 Identity (JWT / API key)<br/>🎯 Intent (method + URL)<br/>📦 Data (headers + body)<br/>🔒 Auth (signature / token)"]

    DEF --> WHAT
```

> The server has no memory of you. Each request stands alone.

---

## Stateful vs Stateless

```mermaid
graph TB
    subgraph Stateful["❌ STATEFUL — Session on Server"]
        SC1["📱 Client: Login"]
        SS1["💾 Server A<br/>stores session_id=abc"]
        SC2["📱 Client: GET /accounts"]
        SS2["❓ Server B<br/>No session → 401"]
        SC1 --> SS1
        SC2 --> SS2
        SS2 -->|"Sticky sessions required"| STICKY["📌 Load balancer<br/>must pin client to Server A"]
    end

    subgraph Stateless["✅ STATELESS — Context in Token"]
        SPACE1["  "]
        SLC["📱 Client sends JWT<br/>in EVERY request"]
        SLA["✅ Server A validates JWT → ok"]
        SLB["✅ Server B validates JWT → ok"]
        SLC2["✅ Server C validates JWT → ok"]
        SLC --> SLA
        SLC --> SLB
        SLC --> SLC2
    end

    style SPACE1 fill:none,stroke:none
```

> Stateless = any server instance can handle any request. No pinning, no shared memory.

---

## Scaling: The Stateless Advantage

```mermaid
graph TB
    C1["📱 Mobile"] --> LB
    C2["🖥️ Web"] --> LB
    C3["🤖 Agent"] --> LB

    LB["⚖️ **Load Balancer**<br/>Round-robin<br/>No sticky sessions needed"]

    LB --> I1["⚙️ Instance 1"]
    LB --> I2["⚙️ Instance 2"]
    LB --> I3["⚙️ Instance 3"]
    LB --> I4["⚙️ Instance N..."]

    I1 --> DB[("🐘 PostgreSQL")]
    I2 --> DB
    I3 --> DB
    I4 --> DB

    I1 --> CACHE[("⚡ Redis")]
    I2 --> CACHE
    I3 --> CACHE
    I4 --> CACHE
```

> Add an instance → more capacity, instantly. Remove an instance → no sessions lost.

---

## What Goes Where

```mermaid
graph TB

    subgraph NeverOnServer["❌ State That MUST NOT Live on Server Instance"]
        SPACE1["  "]
        N1["💾 Session memory"]
        N2["🛒 In-flight request context stored in memory"]
        N3["🔐 Per-user login state"]
    end

    subgraph InSharedStore["✅ State That Lives in a Shared Store"]
        SPACE2["  "]
        S1["💰 Account balances → PostgreSQL"]
        S2["⚡ Rate limit counters → Redis"]
        S3["🔁 Idempotency keys → Redis"]
        S4["🔐 Revoked tokens → Redis"]
    end

    subgraph InRequest["✅ State That Lives in the Request"]
        SPACE3["  "]
        R1["🪪 JWT — who you are, what roles you have"]
        R2["🔑 API Key — service identity"]
        R3["🎯 Path params — which resource"]
        R4["📦 Body — what to do"]
        R5["📋 Headers — content type, idempotency key"]
    end

    style SPACE1 fill:none,stroke:none
    style SPACE2 fill:none,stroke:none
    style SPACE3 fill:none,stroke:none
```

---

## JWT: The Stateless Identity Token

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client
    participant Auth as 🔐 Auth Service
    participant API as 🚪 API (any instance)

    Client->>Auth: POST /auth/token {username, password}
    Auth->>Auth: Validate credentials
    Auth->>Auth: Sign JWT {sub, roles, exp}
    Auth-->>Client: 200 OK {"token": "eyJ..."}

    Note over Client,API: Client stores token. Uses it on every request.

    Client->>API: GET /accounts<br/>Authorization: Bearer eyJ...
    API->>API: Decode JWT header + payload
    API->>API: Verify signature (shared secret / public key)
    API->>API: Check exp not in past
    API->>API: Extract sub, roles → attach to context
    API-->>Client: 200 OK {"accounts": [...]}

    Note over Client,API: No session lookup. No database hit. Stateless.
```

---

## Token Revocation: Handling the Stateless Edge Case

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client
    participant API as 🚪 API
    participant Redis as ⚡ Redis

    Note over Client,Redis: User logs out or token is compromised

    Client->>API: POST /auth/logout<br/>Authorization: Bearer eyJ...
    API->>Redis: SETEX revoked:jti_xyz TTL=token_exp
    API-->>Client: 200 OK

    Note over Client,Redis: On next request with same token

    Client->>API: GET /accounts<br/>Authorization: Bearer eyJ...
    API->>API: JWT signature valid ✅
    API->>Redis: GET revoked:jti_xyz
    Redis-->>API: EXISTS → token is revoked
    API-->>Client: 401 Unauthorized

    Note over Client,Redis: Revocation list in Redis — shared, fast, stateless per-instance.
```

> The server is still stateless per-instance. Revocation state lives in the **shared store**, not in server memory.

---

## Cache-Control: Stateless Caching

```mermaid
graph TB
    subgraph NotCacheable["❌ Not Cacheable"]
        N1["POST /payments<br/>Creates a new resource"]
        N2["DELETE /accounts/acc_123<br/>Mutates state"]
        N3["GET /accounts/acc_123/balance<br/>Cache-Control: no-store"]
    end
    
    subgraph Cacheable["✅ Cacheable — Stateless Responses"]
        SPACE1["  "]
        C1["GET /products<br/>Cache-Control: max-age=300"]
        C2["GET /accounts/acc_123<br/>Cache-Control: private, max-age=60"]
        C3["GET /rates/AUD-USD<br/>Cache-Control: public, max-age=30"]
    end

    style SPACE1 fill:none,stroke:none
```

> Stateless responses can be cached by CDNs and proxies. Mutable-state responses must opt out explicitly.

---

## Statelessness and Fault Tolerance

```mermaid
graph TB
    subgraph Before["❌ Instance Failure — Stateful"]
        BL["⚖️ Load Balancer"]
        BI1["⚙️ Instance 1<br/>💾 session for Alice"]
        BI2["💥 Instance 2 CRASHES<br/>💾 session for Bob — LOST"]
        BL --> BI1
        BL --> BI2
    end

    subgraph After["✅ Instance Failure — Stateless"]
        AL["⚖️ Load Balancer"]
        AI1["⚙️ Instance 1"]
        AI2["💥 Instance 2 CRASHES"]
        AI3["⚙️ Instance 3<br/>handles Bob's next request fine"]
        AL --> AI1
        AL --> AI2
        AL --> AI3
    end
```

> When a stateless instance crashes, the load balancer routes to another. **No user session is lost.**
