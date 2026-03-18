# API Design

## What is API Design?

API design is the process of defining how software components communicate — deciding what operations are exposed, how requests and responses are structured, and what contracts clients can rely on.

Good API design is intentional. It considers the consumer's experience first, makes common operations easy, and makes incorrect usage hard.

> **Bad design:** designing for the server's internal structure — exposing a database table as an endpoint, leaking implementation details in field names, or making clients stitch together multiple calls to do one logical thing.

---

## Core API Styles

### REST (Representational State Transfer)
- Resource-oriented: URLs identify *things*, HTTP verbs identify *actions*
- Stateless: each request contains all information needed to process it
- Human-readable, widely understood, excellent tooling support
- Best for: public APIs, browser clients, simple CRUD operations

### gRPC (Google Remote Procedure Call)
- Function-oriented: models API calls as method invocations
- Uses Protocol Buffers (binary) over HTTP/2 — smaller payloads, faster parsing
- Strongly typed contracts via `.proto` files
- Native streaming support (unary, server-stream, client-stream, bidirectional)
- Best for: internal microservice communication, high-throughput systems, polyglot environments

> **Bad design:** using REST for high-throughput internal services where binary efficiency matters, or using gRPC for a public API where browser clients and third-party consumers need broad compatibility.

---

## Key Design Principles

| Principle       | What it means                                                              |
|-----------------|----------------------------------------------------------------------------|
| **Consistency** | Same patterns everywhere — naming, error shapes, pagination, status codes |
| **Simplicity**  | Expose only what's needed; avoid leaking internal implementation detail   |
| **Security**    | Authentication, authorisation, input validation baked in from the start   |
| **Performance** | Consider payload size, caching, and call frequency at design time         |

> **Bad design:** every endpoint returns a different error shape, some use `snake_case` fields and others `camelCase`, pagination works differently per resource. Clients end up writing bespoke handling for every endpoint.

---

## How Protocols Influence API Design

The protocol you choose shapes every design decision downstream:

- **HTTP/1.1** — request/response, text-based headers, one request per connection by default
- **HTTP/2** — multiplexing, header compression, enables gRPC streaming
- **WebSockets** — persistent bidirectional connection, suited for real-time (chat, live feeds)
- **HTTP/3 (QUIC)** — lower latency, better on unreliable networks

Protocol choice determines: what interaction patterns are possible, how errors propagate, whether you can stream, and what clients can realistically connect.

> **Bad design:** polling a REST endpoint every second for live updates instead of using WebSockets — wastes connections, adds latency, and hammers the server unnecessarily.

---

## Design Approaches

### Contract-First
Write the API specification (OpenAPI, `.proto`) before writing any code. Enables parallel team work and forces upfront design decisions.

### Code-First
Write the implementation; the spec is generated from it. Faster to start, but requires discipline to avoid spec drift.

### When to use which
- Multiple teams / external consumers → contract-first
- Single team / internal service → code-first for speed

> **Bad design:** no spec at all — the API is whatever the code happens to do. Breaking changes ship silently, consumers find out in production, and there is no source of truth to review or diff.

---

## API Lifecycle Management

An API isn't a one-time decision — it evolves:

1. **Design** — define contracts, validate with stakeholders
2. **Build** — implement against the contract
3. **Publish** — document, version, expose to consumers
4. **Monitor** — track usage, error rates, latency
5. **Version** — introduce breaking changes under a new version (`/v2/`)
6. **Deprecate** — communicate sunset timelines, provide migration guides
7. **Retire** — remove old versions once consumers have migrated

> **Bad design:** making breaking changes directly to `/v1/` with no notice — renaming fields, removing endpoints, changing response shapes. Consumers break silently and have no migration path.

---

## OSI Model — Application Layer (Layer 7)

The OSI model has 7 layers. As API engineers, we operate primarily at **Layer 7 — the Application Layer**.

```
Layer 7 — Application   ← HTTP, gRPC, WebSocket (our layer)
Layer 6 — Presentation  ← TLS/SSL encryption, data encoding
Layer 5 — Session       ← Connection management
Layer 4 — Transport     ← TCP / UDP
Layer 3 — Network       ← IP routing
Layer 2 — Data Link     ← Ethernet, MAC addresses
Layer 1 — Physical      ← Cables, hardware
```

**Why it matters:** when debugging API issues, understanding which layer the problem is at (DNS resolution, TLS handshake failure, malformed HTTP, application logic) tells you where to look.

