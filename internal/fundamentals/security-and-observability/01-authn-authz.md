# AuthN / AuthZ

---

## Authentication vs Authorisation

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client
    participant AuthN as 🔑 Authentication<br/>(JWT, OAuth 2.0, API Key)
    participant AuthZ as 🛡️ Authorization<br/>(RBAC, ABAC, OPA)
    participant API as ⚙️ API

    Client->>AuthN: Verify caller identity
    AuthN-->>Client: ✅ Verified identity claim

    Client->>AuthZ: Check permissions for verified identity
    AuthZ-->>Client: ✅ Allow / ❌ Deny

    alt Authorized
        Client->>API: Request allowed
        API-->>Client: 200 OK
    else Not Authorized
        API-->>Client: 403 Forbidden / 404 Not Found
    end
```

> Authentication must succeed **before** authorisation can be evaluated. You cannot check permissions for an unknown caller.

---

## JWT: A Self-Contained Identity Token

**JWT Structure: `header.payload.signature`**

| Component | Purpose | Example |
|-----------|---------|---------|
| **📋 HEADER** | Algorithm & token type | <code>{<br/>&nbsp;&nbsp;"alg":&nbsp;&nbsp;"HS256",<br/>&nbsp;&nbsp;"typ": "JWT"<br/>}</code> |
| **📦 PAYLOAD** | Identity claims | <code>{<br/>&nbsp;&nbsp;"sub": "user_alice",<br/>&nbsp;&nbsp;"roles": ["admin"],<br/>&nbsp;&nbsp;"exp": 1735689600<br/>}</code> |
| **🔐 SIGNATURE** | Proof of authenticity | `HMAC-SHA256(header.payload, secret_key)` |

**Encoded:** `eyJhbGc...Ew.eyJzdWI...NX0.SflKxw...`

```mermaid
sequenceDiagram
    participant Client as 📱 Client
    participant API as 🚪 API
    participant Handler as ⚙️ Handler

    Client->>API: GET /api/v1/accounts<br/>Authorization: Bearer eyJ...

    API->>API: Validate signature<br/>Check expiry<br/>Extract claims

    alt Valid & Not Expired
        API->>Handler: User verified<br/>userID: user_alice<br/>roles: [admin]
        Handler-->>Client: 200 OK {data...}
    else Invalid or Expired
        API-->>Client: 401 Unauthorized
    end
```


> The signature guarantees the payload was **not tampered with**. The expiry prevents stolen tokens from being valid forever.

---

## JWT Validation Flow

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client
    participant MW as 🔐 JWT Middleware
    participant CTX as 🗂️ context.Context
    participant Handler as ⚙️ Handler

    Client->>MW: GET /api/v1/accounts<br/>Authorization: Bearer eyJ...

    MW->>MW: strings.TrimPrefix(header, "Bearer ")
    MW->>MW: Verify signature
    MW->>MW: Check exp not in past

    alt Invalid or expired
        MW-->>Client: 401 Unauthorized {"code":"INVALID_TOKEN"}
    else Valid
        MW->>MW: Extract claims: userID, roles, exp
        MW->>CTX: context.WithValue(r.Context(), identityKey, UserIdentity{...})
        MW->>Handler: ServeHTTP with enriched context
        Handler->>CTX: identity := r.Context().Value(identityKey)
        CTX-->>Handler: UserIdentity{ID:"user_alice", Roles:["admin"]}
        Handler-->>Client: 200 OK {"data"...}
    end

    Note over Client,Handler: No session lookup. No database call.<br/>Identity travels IN the request.
```

---

## Stateless: Why JWT Scales

```mermaid
graph TB
    subgraph Stateful["❌ SESSION-BASED — Does NOT Scale"]
        SC["📱 Client"] -->|"session_id=abc"| SA["⚙️ Server A<br/>stores session in memory"]
        SC2["📱 Client"] -->|"session_id=abc"| SB["⚙️ Server B<br/>❓ unknown session → FAIL"]
    end

    subgraph Stateless["✅ JWT-BASED — Scales Freely"]
        JC["📱 Client<br/>carries JWT"] -->|"Bearer eyJ..."| JA["⚙️ Server A<br/>validates signature ✅"]
        JC -->|"Bearer eyJ..."| JB["⚙️ Server B<br/>validates signature ✅"]
        JC -->|"Bearer eyJ..."| JN["⚙️ Server N<br/>validates signature ✅"]
    end
```

---

## Go Typed Context Keys — No String Collisions

### ❌ BAD — String Keys (Collisions & Accidental Overwrites)

