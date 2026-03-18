# Noteworthy API Incidents: Lessons in Bad API Design

A reference document for understanding real-world consequences of poor API design decisions.

---

## 1. Optus Data Breach (2022)

**Category:** Broken Object Level Authorization (BOLA / IDOR) + Missing Authentication

### What happened

An unauthenticated API endpoint at `/users/{userId}` was publicly accessible. An attacker iterated through sequential user IDs and downloaded the personal records of approximately **9.8 million Australians** — nearly 40% of the country's population.

### The API design failure

```
GET /users/10001   → returns full user record
GET /users/10002   → returns full user record
GET /users/10003   → returns full user record
...
```

- **No authentication required** — the endpoint was reachable without any token or session
- **Sequential, predictable identifiers** — IDs were numeric and incremental, making enumeration trivial
- **No rate limiting** — thousands of requests could be made without triggering any throttle or alert
- **Excessive data exposure** — the response returned the full user object including passport numbers, driver's licence numbers, dates of birth, and home addresses

### Data exposed

- Full name, date of birth, home address
- Email address, phone number
- Government-issued ID numbers (passport, driver's licence)

### Key design principles violated

| Principle | Violation |
|---|---|
| Authentication on all endpoints | Endpoint required no credentials |
| Authorisation (object-level) | No check that the caller owned the requested resource |
| Non-sequential identifiers | Integer IDs made enumeration trivial |
| Rate limiting | No throttle prevented bulk harvesting |
| Minimal response payload | Full PII returned regardless of consumer need |

### How it should have been designed

```
// Require authentication
GET /users/{userId}
Authorization: Bearer <token>

// Server-side: verify the authenticated user IS userId (or has admin role)
if authenticatedUser.ID != userId && !authenticatedUser.IsAdmin {
    return 403 Forbidden
}

// Use UUIDs, not sequential integers
GET /users/f47ac10b-58cc-4372-a567-0e02b2c3d479
```

---

## 2. Coinbase Trading API Vulnerability (2022)

**Category:** Broken Function Level Authorization + Insufficient Input Validation

### What happened

A researcher discovered a flaw in Coinbase's Advanced Trade API that allowed a user to place an order that manipulated their own portfolio in an economically invalid way — swapping an asset they **owned** for an asset they **did not own**, at a net gain. By crafting a specific API request, they were able to effectively "sell" Ethereum they held and "buy" Bitcoin they did not hold, with the trade resolving in their favour.

### The API design failure

The order endpoint accepted a request body where the asset being sold and the asset being bought could be set independently, without the backend fully validating cross-asset ownership and balance constraints atomically:

```json
POST /api/v3/brokerage/orders
{
  "side": "SELL",
  "product_id": "ETH-BTC",
  "base_size": "1.0"
}
```

- **Insufficient server-side validation** — the API accepted combinations of parameters that should have been rejected before any trade logic executed
- **Non-atomic balance checks** — ownership and balance of both legs of the trade were not verified together in a single consistent transaction
- **Missing cross-field invariant enforcement** — the API did not enforce that a sell order requires ownership of the asset being sold *and* that a buy order requires sufficient funds

### Key design principles violated

| Principle | Violation |
|---|---|
| Atomic validation | Balance checks on both legs were not performed together |
| Invariant enforcement | API allowed economically impossible state |
| Input validation at boundary | Invalid asset combinations were accepted |
| Idempotency safeguards | No secondary check before settlement |

### How it should have been designed

```
// Before executing any trade:
// 1. Validate the caller owns sufficient base asset (sell leg)
// 2. Validate the caller has sufficient quote asset or credit (buy leg)
// 3. Lock both balances atomically before order processing
// 4. Roll back entirely if any invariant fails

BEGIN TRANSACTION
  CHECK seller_balance >= order.base_size         → fail fast if not met
  CHECK buyer_balance  >= order.quote_size        → fail fast if not met
  EXECUTE trade
COMMIT (or ROLLBACK on any error)
```

---

## 3. First American Financial Document API (2019)

**Category:** Broken Object Level Authorization (IDOR) — Financial Services

### What happened

First American Financial Corporation, one of the largest title insurance and real estate settlement companies in the United States, exposed approximately **885 million sensitive financial documents** through a single API design flaw. The vulnerability was discovered by a real estate developer who noticed he could increment a number in the URL and access documents belonging to other customers.

### The API design failure

```
GET /title/search?map=000000001   → returns mortgage document
GET /title/search?map=000000002   → returns someone else's mortgage document
GET /title/search?map=000000003   → returns someone else's bank records
...
```

- **Sequential numeric document IDs** in the query parameter — trivially enumerable
- **No authorisation check** — the API did not verify whether the authenticated user had any relationship to the requested document
- **Authenticated but not authorised** — a user could log in legitimately and then access any document in the system

### Data exposed

- Bank account numbers and statements
- Mortgage and loan documents
- Social Security numbers
- Tax records
- Property deeds and closing statements
- Wire transfer instructions (high value for fraud)

### Key design principles violated

| Principle | Violation |
|---|---|
| Object-level authorisation | No ownership check on each document request |
| Non-guessable identifiers | Sequential integers allowed full enumeration |
| Least privilege | Any authenticated user could access any document |
| Audit logging | Mass enumeration went undetected |

### How it should have been designed

```
// Use non-sequential, non-guessable identifiers
GET /documents/f47ac10b-58cc-4372-a567-0e02b2c3d479

// Server-side: verify the document belongs to the authenticated user
document := db.GetDocument(documentID)
if document.OwnerID != authenticatedUser.ID {
    auditLog.Warn("unauthorized document access attempt", ...)
    return 403 Forbidden
}
```

---

## Common Themes

All three incidents share the same underlying failures despite being in different companies and contexts:

| Theme | Optus | Coinbase | First American |
|---|---|---|---|
| Missing / insufficient authorisation | ✓ | ✓ | ✓ |
| Predictable / enumerable identifiers | ✓ | | ✓ |
| Insufficient server-side validation | ✓ | ✓ | |
| No rate limiting / anomaly detection | ✓ | | ✓ |
| Excessive data in response | ✓ | | ✓ |

### OWASP API Security Top 10 mappings

- **API1:2023 — Broken Object Level Authorization** — Optus, First American
- **API3:2023 — Broken Object Property Level Authorization** — Coinbase
- **API5:2023 — Broken Function Level Authorization** — Coinbase
- **API6:2023 — Unrestricted Access to Sensitive Business Flows** — Coinbase
- **API8:2023 — Security Misconfiguration** — Optus (no auth on endpoint)

---

> These incidents are why API security must be treated as a first-class design concern — not a post-deployment checklist item.