> **Bad design:** serving an API over plain HTTP (no TLS) because "it's internal" — traffic crosses Layer 6 unencrypted, meaning any network-level observer can read credentials and payloads.

---

## HTTP

### Common Methods
| Method    | Purpose                            | Idempotent | Safe |
|-----------|------------------------------------|------------|------|
| `GET`     | Retrieve a resource                | Yes        | Yes  |
| `POST`    | Create a resource / trigger action | No         | No   |
| `PUT`     | Replace a resource entirely        | Yes        | No   |
| `PATCH`   | Partially update a resource        | No         | No   |
| `DELETE`  | Remove a resource                  | Yes        | No   |

### Status Codes
| Range | Category      | Common examples                                                                                      |
|-------|---------------|------------------------------------------------------------------------------------------------------|
| 2xx   | Success       | `200 OK`, `201 Created`, `204 No Content`                                                            |
| 3xx   | Redirect      | `301 Moved Permanently`, `304 Not Modified`                                                          |
| 4xx   | Client error  | `400 Bad Request`, `401 Unauthorised`, `403 Forbidden`, `404 Not Found`, `429 Too Many Requests`    |
| 5xx   | Server error  | `500 Internal Server Error`, `502 Bad Gateway`, `503 Service Unavailable`                           |

### Common Headers
| Header            | Direction        | Purpose                                                      |
|-------------------|------------------|--------------------------------------------------------------|
| `Content-Type`    | Request/Response | Media type of the body (`application/json`)                  |
| `Accept`          | Request          | Media types the client can handle                            |
| `Authorization`   | Request          | Credentials (`Bearer <token>`)                               |
| `Cache-Control`   | Both             | Caching directives                                           |
| `X-Request-ID`    | Both             | Correlation ID for tracing requests                          |
| `Retry-After`     | Response         | How long to wait before retrying (used with 429/503)         |

> **Bad design:** returning `200 OK` with `{ "error": "not found" }` in the body — clients can't branch on HTTP status and must parse every response body to detect failure.

---

## Choosing the Right Protocol

| Consideration              | Lean REST/HTTP                         | Lean gRPC                              |
|----------------------------|----------------------------------------|----------------------------------------|
| **Interaction pattern**    | Request/response, public-facing        | Streaming, internal service mesh       |
| **Payload size**           | Small-medium, human-readable fine      | Large volumes, binary efficiency needed |
| **Performance**            | Moderate latency acceptable            | Low latency, high throughput required  |
| **Security needs**         | Standard TLS + JWT sufficient          | mTLS for service-to-service            |
| **Client compatibility**   | Browser, mobile, third-party tools     | Controlled internal clients            |
| **Developer experience**   | Wide familiarity, easy to debug        | Strongly typed, generated clients      |

> **Bad design:** defaulting to REST for everything — using it for a high-frequency internal telemetry pipeline where the JSON parsing overhead and HTTP/1.1 connection overhead becomes a real bottleneck.

---

## Resource Modelling

Good API design starts with modelling your business domain as resources.

### Business Domain → Resources
Identify the *nouns* in your domain: `user`, `order`, `payment`, `product`. These become your resources.

### REST Resource URL Patterns
```
/users                    collection
/users/{id}               single resource
/users/{id}/orders        nested resource (owned by user)
/orders?status=pending    filtering via query params
```

### Rules
- URLs identify *resources* (nouns), not actions (verbs)
- Use **plural nouns** for collections: `/users` not `/user`
- Nest only one level deep where possible — deeper nesting gets unwieldy
- Keep URLs lowercase and hyphen-separated: `/payment-methods` not `/paymentMethods`

> **Bad design:** verb-based URLs like `POST /getUser`, `POST /createOrder`, `GET /deleteAccount?id=1` — ignores HTTP semantics, breaks caching, and produces an unpredictable surface area.

---

## RESTful API Best Practices

- **Plural nouns** for all collection endpoints (`/orders`, `/products`)
- **Consistent URL structure** — same patterns across every resource
- **Version your API** from day one: `/v1/users` — changing a URL later is a breaking change
- **Support filtering, sorting, pagination** on collection endpoints:
  ```
  GET /orders?status=shipped&sort=created_at&limit=20&offset=40
  ```
- **Use the right status codes** — don't return `200 OK` with an error body
- **Consistent error shape** — clients should handle errors the same way everywhere:
  ```json
  {
    "error": {
      "code": "RESOURCE_NOT_FOUND",
      "message": "Order 123 does not exist"
    }
  }
  ```