```go
// Package: auth
ctx = context.WithValue(ctx, "userID", "user_123")

// Package: logging (different developer)
ctx = context.WithValue(ctx, "userID", "log_456")  // OOPS! Overwrite!

// Package: auth tries to read
userID := ctx.Value("userID")  // Gets "log_456", not "user_123" ❌
```

**Problems:**
- Any package can use the string key `"userID"`
- Name collisions are invisible — one package silently overwrites another
- Type assertion to `interface{}` requires runtime conversion with risk of panics
- No compiler protection — mistakes discovered only at runtime

---

### ✅ GOOD — Private Typed Keys (No Collisions Possible)

```go
// Package: auth (internal package)
package auth

type contextKey struct{}
var identityKey = contextKey{}  // Unexported — only this package can use it

// Set identity in context
func WithIdentity(ctx context.Context, identity UserIdentity) context.Context {
    return context.WithValue(ctx, identityKey, identity)
}

// Get identity from context
func IdentityFromContext(ctx context.Context) (UserIdentity, bool) {
    identity, ok := ctx.Value(identityKey).(UserIdentity)
    return identity, ok
}
```

**Advantages:**
- The type `contextKey` is **private** (unexported) — only this package can instantiate it
- No other package can accidentally use the same key — impossible to collide
- Type safety — Go's compiler prevents type mismatches
- Clear API — consumers call `IdentityFromContext()` instead of guessing string keys

> A private unexported type as the context key is **impossible to collide with** — only that package can set or read it.

---

## OAuth 2.0: Delegated Access

```mermaid
sequenceDiagram
    autonumber
    participant User as 👤 User
    participant App as 📱 Bank App
    participant Auth as 🔑 Auth Server<br/>(Okta / Azure AD)
    participant API as 🚪 Bank API

    User->>App: Login
    App->>Auth: Redirect /authorize?client_id=...&scope=accounts:read
    Auth->>User: Login + consent screen
    User->>Auth: Approves
    Auth-->>App: Authorization code

    App->>Auth: POST /token {code, client_secret}
    Auth-->>App: access_token (JWT) + refresh_token

    App->>API: GET /api/v1/accounts<br/>Authorization: Bearer <access_token>
    API->>API: Validate JWT against Auth Server public key
    API-->>App: 200 OK {"data": [...]}

    Note over App,API: The API trusts the Auth Server's signature.<br/>No user password ever touches the API.
```

---

## RBAC vs ABAC

```mermaid
graph TB
    subgraph RBAC["Role-Based Access Control"]
        direction LR
        RB1["User has roles: [admin, viewer]"]
        RB2["Role has permissions: admin → can_delete"]
        RB3["Simple · Fast · Easy to audit"]
        RB1 --> RB2 --> RB3
    end

    subgraph ABAC["Attribute-Based Access Control"]
        direction LR
        AB1["Policy uses ANY attribute:<br/>user.department · resource.owner<br/>time.hour · request.ip"]
        AB2["allow if user.dept == resource.dept<br/>AND time.now < 18:00<br/>AND resource.classification != 'secret'"]
        AB3["Fine-grained · Powerful · Complex"]
        AB1 --> AB2 --> AB3
    end

    CHOOSE{"Which to use?"}
    RBAC --> CHOOSE
    ABAC --> CHOOSE
    CHOOSE -->|"Simple roles, clear boundaries"| RBAC
    CHOOSE -->|"Complex, contextual decisions"| ABAC
```

---

## Token Revocation: The Stateless Edge Case

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

    Note over Client,Redis: On next request with the same token

    Client->>API: GET /accounts<br/>Authorization: Bearer eyJ...
    API->>API: JWT signature valid ✅
    API->>Redis: GET revoked:jti_xyz
    Redis-->>API: EXISTS → token is revoked
    API-->>Client: 401 Unauthorized

    Note over Client,Redis: Revocation state lives in shared store, not server memory.
```

> The server remains stateless per-instance. Revocation state lives in the **shared store**, not in server memory.

---

## HTTP Status Codes for Auth Failures

```mermaid
graph TB
    subgraph Codes["Auth Error Responses"]
        C401["**401 Unauthorized**<br/>No valid credentials provided<br/>or token expired / invalid<br/>→ Client must re-authenticate"]
        C403["**403 Forbidden**<br/>Valid identity, insufficient permissions<br/>→ Access denied, don't retry same request"]
        C404["**404 Not Found**<br/>Use instead of 403 when existence<br/>must not be revealed to caller"]
    end

    RULE["Rule: 401 = identity problem<br/>403 = permission problem<br/>404 = existence hiding"]
    Codes --> RULE
```

> Return `404` instead of `403` when revealing that a resource *exists* would itself be a security leak.