- **Never expose internal IDs** — use UUIDs or opaque identifiers, not sequential integers

> **Bad design:** sequential integer IDs in URLs (`/users/1`, `/users/2`) — makes it trivial to enumerate all records, a pattern that led directly to the Optus and First American breaches. See [Notable API Incidents](notable-api-incidents.md).

---

## Authentication

**Authentication (AuthN)** = verifying *who* the caller is.

| Method                      | How it works                                                                                      |
|-----------------------------|---------------------------------------------------------------------------------------------------|
| **Basic Auth**              | `username:password` base64-encoded in every request header                                        |
| **Bearer tokens**           | `Authorization: Bearer <token>` — server validates the token on each request                      |
| **OAuth2**                  | Delegated auth protocol — user grants a client limited access via an auth server                  |
| **JWT**                     | Self-contained token: `header.payload.signature` — server verifies signature without a DB lookup  |
| **Access + Refresh tokens** | Short-lived access token (minutes/hours) paired with a long-lived refresh token                   |

> **Bad design:** tokens with no expiry, tokens passed in query parameters (they appear in server logs and browser history), or JWTs with sensitive data in an unencrypted payload — anyone who decodes the base64 can read it.

---

## Authorisation

**Authorisation (AuthZ)** = verifying *what* the authenticated caller is allowed to do.

> AuthN answers "who are you?" — AuthZ answers "what are you allowed to do?"

**RBAC — Role-Based Access Control**
Users are assigned roles (`admin`, `editor`, `viewer`); roles carry permissions. Simple, auditable, widely understood.
- When to use: most applications with clearly defined user types.

**ABAC — Attribute-Based Access Control**
Permissions are evaluated from attributes — user department, resource owner, time of day, IP range. Flexible but more complex to reason about.
- When to use: complex enterprise rules ("finance team can read reports during business hours").

**ACL — Access Control List**
An explicit list of who can do what on a specific resource — like file system permissions.
- When to use: per-document or per-object sharing (Google Docs model), object storage.

**Implementation patterns**
- Enforce AuthZ in middleware/interceptors, not scattered across business logic
- Centralise policy with a policy engine (e.g. OPA — Open Policy Agent)
- Always check at the **object level**, not just the route — failing to do so is BOLA/IDOR, the most common API breach

> **Bad design:** checking that a user is authenticated but not that they're authorised to access *that specific resource* — e.g. `GET /invoices/8842` returns the invoice to any logged-in user, not just the owner. This is exactly what caused the Optus and First American incidents.

---

## Security

**Rate Limiting**
Cap requests per client/IP per time window to prevent abuse, credential stuffing, and accidental or intentional DoS.
> **Bad design:** no rate limit on `POST /login` or `POST /password-reset` — an attacker can try millions of passwords or spam password reset emails.

**CORS (Cross-Origin Resource Sharing)**
Browser mechanism controlling which origins can make requests to your API. Misconfigured CORS on an authenticated API allows any website to make requests on behalf of a logged-in user.
> **Bad design:** `Access-Control-Allow-Origin: *` on a private API, or adding CORS headers to silence browser errors without understanding what access is being granted.

**SQL Injection**
Unsanitised user input concatenated into SQL queries lets an attacker execute arbitrary SQL — reading, modifying, or deleting data.
> **Bad design:** `"SELECT * FROM users WHERE id = " + userInput` — always use parameterised queries / prepared statements.

**CSRF (Cross-Site Request Forgery)**
Tricks a user's browser into making an unintended authenticated request — e.g. a fund transfer — by exploiting an active session cookie.
> **Bad design:** stateful cookie auth with no CSRF token and no `SameSite` cookie attribute on state-changing endpoints.

**Firewalls**
Network-level rules that block traffic before it reaches your application. Internal services — databases, admin APIs, message brokers — should never be publicly reachable.
> **Bad design:** database port `5432` open to `0.0.0.0/0` in a cloud security group — your database is on the public internet.

**VPN**
Encrypted tunnel for accessing private networks. Internal APIs that aren't meant to be public should sit behind a VPN, not just an undocumented URL.
> **Bad design:** "security through obscurity" — an internal API assumed safe because its URL isn't published. URLs get leaked through logs, referrer headers, and browser history.

---

> For real-world consequences of ignoring these principles, see the [Notable API Incidents](notable-api-incidents.md)
